package controllers

import "github.com/goravel/framework/contracts/http"

type PageMeta struct {
	Title        string
	PageTitle    string
	PageSubtitle string
	ActiveRoute  string
}

type PageController struct{}

func NewPageController() *PageController {
	return &PageController{}
}

func (r *PageController) Account(ctx http.Context) http.Response {
	return ctx.Response().View().Make("account.tmpl", PageMeta{
		Title:        "Cuenta — Axolotl",
		PageTitle:    "Cuenta",
		PageSubtitle: "Perfil, preferencias, apariencia y seguridad de la aplicación.",
		ActiveRoute:  "account",
	})
}

func (r *PageController) Server(ctx http.Context) http.Response {
	return ctx.Response().View().Make("server.tmpl", PageMeta{
		Title:        "Servidor — Axolotl",
		PageTitle:    "Servidor",
		PageSubtitle: "CPU, RAM, red, almacenamiento y estado técnico en una página dedicada.",
		ActiveRoute:  "server",
	})
}

func (r *PageController) Components(ctx http.Context) http.Response {
	return ctx.Response().View().Make("components.tmpl", PageMeta{
		Title:        "Componentes — Axolotl",
		PageTitle:    "Componentes UI",
		PageSubtitle: "Laboratorio de componentes, estados y patrones reutilizables.",
		ActiveRoute:  "components",
	})
}
