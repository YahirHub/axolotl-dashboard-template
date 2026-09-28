package controllers

import (
	"fmt"
	stdhttp "net/http"
	"net/mail"
	"strings"
	"unicode"

	"github.com/goravel/framework/contracts/http"

	"github.com/YahirHub/axolotl-dashboard-template/app/facades"
	"github.com/YahirHub/axolotl-dashboard-template/app/models"
	"github.com/YahirHub/axolotl-dashboard-template/app/security"
)

type AuthPageData struct {
	Title               string
	Mode                string
	Heading             string
	Subtitle            string
	CSRFToken           string
	Error               string
	OldName             string
	OldEmail            string
	OldAlias            string
	RegistrationEnabled bool
}

type AuthController struct{}

func NewAuthController() *AuthController {
	return &AuthController{}
}

func (r *AuthController) Home(ctx http.Context) http.Response {
	if !security.AdminExists() {
		return ctx.Response().Redirect(stdhttp.StatusSeeOther, "/setup")
	}
	if facades.Auth(ctx).Check() {
		return ctx.Response().Redirect(stdhttp.StatusSeeOther, "/dashboard")
	}
	return ctx.Response().Redirect(stdhttp.StatusSeeOther, "/login")
}

func (r *AuthController) Setup(ctx http.Context) http.Response {
	if security.AdminExists() {
		if facades.Auth(ctx).Check() {
			return ctx.Response().Redirect(stdhttp.StatusSeeOther, "/dashboard")
		}
		return ctx.Response().Redirect(stdhttp.StatusSeeOther, "/login")
	}
	if _, err := security.InitializeAdminPIN(); err != nil {
		facades.Log().Errorf("no se pudo asegurar el PIN de setup: %v", err)
		return r.authView(ctx, "setup", "No fue posible preparar el PIN temporal. Revisa el log del servidor.")
	}
	return r.authView(ctx, "setup", "")
}

func (r *AuthController) StoreSetup(ctx http.Context) http.Response {
	if security.AdminExists() {
		return ctx.Response().Redirect(stdhttp.StatusSeeOther, "/login")
	}

	name, email, alias, password, errMessage := readAccountInput(ctx)
	if errMessage != "" {
		return r.authViewWithOld(ctx, "setup", errMessage, name, email, alias)
	}

	pin := ""
	for i := 1; i <= 6; i++ {
		pin += ctx.Request().Input(fmt.Sprintf("pin_%d", i))
	}
	if len(pin) != 6 || strings.IndexFunc(pin, func(r rune) bool { return !unicode.IsDigit(r) }) >= 0 {
		return r.authViewWithOld(ctx, "setup", "El PIN debe contener exactamente 6 dígitos.", name, email, alias)
	}

	var user models.User
	err := security.CompleteWithPIN(pin, func() error {
		exists, err := facades.Orm().Query().Model(&models.User{}).Where("email", email).Exists()
		if err != nil {
			return err
		}
		if exists {
			return fmt.Errorf("el correo ya está registrado")
		}

		hashed, err := facades.Hash().Make(password)
		if err != nil {
			return err
		}

		user = models.User{
			Name:     name,
			Email:    email,
			Alias:    alias,
			Password: hashed,
			Role:     models.RoleAdmin,
			Active:   true,
		}
		return facades.Orm().Query().Create(&user)
	})
	if err != nil {
		message := "No fue posible completar el setup."
		if strings.Contains(err.Error(), "pin") {
			message = "PIN temporal incorrecto."
		}
		if strings.Contains(err.Error(), "correo") {
			message = err.Error()
		}
		return r.authViewWithOld(ctx, "setup", message, name, email, alias)
	}

	if _, err := facades.Auth(ctx).Login(&user); err != nil {
		facades.Log().Errorf("admin creado pero no se pudo iniciar sesión: %v", err)
		return ctx.Response().Redirect(stdhttp.StatusSeeOther, "/login")
	}

	facades.Log().Infof("Setup inicial completado. Administrador creado: %s", user.Email)
	return ctx.Response().Redirect(stdhttp.StatusSeeOther, "/dashboard")
}

func (r *AuthController) Login(ctx http.Context) http.Response {
	if !security.AdminExists() {
		return ctx.Response().Redirect(stdhttp.StatusSeeOther, "/setup")
	}
	if facades.Auth(ctx).Check() {
		return ctx.Response().Redirect(stdhttp.StatusSeeOther, "/dashboard")
	}
	return r.authView(ctx, "login", "")
}

func (r *AuthController) StoreLogin(ctx http.Context) http.Response {
	if !security.AdminExists() {
		return ctx.Response().Redirect(stdhttp.StatusSeeOther, "/setup")
	}

	email := normalizeEmail(ctx.Request().Input("email"))
	password := ctx.Request().Input("password")
	if email == "" || len(email) > 190 || password == "" || len(password) > 128 {
		return r.authViewWithOld(ctx, "login", "Credenciales incorrectas.", "", email, "")
	}

	var user models.User
	if err := facades.Orm().Query().Where("email", email).Where("active", true).First(&user); err != nil {
		// Ejecuta una operación Argon2id incluso cuando el correo no existe
		// para reducir diferencias de tiempo que permitan enumerar usuarios.
		_, _ = facades.Hash().Make("AxolotlDummyPassword123!")
		return r.authViewWithOld(ctx, "login", "Credenciales incorrectas.", "", email, "")
	}
	if !facades.Hash().Check(password, user.Password) {
		return r.authViewWithOld(ctx, "login", "Credenciales incorrectas.", "", email, "")
	}

	if facades.Hash().NeedsRehash(user.Password) {
		if hashed, err := facades.Hash().Make(password); err == nil {
			user.Password = hashed
			_ = facades.Orm().Query().Save(&user)
		}
	}

	if _, err := facades.Auth(ctx).Login(&user); err != nil {
		facades.Log().Errorf("error iniciando sesión para %s: %v", email, err)
		return r.authViewWithOld(ctx, "login", "No fue posible iniciar sesión.", "", email, "")
	}

	return ctx.Response().Redirect(stdhttp.StatusSeeOther, "/dashboard")
}

func (r *AuthController) Register(ctx http.Context) http.Response {
	if !security.AdminExists() {
		return ctx.Response().Redirect(stdhttp.StatusSeeOther, "/setup")
	}
	if !facades.Config().GetBool("auth.registration.enabled") {
		return ctx.Response().Redirect(stdhttp.StatusSeeOther, "/login")
	}
	if facades.Auth(ctx).Check() {
		return ctx.Response().Redirect(stdhttp.StatusSeeOther, "/dashboard")
	}
	return r.authView(ctx, "register", "")
}

func (r *AuthController) StoreRegister(ctx http.Context) http.Response {
	if !security.AdminExists() {
		return ctx.Response().Redirect(stdhttp.StatusSeeOther, "/setup")
	}
	if !facades.Config().GetBool("auth.registration.enabled") {
		return ctx.Response().Redirect(stdhttp.StatusSeeOther, "/login")
	}

	name, email, alias, password, errMessage := readAccountInput(ctx)
	if errMessage != "" {
		return r.authViewWithOld(ctx, "register", errMessage, name, email, alias)
	}

	exists, err := facades.Orm().Query().Model(&models.User{}).Where("email", email).Exists()
	if err != nil {
		return r.authViewWithOld(ctx, "register", "No fue posible validar el correo.", name, email, alias)
	}
	if exists {
		return r.authViewWithOld(ctx, "register", "Ese correo ya está registrado.", name, email, alias)
	}

	hashed, err := facades.Hash().Make(password)
	if err != nil {
		return r.authViewWithOld(ctx, "register", "No fue posible proteger la contraseña.", name, email, alias)
	}

	user := models.User{
		Name:     name,
		Email:    email,
		Alias:    alias,
		Password: hashed,
		Role:     models.RoleUser,
		Active:   true,
	}
	if err := facades.Orm().Query().Create(&user); err != nil {
		return r.authViewWithOld(ctx, "register", "No fue posible crear la cuenta.", name, email, alias)
	}

	if _, err := facades.Auth(ctx).Login(&user); err != nil {
		return ctx.Response().Redirect(stdhttp.StatusSeeOther, "/login")
	}

	return ctx.Response().Redirect(stdhttp.StatusSeeOther, "/dashboard")
}

func (r *AuthController) Logout(ctx http.Context) http.Response {
	_ = facades.Auth(ctx).Logout()
	return ctx.Response().Redirect(stdhttp.StatusSeeOther, "/login")
}

func (r *AuthController) authView(ctx http.Context, mode, errorMessage string) http.Response {
	return r.authViewWithOld(ctx, mode, errorMessage, "", "", "")
}

func (r *AuthController) authViewWithOld(ctx http.Context, mode, errorMessage, name, email, alias string) http.Response {
	title, heading, subtitle := "Acceso — Axolotl", "Bienvenido", "Inicia sesión para continuar."
	switch mode {
	case "setup":
		title = "Setup inicial — Axolotl"
		heading = "Crear administrador"
		subtitle = "Configura la primera cuenta administrativa usando el PIN temporal mostrado en el log del servidor."
	case "register":
		title = "Registro — Axolotl"
		heading = "Crear cuenta"
		subtitle = "Registra un usuario estándar para esta instalación de Axolotl."
	}

	return ctx.Response().View().Make("auth.tmpl", AuthPageData{
		Title:               title,
		Mode:                mode,
		Heading:             heading,
		Subtitle:            subtitle,
		CSRFToken:           ctx.Request().Session().Token(),
		Error:               errorMessage,
		OldName:             name,
		OldEmail:            email,
		OldAlias:            alias,
		RegistrationEnabled: facades.Config().GetBool("auth.registration.enabled"),
	})
}

func readAccountInput(ctx http.Context) (name, email, alias, password, errorMessage string) {
	name = strings.TrimSpace(ctx.Request().Input("name"))
	email = normalizeEmail(ctx.Request().Input("email"))
	alias = strings.ToUpper(strings.TrimSpace(ctx.Request().Input("alias")))
	password = ctx.Request().Input("password")
	confirmation := ctx.Request().Input("password_confirmation")

	if len(name) < 2 || len(name) > 100 {
		return name, email, alias, password, "El nombre debe tener entre 2 y 100 caracteres."
	}
	if len(email) > 190 || !validEmail(email) {
		return name, email, alias, password, "Ingresa un correo válido."
	}
	if len(alias) < 2 || len(alias) > 3 || strings.IndexFunc(alias, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r)
	}) >= 0 {
		return name, email, alias, password, "El alias debe tener 2 o 3 letras o números."
	}
	if len(password) < 10 || len(password) > 128 {
		return name, email, alias, password, "La contraseña debe tener entre 10 y 128 caracteres."
	}
	if password != confirmation {
		return name, email, alias, password, "Las contraseñas no coinciden."
	}

	return name, email, alias, password, ""
}

func normalizeEmail(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}

func validEmail(value string) bool {
	address, err := mail.ParseAddress(value)
	return err == nil && strings.EqualFold(address.Address, value)
}
