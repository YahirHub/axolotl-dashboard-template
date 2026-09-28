package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"

	"github.com/YahirHub/axolotl-dashboard-template/app/facades"
)

type M20260928000300AddAuthVersionAndPasswordResetTokens struct{}

func (r *M20260928000300AddAuthVersionAndPasswordResetTokens) Signature() string {
	return "20260928000300_add_auth_version_and_password_reset_tokens"
}

func (r *M20260928000300AddAuthVersionAndPasswordResetTokens) Up() error {
	if facades.Schema().HasTable("users") && !facades.Schema().HasColumn("users", "auth_version") {
		if err := facades.Schema().Table("users", func(table schema.Blueprint) {
			table.UnsignedInteger("auth_version").Default(1)
		}); err != nil {
			return err
		}
	}

	if facades.Schema().HasTable("password_reset_tokens") {
		return nil
	}

	return facades.Schema().Create("password_reset_tokens", func(table schema.Blueprint) {
		table.ID()
		table.UnsignedBigInteger("user_id")
		table.String("token_hash", 64)
		table.DateTime("expires_at")
		table.Boolean("used").Default(false)
		table.Timestamps()
		table.Unique("token_hash")
		table.Index("user_id")
		table.Index("expires_at")
	})
}

func (r *M20260928000300AddAuthVersionAndPasswordResetTokens) Down() error {
	if err := facades.Schema().DropIfExists("password_reset_tokens"); err != nil {
		return err
	}
	if facades.Schema().HasTable("users") && facades.Schema().HasColumn("users", "auth_version") {
		return facades.Schema().DropColumns("users", []string{"auth_version"})
	}
	return nil
}
