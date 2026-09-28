package controllers

import (
	"strings"

	"github.com/goravel/framework/contracts/http"

	"github.com/YahirHub/axolotl-dashboard-template/app/facades"
	"github.com/YahirHub/axolotl-dashboard-template/app/models"
)

type CurrentUserView struct {
	Name    string
	Email   string
	Alias   string
	Role    string
	IsAdmin bool
}

type PageMeta struct {
	Title        string
	PageTitle    string
	PageSubtitle string
	ActiveRoute  string
	CSRFToken    string
	CurrentUser  CurrentUserView
}

func pageMeta(ctx http.Context, title, pageTitle, pageSubtitle, activeRoute string) PageMeta {
	meta := PageMeta{
		Title:        title,
		PageTitle:    pageTitle,
		PageSubtitle: pageSubtitle,
		ActiveRoute:  activeRoute,
	}

	if ctx.Request().HasSession() {
		meta.CSRFToken = ctx.Request().Session().Token()
	}

	var user models.User
	if err := facades.Auth(ctx).User(&user); err == nil {
		role := "Usuario"
		if user.IsAdmin() {
			role = "Administrador"
		}
		meta.CurrentUser = CurrentUserView{
			Name:    user.Name,
			Email:   user.Email,
			Alias:   strings.ToUpper(user.Alias),
			Role:    role,
			IsAdmin: user.IsAdmin(),
		}
	}

	return meta
}

type PageController struct{}

func NewPageController() *PageController {
	return &PageController{}
}

func (r *PageController) Account(ctx http.Context) http.Response {
	return ctx.Response().View().Make("account.tmpl", pageMeta(
		ctx,
		"Cuenta — Axolotl",
		"Cuenta",
		"Perfil, preferencias, apariencia y seguridad de la aplicación.",
		"account",
	))
}

func (r *PageController) Server(ctx http.Context) http.Response {
	return ctx.Response().View().Make("server.tmpl", pageMeta(
		ctx,
		"Servidor — Axolotl",
		"Servidor",
		"CPU, RAM, red, almacenamiento y estado técnico en una página dedicada.",
		"server",
	))
}

func (r *PageController) Components(ctx http.Context) http.Response {
	return ctx.Response().View().Make("components.tmpl", pageMeta(
		ctx,
		"Componentes — Axolotl",
		"Componentes UI",
		"Laboratorio de componentes, estados y patrones reutilizables.",
		"components",
	))
}
