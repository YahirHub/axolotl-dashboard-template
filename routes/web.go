package routes

import (
	httpmiddleware "github.com/goravel/framework/http/middleware"

	"github.com/YahirHub/axolotl-dashboard-template/app/facades"
	"github.com/YahirHub/axolotl-dashboard-template/app/http/controllers"
	appmiddleware "github.com/YahirHub/axolotl-dashboard-template/app/http/middleware"
)

func Web() {
	authController := controllers.NewAuthController()
	dashboardController := controllers.NewDashboardController()
	pageController := controllers.NewPageController()
	fileController := controllers.NewFileController()
	healthController := controllers.NewHealthController()
	requireAuth := appmiddleware.NewRequireAuth()

	facades.Route().Get("/", authController.Home).Name("home")

	facades.Route().Get("/setup", authController.Setup).Name("setup")
	facades.Route().Middleware(httpmiddleware.Throttle("auth-setup")).
		Post("/setup", authController.StoreSetup).Name("setup.store")

	facades.Route().Get("/login", authController.Login).Name("login")
	facades.Route().Middleware(httpmiddleware.Throttle("auth-login")).
		Post("/login", authController.StoreLogin).Name("login.store")

	facades.Route().Get("/register", authController.Register).Name("register")
	facades.Route().Middleware(httpmiddleware.Throttle("auth-register")).
		Post("/register", authController.StoreRegister).Name("register.store")

	facades.Route().Get("/forgot-password", authController.ForgotPassword).Name("password.forgot")
	facades.Route().Middleware(httpmiddleware.Throttle("auth-recovery")).
		Post("/forgot-password", authController.StoreForgotPassword).Name("password.email")
	facades.Route().Get("/reset-password", authController.ResetPassword).Name("password.reset")
	facades.Route().Middleware(httpmiddleware.Throttle("auth-recovery")).
		Post("/reset-password", authController.StoreResetPassword).Name("password.update")

	facades.Route().Middleware(requireAuth).Post("/logout", authController.Logout).Name("logout")

	facades.Route().Middleware(requireAuth).Get("/dashboard", dashboardController.Index).Name("dashboard")
	facades.Route().Middleware(requireAuth).Get("/cuenta", pageController.Account).Name("account")
	facades.Route().Middleware(requireAuth).Get("/servidor", pageController.Server).Name("server")
	facades.Route().Middleware(requireAuth).Get("/componentes", pageController.Components).Name("components")
	facades.Route().Middleware(requireAuth).Get("/archivos", fileController.Index).Name("files")
	facades.Route().Middleware(requireAuth).Get("/archivos/*path", fileController.Index).Name("files.path")
	facades.Route().Middleware(requireAuth).Get("/api/files", fileController.List).Name("files.list")
	facades.Route().Middleware(requireAuth).Get("/api/files/read", fileController.Read).Name("files.read")
	facades.Route().Middleware(requireAuth).Post("/api/files/create", fileController.Create).Name("files.create")
	facades.Route().Middleware(requireAuth).Post("/api/files/save", fileController.Save).Name("files.save")
	facades.Route().Middleware(requireAuth).Post("/api/files/delete", fileController.Delete).Name("files.delete")

	facades.Route().Get("/health", healthController.Show).Name("health")
	facades.Route().Static("assets", "./public")
}
