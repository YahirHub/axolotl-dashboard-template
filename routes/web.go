package routes

import (
	"github.com/YahirHub/axolotl-dashboard-template/app/facades"
	"github.com/YahirHub/axolotl-dashboard-template/app/http/controllers"
)

func Web() {
	dashboardController := controllers.NewDashboardController()

	facades.Route().Get("/", dashboardController.Index)
	facades.Route().Static("assets", "./public")
}
