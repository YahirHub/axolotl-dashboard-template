package middleware

import (
	"fmt"
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
		facades.Log().Channel("security").Warningf(
			"acceso privado sin sesión ip=%s path=%s",
			security.ClientAddress(ctx),
			ctx.Request().Path(),
		)
		_ = ctx.Response().Redirect(stdhttp.StatusSeeOther, "/login").Abort()
		return
	}

	var user models.User
	if err := facades.Auth(ctx).User(&user); err != nil || !user.Active {
		_ = facades.Auth(ctx).Logout()
		_ = ctx.Response().Redirect(stdhttp.StatusSeeOther, "/login").Abort()
		return
	}

	sessionVersion := fmt.Sprint(ctx.Request().Session().Get("auth_version", ""))
	if user.AuthVersion == 0 || sessionVersion != fmt.Sprint(user.AuthVersion) {
		facades.Log().Channel("security").Warningf(
			"sesión invalidada por versión ip=%s user_id=%d",
			security.ClientAddress(ctx),
			user.ID,
		)
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
