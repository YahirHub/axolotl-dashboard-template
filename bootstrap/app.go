package bootstrap

import (
	contractsfoundation "github.com/goravel/framework/contracts/foundation"
	"github.com/goravel/framework/foundation"

	"github.com/YahirHub/axolotl-dashboard-template/config"
	"github.com/YahirHub/axolotl-dashboard-template/routes"
)

func Boot() contractsfoundation.Application {
	return foundation.Setup().
		WithMigrations(Migrations).
		WithRouting(func() {
			routes.Web()
			routes.Grpc()
		}).
		WithProviders(Providers).
		WithConfig(config.Boot).
		Create()
}
