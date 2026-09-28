package feature

import (
	"fmt"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/YahirHub/axolotl-dashboard-template/tests"
)

func TestFileExplorerAPI(t *testing.T) {
	resetUsers(t)
	admin := createAdmin(t)
	var testCase tests.TestCase
	session := map[string]any{
		"auth_user_id": fmt.Sprint(admin.ID),
		"auth_version": fmt.Sprint(admin.AuthVersion),
	}

	page, err := testCase.Http(t).WithSession(session).Get("/archivos")
	require.NoError(t, err)
	page.
		AssertOk().
		AssertHeader("Content-Security-Policy", contentSecurityPolicyForTest()).
		AssertSee([]string{
			"Explorador de archivos",
			"Nueva carpeta",
			"Nuevo archivo",
			"sidebar-collapsed",
		})

	body, err := page.Content()
	require.NoError(t, err)
	require.True(t, strings.Contains(body, `data-route="files"`))
	require.True(t, strings.Contains(body, `aria-current="page"`))

	csrf := page.Headers().Get("X-CSRF-TOKEN")
	require.NotEmpty(t, csrf)
	cookies := page.Cookies()

	create := postForm(t, "/api/files/create", url.Values{
		"parent": {""},
		"name":   {"api-note.txt"},
		"type":   {"file"},
	}, cookies, csrf)
	create.AssertCreated()

	routeFolder := postForm(t, "/api/files/create", url.Values{
		"parent": {""},
		"name":   {"route-test"},
		"type":   {"folder"},
	}, cookies, csrf)
	routeFolder.AssertCreated()

	nested, err := testCase.Http(t).WithCookies(cookies).Get("/archivos/route-test")
	require.NoError(t, err)
	nested.AssertOk().AssertSee([]string{"sidebar-collapsed"})
	nestedBody, err := nested.Content()
	require.NoError(t, err)
	require.Contains(t, nestedBody, `data-current-path="route-test"`)

	save := postForm(t, "/api/files/save", url.Values{
		"path":    {"api-note.txt"},
		"content": {"Contenido desde la prueba del backend Go."},
	}, cookies, csrf)
	save.AssertOk()

	read, err := testCase.Http(t).
		WithCookies(cookies).
		Get("/api/files/read?path=api-note.txt")
	require.NoError(t, err)
	read.AssertOk()

	payload, err := read.Json()
	require.NoError(t, err)
	require.Equal(t, "api-note.txt", payload["path"])
	require.Equal(t, "Contenido desde la prueba del backend Go.", payload["content"])
	require.Equal(t, "plaintext", payload["language"])

	traversal, err := testCase.Http(t).
		WithCookies(cookies).
		Get("/api/files/read?path=..%2Fsecret.txt")
	require.NoError(t, err)
	traversal.AssertBadRequest()

	doubleEncoded, err := testCase.Http(t).
		WithCookies(cookies).
		Get("/api/files/read?path=%252e%252e%252fsecret.txt")
	require.NoError(t, err)
	doubleEncoded.AssertBadRequest()

	absolute, err := testCase.Http(t).
		WithCookies(cookies).
		Get("/api/files?path=%2Fetc")
	require.NoError(t, err)
	absolute.AssertBadRequest()

	reserved := postForm(t, "/api/files/create", url.Values{
		"parent": {""},
		"name":   {"CON.txt"},
		"type":   {"file"},
	}, cookies, csrf)
	reserved.AssertBadRequest()

	xssName := postForm(t, "/api/files/create", url.Values{
		"parent": {""},
		"name":   {"<script>.txt"},
		"type":   {"file"},
	}, cookies, csrf)
	xssName.AssertBadRequest()

	remove := postForm(t, "/api/files/delete", url.Values{
		"path": {"api-note.txt"},
	}, cookies, csrf)
	remove.AssertOk()

	removeFolder := postForm(t, "/api/files/delete", url.Values{
		"path": {"route-test"},
	}, cookies, csrf)
	removeFolder.AssertOk()

	missing, err := testCase.Http(t).
		WithCookies(cookies).
		Get("/api/files/read?path=api-note.txt")
	require.NoError(t, err)
	missing.AssertNotFound()
}

func contentSecurityPolicyForTest() string {
	return "default-src 'self'; base-uri 'self'; form-action 'self'; frame-ancestors 'self'; object-src 'none'; script-src 'self' https://cdn.jsdelivr.net https://cdnjs.cloudflare.com; style-src 'self' 'unsafe-inline' https://cdn.jsdelivr.net https://cdnjs.cloudflare.com; font-src 'self' data: https://cdnjs.cloudflare.com; img-src 'self' data:; connect-src 'self'; media-src 'self'; worker-src 'self'; manifest-src 'self'"
}
