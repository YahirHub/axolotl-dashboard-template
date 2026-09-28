package config

import (
	"fmt"
	"path/filepath"

	"github.com/goravel/framework/contracts/database/driver"
	"github.com/goravel/framework/support/path"
	sqlitefacades "github.com/goravel/sqlite/facades"

	"github.com/YahirHub/axolotl-dashboard-template/app/facades"
)

func init() {
	config := facades.Config()

	database := fmt.Sprint(config.Env("DB_DATABASE", "storage/axolotl.sqlite"))
	if !filepath.IsAbs(database) {
		database = path.Base(database)
	}

	config.Add("database", map[string]any{
		"default": config.Env("DB_CONNECTION", "sqlite"),
		"connections": map[string]any{
			"sqlite": map[string]any{
				"database": database,
				"prefix":   "",
				"singular": false,
				"via": func() (driver.Driver, error) {
					return sqlitefacades.Sqlite("sqlite")
				},
			},
		},
		"pool": map[string]any{
			"max_idle_conns":    1,
			"max_open_conns":    1,
			"conn_max_idletime": 3600,
			"conn_max_lifetime": 3600,
		},
		"slow_threshold": 200,
		"migrations": map[string]any{
			"table": "migrations",
		},
	})
}
