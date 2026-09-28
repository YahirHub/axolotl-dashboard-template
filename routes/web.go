package routes

import (
	"github.com/YahirHub/axolotl-dashboard-template/app/facades"
	"github.com/YahirHub/axolotl-dashboard-template/app/http/controllers"
)

func Web() {
	dashboardController := controllers.NewDashboardController()
	healthController := controllers.NewHealthController()

	facades.Route().Get("/", dashboardController.Index).Name("dashboard")
	facades.Route().Get("/health", healthController.Show).Name("health")
	facades.Route().Static("assets", "./public")
}
