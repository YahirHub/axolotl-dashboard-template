package routes

import (
	"github.com/YahirHub/axolotl-dashboard-template/app/facades"
	"github.com/YahirHub/axolotl-dashboard-template/app/http/controllers"
)

func Web() {
	dashboardController := controllers.NewDashboardController()
	pageController := controllers.NewPageController()
	healthController := controllers.NewHealthController()

	facades.Route().Get("/", dashboardController.Index).Name("home")
	facades.Route().Get("/dashboard", dashboardController.Index).Name("dashboard")
	facades.Route().Get("/cuenta", pageController.Account).Name("account")
	facades.Route().Get("/servidor", pageController.Server).Name("server")
	facades.Route().Get("/componentes", pageController.Components).Name("components")
	facades.Route().Get("/health", healthController.Show).Name("health")
	facades.Route().Static("assets", "./public")
}
