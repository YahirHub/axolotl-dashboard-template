package feature

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/YahirHub/axolotl-dashboard-template/tests"
)

func TestDashboardRoutes(t *testing.T) {
	var testCase tests.TestCase

	for _, path := range []string{"/", "/dashboard"} {
		response, err := testCase.Http(t).Get(path)
		require.NoError(t, err)

		response.
			AssertOk().
			AssertHeader("X-Content-Type-Options", "nosniff").
			AssertHeader("X-Frame-Options", "SAMEORIGIN").
			AssertSee([]string{
				"Dashboard",
				"Proyectos destacados",
				"Tareas prioritarias",
				"Equipo",
				"Agenda",
				"Integraciones",
				"/assets/css/app.css",
				"/assets/js/app.js",
				"/assets/brand/axolotl.svg",
				"/assets/brand/axolotl.png",
			})
	}
}

func TestAccountRoute(t *testing.T) {
	var testCase tests.TestCase

	response, err := testCase.Http(t).Get("/cuenta")
	require.NoError(t, err)

	response.
		AssertOk().
		AssertSee([]string{
			"Cuenta",
			"Perfil",
			"Preferencias",
			"Seguridad",
			"Guardar perfil",
		})
}

func TestServerRoute(t *testing.T) {
	var testCase tests.TestCase

	response, err := testCase.Http(t).Get("/servidor")
	require.NoError(t, err)

	response.
		AssertOk().
		AssertSee([]string{
			"Servidor",
			"Métricas del servidor",
			"CPU",
			"RAM",
			"Red",
			"Almacenamiento",
		})
}

func TestComponentsRoute(t *testing.T) {
	var testCase tests.TestCase

	response, err := testCase.Http(t).Get("/componentes")
	require.NoError(t, err)

	response.
		AssertOk().
		AssertSee([]string{
			"Componentes UI",
			"Design system reutilizable",
			"Laboratorio UI",
			"Botones y tokens semánticos",
		})
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
