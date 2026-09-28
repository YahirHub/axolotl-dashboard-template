package feature

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/YahirHub/axolotl-dashboard-template/tests"
)

func TestDashboardRoute(t *testing.T) {
	var testCase tests.TestCase

	response, err := testCase.Http(t).Get("/")
	require.NoError(t, err)

	response.
		AssertOk().
		AssertHeader("X-Content-Type-Options", "nosniff").
		AssertHeader("X-Frame-Options", "SAMEORIGIN").
		AssertSee([]string{
			"Axolotl",
			"Centro de operaciones",
			"Proyectos destacados",
			"Tareas prioritarias",
			"Equipo",
			"Agenda",
			"Integraciones",
			"Métricas del servidor",
			"Laboratorio UI",
			"/assets/css/app.css",
			"/assets/js/app.js",
			"/assets/brand/axolotl.svg",
			"/assets/brand/axolotl.png",
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
