package main

import (
	"github.com/YahirHub/axolotl-dashboard-template/app/facades"
	"github.com/YahirHub/axolotl-dashboard-template/app/security"
	"github.com/YahirHub/axolotl-dashboard-template/bootstrap"
)

func main() {
	app := bootstrap.Boot()

	if _, err := security.InitializeAdminPIN(); err != nil {
		facades.Log().Errorf(
			"no se pudo inicializar el setup de administrador: %v. Si la base es nueva ejecuta: go run . artisan migrate",
			err,
		)
		return
	}

	app.Start()
}
