package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"

	"github.com/YahirHub/axolotl-dashboard-template/app/facades"
)

type M20260928000200CreateUsersTable struct{}

func (r *M20260928000200CreateUsersTable) Signature() string {
	return "20260928000200_create_users_table"
}

func (r *M20260928000200CreateUsersTable) Up() error {
	if facades.Schema().HasTable("users") {
		return nil
	}

	return facades.Schema().Create("users", func(table schema.Blueprint) {
		table.ID()
		table.String("name", 100)
		table.String("email", 190)
		table.String("alias", 3)
		table.String("password", 255)
		table.String("role", 32)
		table.Boolean("active")
		table.UnsignedInteger("auth_version").Default(1)
		table.Timestamps()
		table.Unique("email")
		table.Index("role")
	})
}

func (r *M20260928000200CreateUsersTable) Down() error {
	return facades.Schema().DropIfExists("users")
}
