package feature

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/YahirHub/axolotl-dashboard-template/app/facades"
	"github.com/YahirHub/axolotl-dashboard-template/app/models"
	"github.com/YahirHub/axolotl-dashboard-template/app/security"
	"github.com/YahirHub/axolotl-dashboard-template/database/migrations"
	"github.com/YahirHub/axolotl-dashboard-template/tests"
)

func resetUsers(t *testing.T) {
	t.Helper()
	security.LoginAttempts().Reset()

	recoveryMigration := &migrations.M20260928000300AddAuthVersionAndPasswordResetTokens{}
	userMigration := &migrations.M20260928000200CreateUsersTable{}

	_ = recoveryMigration.Down()
	_ = userMigration.Down()
	require.NoError(t, userMigration.Up())
	require.NoError(t, recoveryMigration.Up())
}

func createAdmin(t *testing.T) models.User {
	t.Helper()
	hashed, err := facades.Hash().Make("AdminPass123!")
	require.NoError(t, err)

	user := models.User{
		Name:        "Admin Test",
		Email:       "admin@example.test",
		Alias:       "AT",
		Password:    hashed,
		Role:        models.RoleAdmin,
		Active:      true,
		AuthVersion: 1,
	}
	require.NoError(t, facades.Orm().Query().Create(&user))
	return user
}

func TestFirstRunRedirectsToSetup(t *testing.T) {
	resetUsers(t)
	var testCase tests.TestCase

	response, err := testCase.Http(t).Get("/")
	require.NoError(t, err)
	response.AssertStatus(303).AssertSee([]string{"/setup"})

	setup, err := testCase.Http(t).Get("/setup")
	require.NoError(t, err)
	setup.AssertOk().AssertSee([]string{
		"Crear administrador",
		"PIN temporal",
		"pin_1",
		"pin_6",
	})
}

func TestGuestRoutesAfterAdminExists(t *testing.T) {
	resetUsers(t)
	createAdmin(t)
	var testCase tests.TestCase

	root, err := testCase.Http(t).Get("/")
	require.NoError(t, err)
	root.AssertStatus(303).AssertSee([]string{"/login"})

	login, err := testCase.Http(t).Get("/login")
	require.NoError(t, err)
	login.AssertOk().AssertSee([]string{"Iniciar sesión", "Correo", "Contraseña"})

	register, err := testCase.Http(t).Get("/register")
	require.NoError(t, err)
	register.AssertOk().AssertSee([]string{"Crear cuenta", "rol", "Usuario"})
}

func TestProtectedPagesRequireAuthentication(t *testing.T) {
	resetUsers(t)
	createAdmin(t)
	var testCase tests.TestCase

	response, err := testCase.Http(t).Get("/dashboard")
	require.NoError(t, err)
	response.AssertStatus(303).AssertSee([]string{"/login"})
}

func TestAuthenticatedPages(t *testing.T) {
	resetUsers(t)
	admin := createAdmin(t)
	var testCase tests.TestCase
	session := map[string]any{
		"auth_user_id": fmt.Sprint(admin.ID),
		"auth_version": fmt.Sprint(admin.AuthVersion),
	}

	for _, item := range []struct {
		path   string
		active string
		see    []string
	}{
		{"/dashboard", "dashboard", []string{"Dashboard", "Proyectos destacados", "Tareas prioritarias", "Equipo", "Agenda", "Integraciones"}},
		{"/cuenta", "account", []string{"Cuenta", "Perfil", "Preferencias", "Seguridad"}},
		{"/servidor", "server", []string{"Servidor", "Métricas del servidor", "CPU", "RAM", "Red", "Almacenamiento"}},
		{"/componentes", "components", []string{"Componentes UI", "Design system reutilizable", "Laboratorio UI"}},
		{"/archivos", "files", []string{"Explorador de archivos", "Nueva carpeta", "Nuevo archivo", "Acceso rápido"}},
	} {
		response, err := testCase.Http(t).WithSession(session).Get(item.path)
		require.NoError(t, err)
		response.
			AssertOk().
			AssertHeader("X-Content-Type-Options", "nosniff").
			AssertSee(item.see)

		body, err := response.Content()
		require.NoError(t, err)
		require.Contains(t, body, `data-route="`+item.active+`"`)
		require.Contains(t, body, `aria-current="page"`)
		require.True(t, strings.Contains(body, "nav-item active"), "la ruta %s debe renderizar una opción activa", item.path)
	}
}

func TestHealthRoute(t *testing.T) {
	var testCase tests.TestCase

	response, err := testCase.Http(t).Get("/health")
	require.NoError(t, err)

	response.
		AssertOk().
		AssertHeader("X-Content-Type-Options", "nosniff").
		AssertJson(map[string]any{
			"status":   "ok",
			"database": "ok",
		})
}
