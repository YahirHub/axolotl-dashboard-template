package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"

	"github.com/YahirHub/axolotl-dashboard-template/app/facades"
)

type M20260928000100CreateAppSettingsTable struct{}

func (r *M20260928000100CreateAppSettingsTable) Signature() string {
	return "20260928000100_create_app_settings_table"
}

func (r *M20260928000100CreateAppSettingsTable) Up() error {
	if facades.Schema().HasTable("app_settings") {
		return nil
	}

	return facades.Schema().Create("app_settings", func(table schema.Blueprint) {
		table.ID()
		table.String("key")
		table.LongText("value")
		table.Timestamps()
		table.Unique("key")
	})
}

func (r *M20260928000100CreateAppSettingsTable) Down() error {
	return facades.Schema().DropIfExists("app_settings")
}
