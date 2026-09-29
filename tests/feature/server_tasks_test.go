package feature

import (
	"fmt"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/YahirHub/axolotl-dashboard-template/tests"
)

func TestServerRoutesAndTaskTemplates(t *testing.T) {
	resetUsers(t)
	admin := createAdmin(t)
	var testCase tests.TestCase
	session := map[string]any{
		"auth_user_id": fmt.Sprint(admin.ID),
		"auth_version": fmt.Sprint(admin.AuthVersion),
	}

	page, err := testCase.Http(t).WithSession(session).Get("/servidor/tareas")
	require.NoError(t, err)
	page.AssertOk().AssertSee([]string{
		"Administrador de tareas",
		"Nueva plantilla",
		"systemd",
	})

	body, err := page.Content()
	require.NoError(t, err)
	require.Contains(t, body, "data-route=\"server-tasks\"")
	require.NotContains(t, sidebarHTML(body), "#")
	require.NotContains(t, sidebarHTML(body), "data-action=")
	require.Contains(t, sidebarHTML(body), "href=\"/servidor\"")
	require.Contains(t, sidebarHTML(body), "href=\"/servidor/tareas\"")
	require.Contains(t, sidebarHTML(body), "href=\"/servidor/entorno\"")

	csrf := page.Headers().Get("X-CSRF-TOKEN")
	require.NotEmpty(t, csrf)
	cookies := page.Cookies()

	createService := postForm(t, "/api/server/tasks/create", url.Values{
		"name":              {"demo-worker"},
		"kind":              {"service"},
		"description":       {"Worker de prueba"},
		"exec_start":        {"/usr/bin/php /var/www/app/artisan queue:work"},
		"working_directory": {"/var/www/app"},
		"user":              {"www-data"},
		"restart":           {"on-failure"},
		"after":             {"network.target"},
	}, cookies, csrf)
	createService.AssertCreated()

	list, err := testCase.Http(t).WithCookies(cookies).Get("/api/server/tasks")
	require.NoError(t, err)
	list.AssertOk().AssertSee([]string{"demo-worker.service"})

	read, err := testCase.Http(t).WithCookies(cookies).Get("/api/server/tasks/read?filename=demo-worker.service")
	require.NoError(t, err)
	read.AssertOk().AssertSee([]string{
		"demo-worker.service",
		"ExecStart=/usr/bin/php /var/www/app/artisan queue:work",
	})

	unsafe := postForm(t, "/api/server/tasks/create", url.Values{
		"name":       {"unsafe"},
		"kind":       {"service"},
		"exec_start": {"/usr/bin/true\nEnvironment=BAD=1"},
	}, cookies, csrf)
	unsafe.AssertBadRequest()

	remove := postForm(t, "/api/server/tasks/delete", url.Values{
		"filename": {"demo-worker.service"},
	}, cookies, csrf)
	remove.AssertOk()

	environment, err := testCase.Http(t).WithCookies(cookies).Get("/servidor/entorno")
	require.NoError(t, err)
	environment.AssertOk().AssertSee([]string{
		"Entorno",
		"Runtime",
		"GOMAXPROCS",
	})
	environmentBody, err := environment.Content()
	require.NoError(t, err)
	require.Contains(t, environmentBody, "data-route=\"server-environment\"")
}

func TestSidebarContainsOnlyRouteNavigation(t *testing.T) {
	resetUsers(t)
	admin := createAdmin(t)
	var testCase tests.TestCase
	session := map[string]any{
		"auth_user_id": fmt.Sprint(admin.ID),
		"auth_version": fmt.Sprint(admin.AuthVersion),
	}

	response, err := testCase.Http(t).WithSession(session).Get("/dashboard")
	require.NoError(t, err)
	response.AssertOk()

	body, err := response.Content()
	require.NoError(t, err)
	sidebar := sidebarHTML(body)
	require.NotEmpty(t, sidebar)
	require.NotContains(t, sidebar, "href=\"#")
	require.NotContains(t, sidebar, "/dashboard#")
	require.NotContains(t, sidebar, "data-action=")
	require.Contains(t, sidebar, "href=\"/dashboard\"")
	require.Contains(t, sidebar, "href=\"/servidor\"")
	require.Contains(t, sidebar, "href=\"/servidor/tareas\"")
	require.Contains(t, sidebar, "href=\"/servidor/entorno\"")
	require.Contains(t, sidebar, "href=\"/archivos\"")
	require.Contains(t, sidebar, "href=\"/cuenta\"")
	require.Contains(t, sidebar, "href=\"/componentes\"")
}

func sidebarHTML(body string) string {
	start := strings.Index(body, "<nav class=\"nav\"")
	if start < 0 {
		return ""
	}
	end := strings.Index(body[start:], "</nav>")
	if end < 0 {
		return ""
	}
	return body[start : start+end+len("</nav>")]
}
