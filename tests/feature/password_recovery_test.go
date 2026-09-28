package feature

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/YahirHub/axolotl-dashboard-template/app/facades"
	"github.com/YahirHub/axolotl-dashboard-template/app/models"
	"github.com/YahirHub/axolotl-dashboard-template/app/security"
	"github.com/YahirHub/axolotl-dashboard-template/tests"
)

func TestPasswordRecoveryIsHiddenWhenMailIsNotConfigured(t *testing.T) {
	oldMail := facades.Config().Get("mail")
	t.Cleanup(func() { facades.Config().Add("mail", oldMail) })
	facades.Config().Add("mail", map[string]any{
		"host":       "",
		"port":       1025,
		"username":   "",
		"password":   "",
		"encryption": "auto",
		"from": map[string]any{
			"address": "",
			"name":    "Axolotl",
		},
	})

	resetUsers(t)
	createAdmin(t)
	var testCase tests.TestCase

	forgot, err := testCase.Http(t).Get("/forgot-password")
	require.NoError(t, err)
	forgot.AssertStatus(303).AssertHeader("Location", "/login")

	login, err := testCase.Http(t).Get("/login")
	require.NoError(t, err)
	login.AssertOk().AssertDontSee([]string{"¿Olvidaste tu contraseña?"})
}

func TestPasswordRecoveryEmailAndOneTimeReset(t *testing.T) {
	port, messages := startFakeSMTP(t)
	oldMail := facades.Config().Get("mail")
	t.Cleanup(func() { facades.Config().Add("mail", oldMail) })
	facades.Config().Add("mail", map[string]any{
		"host":       "127.0.0.1",
		"port":       port,
		"username":   "",
		"password":   "",
		"encryption": "none",
		"from": map[string]any{
			"address": "no-reply@example.test",
			"name":    "Axolotl Test",
		},
	})

	resetUsers(t)
	admin := createAdmin(t)
	facades.Cache().Forget("password-recovery:" + security.IdentityFingerprint(admin.Email))

	var testCase tests.TestCase
	forgotPage, err := testCase.Http(t).Get("/forgot-password")
	require.NoError(t, err)
	forgotPage.AssertOk().AssertSee([]string{"Recuperar cuenta", "Enviar enlace"})

	request := postForm(
		t,
		"/forgot-password",
		url.Values{"email": {admin.Email}},
		forgotPage.Cookies(),
		forgotPage.Headers().Get("X-CSRF-TOKEN"),
	)
	request.AssertStatus(303).AssertHeader("Location", "/forgot-password")

	var message string
	select {
	case message = <-messages:
	case <-time.After(5 * time.Second):
		t.Fatal("no se recibió correo SMTP de recuperación")
	}

	require.Contains(t, message, "Restablecer")
	require.NotContains(t, message, admin.Password)

	tokenPattern := regexp.MustCompile(`#token=([A-Za-z0-9_-]{40,128})`)
	match := tokenPattern.FindStringSubmatch(message)
	require.Len(t, match, 2)
	token := match[1]

	resetPage, err := testCase.Http(t).Get("/reset-password")
	require.NoError(t, err)
	resetPage.AssertOk().
		AssertHeader("Cache-Control", "no-store, max-age=0").
		AssertHeader("Referrer-Policy", "no-referrer").
		AssertSee([]string{"Crear nueva contraseña", "Actualizar contraseña", "passwordResetToken"})

	oldSession := map[string]any{
		"auth_user_id": fmt.Sprint(admin.ID),
		"auth_version": fmt.Sprint(admin.AuthVersion),
	}

	resetResponse := postForm(
		t,
		"/reset-password",
		url.Values{
			"token":                 {token},
			"password":              {"NewStrongPass123!"},
			"password_confirmation": {"NewStrongPass123!"},
		},
		resetPage.Cookies(),
		resetPage.Headers().Get("X-CSRF-TOKEN"),
	)
	resetResponse.AssertStatus(303).AssertHeader("Location", "/login")

	var updated models.User
	require.NoError(t, facades.Orm().Query().Where("id", admin.ID).First(&updated))
	require.True(t, facades.Hash().Check("NewStrongPass123!", updated.Password))
	require.Greater(t, updated.AuthVersion, admin.AuthVersion)

	manager := security.NewPasswordResetManager()
	_, err = manager.Validate(token)
	require.ErrorIs(t, err, security.ErrResetTokenInvalid)

	oldSessionResponse, err := testCase.Http(t).WithSession(oldSession).Get("/dashboard")
	require.NoError(t, err)
	oldSessionResponse.AssertStatus(303).AssertHeader("Location", "/login")

	loginPage, err := testCase.Http(t).Get("/login")
	require.NoError(t, err)
	login := postForm(
		t,
		"/login",
		url.Values{
			"email":    {admin.Email},
			"password": {"NewStrongPass123!"},
		},
		loginPage.Cookies(),
		loginPage.Headers().Get("X-CSRF-TOKEN"),
	)
	login.AssertStatus(303).AssertHeader("Location", "/dashboard")
}

func TestPasswordResetManagerRejectsExpiredToken(t *testing.T) {
	resetUsers(t)
	admin := createAdmin(t)

	manager := security.NewPasswordResetManager()
	token, err := manager.Issue(admin.ID, 5*time.Minute)
	require.NoError(t, err)

	var record models.PasswordResetToken
	require.NoError(t, facades.Orm().Query().Where("token_hash <> ?", "").First(&record))
	record.ExpiresAt = time.Now().UTC().Add(-time.Minute)
	require.NoError(t, facades.Orm().Query().Save(&record))

	_, err = manager.Validate(token)
	require.True(t, errors.Is(err, security.ErrResetTokenInvalid))
}
