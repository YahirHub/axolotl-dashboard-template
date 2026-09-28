package bootstrap

import (
	contractsfoundation "github.com/goravel/framework/contracts/foundation"
	"github.com/goravel/framework/contracts/foundation/configuration"
	contractshttp "github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/foundation"
	"github.com/goravel/framework/http/limit"
	httpmiddleware "github.com/goravel/framework/http/middleware"
	sessionmiddleware "github.com/goravel/framework/session/middleware"

	appfacades "github.com/YahirHub/axolotl-dashboard-template/app/facades"
	appmiddleware "github.com/YahirHub/axolotl-dashboard-template/app/http/middleware"
	"github.com/YahirHub/axolotl-dashboard-template/config"
	"github.com/YahirHub/axolotl-dashboard-template/routes"
)

func Boot() contractsfoundation.Application {
	return foundation.Setup().
		WithMiddleware(func(handler configuration.Middleware) {
			handler.Append(
				appmiddleware.NewSecurityHeaders(),
				sessionmiddleware.StartSession(),
				httpmiddleware.VerifyCsrfToken(),
			)
		}).
		WithMigrations(Migrations).
		WithRouting(func() {
			routes.Web()
			routes.Grpc()
		}).
		WithProviders(Providers).
		WithConfig(config.Boot).
		WithCallback(func() {
			appfacades.RateLimiter().For("auth-login", func(ctx contractshttp.Context) contractshttp.Limit {
				return limit.PerMinute(10).By(ctx.Request().Ip())
			})
			appfacades.RateLimiter().For("auth-register", func(ctx contractshttp.Context) contractshttp.Limit {
				return limit.PerMinute(5).By(ctx.Request().Ip())
			})
			appfacades.RateLimiter().For("auth-setup", func(ctx contractshttp.Context) contractshttp.Limit {
				return limit.PerMinute(5).By(ctx.Request().Ip())
			})
		}).
		Create()
}
