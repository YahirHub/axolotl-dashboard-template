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
	session := map[string]any{"auth_user_id": fmt.Sprint(admin.ID)}

	page, err := testCase.Http(t).WithSession(session).Get("/archivos")
	require.NoError(t, err)
	page.
		AssertOk().
		AssertSee([]string{
			"Explorador de archivos",
			"Nueva carpeta",
			"Nuevo archivo",
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

	remove := postForm(t, "/api/files/delete", url.Values{
		"path": {"api-note.txt"},
	}, cookies, csrf)
	remove.AssertOk()

	missing, err := testCase.Http(t).
		WithCookies(cookies).
		Get("/api/files/read?path=api-note.txt")
	require.NoError(t, err)
	missing.AssertNotFound()
}
