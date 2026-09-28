package controllers

import (
	"fmt"
	stdhttp "net/http"
	"net/mail"
	"strings"
	"time"
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
	Success             string
	OldName             string
	OldEmail            string
	OldAlias            string
	ResetToken          string
	RegistrationEnabled bool
	RecoveryEnabled     bool
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
			Name:        name,
			Email:       email,
			Alias:       alias,
			Password:    hashed,
			Role:        models.RoleAdmin,
			Active:      true,
			AuthVersion: 1,
		}
		return facades.Orm().Query().Create(&user)
	})
	if err != nil {
		facades.Log().Channel("security").Errorf("setup inicial falló: %v", err)

		message := "No fue posible completar el setup."
		if strings.Contains(err.Error(), "pin") {
			message = "PIN temporal incorrecto."
		}
		if strings.Contains(err.Error(), "correo") {
			message = err.Error()
		}
		return r.authViewWithOld(ctx, "setup", message, name, email, alias)
	}

	if err := startAuthenticatedSession(ctx, &user); err != nil {
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
	ip := security.ClientAddress(ctx)
	limiter := security.LoginAttempts()

	if retry, blocked := limiter.Check(ip, email); blocked {
		return r.blockedLogin(ctx, email, ip, retry)
	}

	fail := func() http.Response {
		retry := limiter.RecordFailure(ip, email)
		if retry > 0 {
			return r.blockedLogin(ctx, email, ip, retry)
		}
		return r.authViewWithOld(ctx, "login", "Credenciales incorrectas.", "", email, "")
	}

	if email == "" || len(email) > 190 || password == "" || len(password) > 128 {
		_, _ = facades.Hash().Make("AxolotlDummyPassword123!")
		return fail()
	}

	var user models.User
	if err := facades.Orm().Query().Where("email", email).Where("active", true).First(&user); err != nil {
		// Ejecuta Argon2id aunque el correo no exista para reducir diferencias
		// de tiempo que podrían facilitar enumeración de usuarios.
		_, _ = facades.Hash().Make("AxolotlDummyPassword123!")
		return fail()
	}
	if !facades.Hash().Check(password, user.Password) {
		return fail()
	}

	if facades.Hash().NeedsRehash(user.Password) {
		if hashed, err := facades.Hash().Make(password); err == nil {
			user.Password = hashed
			_ = facades.Orm().Query().Save(&user)
		}
	}

	if err := startAuthenticatedSession(ctx, &user); err != nil {
		facades.Log().Errorf("error iniciando sesión: %v", err)
		return r.authViewWithOld(ctx, "login", "No fue posible iniciar sesión.", "", email, "")
	}

	limiter.RecordSuccess(ip, email)
	return ctx.Response().Redirect(stdhttp.StatusSeeOther, "/dashboard")
}

func (r *AuthController) blockedLogin(ctx http.Context, email, ip string, retry time.Duration) http.Response {
	seconds := int(retry.Seconds())
	if seconds < 1 {
		seconds = 1
	}
	ctx.Response().Header("Retry-After", fmt.Sprint(seconds))
	facades.Log().Channel("security").Warningf(
		"login bloqueado ip=%s identity=%s retry_seconds=%d",
		ip,
		security.IdentityFingerprint(email),
		seconds,
	)
	return r.authViewWithOld(
		ctx,
		"login",
		"Demasiados intentos. Espera un momento antes de volver a intentarlo.",
		"",
		email,
		"",
	)
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
		Name:        name,
		Email:       email,
		Alias:       alias,
		Password:    hashed,
		Role:        models.RoleUser,
		Active:      true,
		AuthVersion: 1,
	}
	if err := facades.Orm().Query().Create(&user); err != nil {
		return r.authViewWithOld(ctx, "register", "No fue posible crear la cuenta.", name, email, alias)
	}

	if err := startAuthenticatedSession(ctx, &user); err != nil {
		return ctx.Response().Redirect(stdhttp.StatusSeeOther, "/login")
	}

	return ctx.Response().Redirect(stdhttp.StatusSeeOther, "/dashboard")
}

func (r *AuthController) ForgotPassword(ctx http.Context) http.Response {
	if !security.AdminExists() {
		return ctx.Response().Redirect(stdhttp.StatusSeeOther, "/setup")
	}
	if !security.PasswordRecoveryEnabled() {
		return ctx.Response().Redirect(stdhttp.StatusSeeOther, "/login")
	}
	return r.authView(ctx, "forgot", "")
}

func (r *AuthController) StoreForgotPassword(ctx http.Context) http.Response {
	if !security.AdminExists() {
		return ctx.Response().Redirect(stdhttp.StatusSeeOther, "/setup")
	}
	if !security.PasswordRecoveryEnabled() {
		return ctx.Response().Redirect(stdhttp.StatusSeeOther, "/login")
	}

	email := normalizeEmail(ctx.Request().Input("email"))
	genericMessage := "Si existe una cuenta con ese correo, recibirás un enlace para restablecer la contraseña."
	ctx.Request().Session().Flash("auth_success", genericMessage)

	if len(email) <= 190 && validEmail(email) {
		cacheKey := "password-recovery:" + security.IdentityFingerprint(email)
		if facades.Cache().Add(cacheKey, true, 2*time.Minute) {
			go r.sendPasswordRecovery(email)
		}
	}

	return ctx.Response().Redirect(stdhttp.StatusSeeOther, "/forgot-password")
}

func (r *AuthController) ResetPassword(ctx http.Context) http.Response {
	if !security.PasswordRecoveryEnabled() {
		return ctx.Response().Redirect(stdhttp.StatusSeeOther, "/login")
	}

	return r.authViewWithToken(ctx, "reset", "", "")
}

func (r *AuthController) StoreResetPassword(ctx http.Context) http.Response {
	if !security.PasswordRecoveryEnabled() {
		return ctx.Response().Redirect(stdhttp.StatusSeeOther, "/login")
	}

	token := strings.TrimSpace(ctx.Request().Input("token"))
	password := ctx.Request().Input("password")
	confirmation := ctx.Request().Input("password_confirmation")
	if message := validatePassword(password, confirmation); message != "" {
		return r.authViewWithToken(ctx, "reset", message, token)
	}

	hashed, err := facades.Hash().Make(password)
	if err != nil {
		return r.authViewWithToken(ctx, "reset", "No fue posible proteger la nueva contraseña.", token)
	}

	manager := security.NewPasswordResetManager()
	user, err := manager.Consume(token, hashed)
	if err != nil {
		return r.authViewWithToken(ctx, "reset", "El enlace de recuperación no es válido o ya expiró.", token)
	}

	security.LoginAttempts().ClearIdentity(user.Email)
	_ = facades.Auth(ctx).Logout()
	ctx.Request().Session().Flash("auth_success", "Contraseña actualizada. Ya puedes iniciar sesión con tu nueva contraseña.")
	facades.Log().Channel("security").Warningf(
		"password reset completado identity=%s",
		security.IdentityFingerprint(user.Email),
	)

	return ctx.Response().Redirect(stdhttp.StatusSeeOther, "/login")
}

func (r *AuthController) sendPasswordRecovery(email string) {
	var user models.User
	if err := facades.Orm().Query().
		Where("email", email).
		Where("active", true).
		First(&user); err != nil {
		return
	}

	manager := security.NewPasswordResetManager()
	ttlMinutes := facades.Config().GetInt("auth.password_reset.ttl_minutes", 30)
	token, err := manager.Issue(user.ID, time.Duration(ttlMinutes)*time.Minute)
	if err != nil {
		facades.Log().Channel("security").Errorf(
			"no se pudo emitir token de recovery identity=%s",
			security.IdentityFingerprint(email),
		)
		return
	}

	link, err := security.PasswordResetLink(token)
	if err != nil {
		facades.Log().Channel("security").Errorf(
			"no se pudo construir enlace de recovery identity=%s",
			security.IdentityFingerprint(email),
		)
		return
	}

	if err := security.SendPasswordResetEmail(user.Email, user.Name, link); err != nil {
		facades.Log().Channel("security").Errorf(
			"no se pudo enviar recovery identity=%s error=%v",
			security.IdentityFingerprint(email),
			err,
		)
	}
}

func (r *AuthController) Logout(ctx http.Context) http.Response {
	_ = facades.Auth(ctx).Logout()
	return ctx.Response().Redirect(stdhttp.StatusSeeOther, "/login")
}

func (r *AuthController) authView(ctx http.Context, mode, errorMessage string) http.Response {
	return r.renderAuthView(ctx, mode, errorMessage, "", "", "", "")
}

func (r *AuthController) authViewWithOld(ctx http.Context, mode, errorMessage, name, email, alias string) http.Response {
	return r.renderAuthView(ctx, mode, errorMessage, name, email, alias, "")
}

func (r *AuthController) authViewWithToken(ctx http.Context, mode, errorMessage, token string) http.Response {
	return r.renderAuthView(ctx, mode, errorMessage, "", "", "", token)
}

func (r *AuthController) renderAuthView(ctx http.Context, mode, errorMessage, name, email, alias, token string) http.Response {
	ctx.Response().
		Header("Cache-Control", "no-store, max-age=0").
		Header("Pragma", "no-cache").
		Header("Referrer-Policy", "no-referrer")

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
	case "forgot":
		title = "Recuperar cuenta — Axolotl"
		heading = "Recuperar cuenta"
		subtitle = "Solicita un enlace temporal para restablecer tu contraseña."
	case "reset":
		title = "Nueva contraseña — Axolotl"
		heading = "Crear nueva contraseña"
		subtitle = "El enlace es temporal y solo puede utilizarse una vez."
	}

	success := ""
	if value := ctx.Request().Session().Pull("auth_success", nil); value != nil {
		success = fmt.Sprint(value)
	}

	return ctx.Response().View().Make("auth.tmpl", AuthPageData{
		Title:               title,
		Mode:                mode,
		Heading:             heading,
		Subtitle:            subtitle,
		CSRFToken:           ctx.Request().Session().Token(),
		Error:               errorMessage,
		Success:             success,
		OldName:             name,
		OldEmail:            email,
		OldAlias:            alias,
		ResetToken:          token,
		RegistrationEnabled: facades.Config().GetBool("auth.registration.enabled"),
		RecoveryEnabled:     security.PasswordRecoveryEnabled(),
	})
}

func startAuthenticatedSession(ctx http.Context, user *models.User) error {
	if user.AuthVersion == 0 {
		user.AuthVersion = 1
	}
	if _, err := facades.Auth(ctx).Login(user); err != nil {
		return err
	}
	ctx.Request().Session().Put("auth_version", fmt.Sprint(user.AuthVersion))
	return ctx.Request().Session().Save()
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
	if message := validatePassword(password, confirmation); message != "" {
		return name, email, alias, password, message
	}

	return name, email, alias, password, ""
}

func validatePassword(password, confirmation string) string {
	if len(password) < 10 || len(password) > 128 {
		return "La contraseña debe tener entre 10 y 128 caracteres."
	}
	if password != confirmation {
		return "Las contraseñas no coinciden."
	}
	return ""
}

func normalizeEmail(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}

func validEmail(value string) bool {
	address, err := mail.ParseAddress(value)
	return err == nil && strings.EqualFold(address.Address, value)
}
