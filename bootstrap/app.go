package bootstrap

import (
	contractsfoundation "github.com/goravel/framework/contracts/foundation"
	"github.com/goravel/framework/contracts/foundation/configuration"
	"github.com/goravel/framework/foundation"

	appmiddleware "github.com/YahirHub/axolotl-dashboard-template/app/http/middleware"
	"github.com/YahirHub/axolotl-dashboard-template/config"
	"github.com/YahirHub/axolotl-dashboard-template/routes"
)

func Boot() contractsfoundation.Application {
	return foundation.Setup().
		WithMiddleware(func(handler configuration.Middleware) {
			handler.Append(
				appmiddleware.NewSecurityHeaders(),
			)
		}).
		WithMigrations(Migrations).
		WithRouting(func() {
			routes.Web()
			routes.Grpc()
		}).
		WithProviders(Providers).
		WithConfig(config.Boot).
		Create()
}
