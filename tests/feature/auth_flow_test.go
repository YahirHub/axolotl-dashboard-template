package feature

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"testing"

	contractshttp "github.com/goravel/framework/contracts/testing/http"
	"github.com/stretchr/testify/require"

	"github.com/YahirHub/axolotl-dashboard-template/app/facades"
	"github.com/YahirHub/axolotl-dashboard-template/app/models"
	"github.com/YahirHub/axolotl-dashboard-template/app/security"
	"github.com/YahirHub/axolotl-dashboard-template/tests"
)

func postForm(t *testing.T, path string, values url.Values, cookies []*http.Cookie, csrf string) contractshttp.Response {
	t.Helper()
	var testCase tests.TestCase
	request := testCase.Http(t).
		WithCookies(cookies).
		WithHeader("Content-Type", "application/x-www-form-urlencoded").
		WithHeader("X-CSRF-TOKEN", csrf)

	response, err := request.Post(path, strings.NewReader(values.Encode()))
	require.NoError(t, err)
	return response
}

func TestSetupCreatesAdminAndStartsSession(t *testing.T) {
	resetUsers(t)

	pin, err := security.InitializeAdminPIN()
	require.NoError(t, err)
	require.Len(t, pin, 6)

	var testCase tests.TestCase
	setupPage, err := testCase.Http(t).Get("/setup")
	require.NoError(t, err)
	setupPage.AssertOk()

	values := url.Values{
		"name":                  {"Primer Admin"},
		"email":                 {"owner@example.test"},
		"alias":                 {"PA"},
		"password":              {"StrongAdmin123!"},
		"password_confirmation": {"StrongAdmin123!"},
	}
	for i, digit := range pin {
		values.Set(fmt.Sprintf("pin_%d", i+1), string(digit))
	}

	setupResponse := postForm(
		t,
		"/setup",
		values,
		setupPage.Cookies(),
		setupPage.Headers().Get("X-CSRF-TOKEN"),
	)
	setupResponse.
		AssertStatus(303).
		AssertHeader("Location", "/dashboard")

	var admin models.User
	require.NoError(t, facades.Orm().Query().Where("email", "owner@example.test").First(&admin))
	require.Equal(t, models.RoleAdmin, admin.Role)
	require.True(t, admin.Active)
	require.NotEqual(t, "StrongAdmin123!", admin.Password)
	require.True(t, facades.Hash().Check("StrongAdmin123!", admin.Password))

	dashboard, err := testCase.Http(t).WithCookies(setupResponse.Cookies()).Get("/dashboard")
	require.NoError(t, err)
	dashboard.AssertOk().AssertSee([]string{"Primer Admin", "Administrador", "Dashboard"})
}

func TestSetupRejectsIncorrectPIN(t *testing.T) {
	resetUsers(t)

	pin, err := security.InitializeAdminPIN()
	require.NoError(t, err)
	require.Len(t, pin, 6)

	wrong := []byte(pin)
	if wrong[0] == '0' {
		wrong[0] = '1'
	} else {
		wrong[0] = '0'
	}

	var testCase tests.TestCase
	setupPage, err := testCase.Http(t).Get("/setup")
	require.NoError(t, err)

	values := url.Values{
		"name":                  {"Blocked Admin"},
		"email":                 {"blocked@example.test"},
		"alias":                 {"BA"},
		"password":              {"StrongAdmin123!"},
		"password_confirmation": {"StrongAdmin123!"},
	}
	for i, digit := range wrong {
		values.Set(fmt.Sprintf("pin_%d", i+1), string(digit))
	}

	response := postForm(
		t,
		"/setup",
		values,
		setupPage.Cookies(),
		setupPage.Headers().Get("X-CSRF-TOKEN"),
	)
	response.AssertOk().AssertSee([]string{"PIN temporal incorrecto"})

	count, err := facades.Orm().Query().Model(&models.User{}).Count()
	require.NoError(t, err)
	require.Zero(t, count)
}

func TestRegistrationCreatesOnlyStandardUser(t *testing.T) {
	resetUsers(t)
	createAdmin(t)

	var testCase tests.TestCase
	registerPage, err := testCase.Http(t).Get("/register")
	require.NoError(t, err)
	registerPage.AssertOk()

	values := url.Values{
		"name":                  {"Usuario Demo"},
		"email":                 {"user@example.test"},
		"alias":                 {"UD"},
		"password":              {"StrongUser123!"},
		"password_confirmation": {"StrongUser123!"},
		"role":                  {"admin"},
	}

	response := postForm(
		t,
		"/register",
		values,
		registerPage.Cookies(),
		registerPage.Headers().Get("X-CSRF-TOKEN"),
	)
	response.
		AssertStatus(303).
		AssertHeader("Location", "/dashboard")

	var user models.User
	require.NoError(t, facades.Orm().Query().Where("email", "user@example.test").First(&user))
	require.Equal(t, models.RoleUser, user.Role)
	require.False(t, user.IsAdmin())
	require.True(t, facades.Hash().Check("StrongUser123!", user.Password))
}

func TestLoginProgressiveBlock(t *testing.T) {
	resetUsers(t)
	admin := createAdmin(t)

	var testCase tests.TestCase
	loginPage, err := testCase.Http(t).Get("/login")
	require.NoError(t, err)
	loginPage.AssertOk()

	values := url.Values{
		"email":    {admin.Email},
		"password": {"WrongPassword123!"},
	}

	var last contractshttp.Response
	for i := 0; i < 5; i++ {
		last = postForm(
			t,
			"/login",
			values,
			loginPage.Cookies(),
			loginPage.Headers().Get("X-CSRF-TOKEN"),
		)
	}

	require.NotNil(t, last)
	last.AssertOk().AssertSee([]string{"Demasiados intentos"})
	require.NotEmpty(t, last.Headers().Get("Retry-After"))

	correctWhileBlocked := postForm(
		t,
		"/login",
		url.Values{
			"email":    {admin.Email},
			"password": {"AdminPass123!"},
		},
		loginPage.Cookies(),
		loginPage.Headers().Get("X-CSRF-TOKEN"),
	)
	correctWhileBlocked.AssertOk().AssertSee([]string{"Demasiados intentos"})
}

func TestLoginLogoutSessionFlow(t *testing.T) {
	resetUsers(t)
	admin := createAdmin(t)

	var testCase tests.TestCase
	loginPage, err := testCase.Http(t).Get("/login")
	require.NoError(t, err)
	loginPage.AssertOk()

	loginValues := url.Values{
		"email":    {admin.Email},
		"password": {"AdminPass123!"},
	}
	loginResponse := postForm(
		t,
		"/login",
		loginValues,
		loginPage.Cookies(),
		loginPage.Headers().Get("X-CSRF-TOKEN"),
	)
	loginResponse.
		AssertStatus(303).
		AssertHeader("Location", "/dashboard")

	dashboard, err := testCase.Http(t).WithCookies(loginResponse.Cookies()).Get("/dashboard")
	require.NoError(t, err)
	dashboard.AssertOk()

	logoutResponse := postForm(
		t,
		"/logout",
		url.Values{},
		loginResponse.Cookies(),
		dashboard.Headers().Get("X-CSRF-TOKEN"),
	)
	logoutResponse.
		AssertStatus(303).
		AssertHeader("Location", "/login")
}
