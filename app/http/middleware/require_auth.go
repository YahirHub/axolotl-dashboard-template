package middleware

import (
	stdhttp "net/http"

	"github.com/goravel/framework/contracts/http"

	"github.com/YahirHub/axolotl-dashboard-template/app/facades"
	"github.com/YahirHub/axolotl-dashboard-template/app/models"
	"github.com/YahirHub/axolotl-dashboard-template/app/security"
)

type RequireAuth struct{}

func NewRequireAuth() *RequireAuth {
	return &RequireAuth{}
}

func (r *RequireAuth) Handle(ctx http.Context) {
	if !security.AdminExists() {
		_ = ctx.Response().Redirect(stdhttp.StatusSeeOther, "/setup").Abort()
		return
	}

	if facades.Auth(ctx).Guest() {
		_ = ctx.Response().Redirect(stdhttp.StatusSeeOther, "/login").Abort()
		return
	}

	var user models.User
	if err := facades.Auth(ctx).User(&user); err != nil || !user.Active {
		_ = facades.Auth(ctx).Logout()
		_ = ctx.Response().Redirect(stdhttp.StatusSeeOther, "/login").Abort()
		return
	}

	ctx.WithValue("user", user)
	ctx.Request().Next()
}

func (r *RequireAuth) Signature() string {
	return "axolotl.require-auth"
}
