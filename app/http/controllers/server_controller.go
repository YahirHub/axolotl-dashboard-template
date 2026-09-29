package controllers

import (
	"errors"
	"fmt"
	"io/fs"
	stdhttp "net/http"
	"os"
	"runtime"
	"runtime/debug"
	"strings"
	"time"

	"github.com/goravel/framework/contracts/http"

	"github.com/YahirHub/axolotl-dashboard-template/app/facades"
	appserver "github.com/YahirHub/axolotl-dashboard-template/app/server"
)

type ServerPageData struct {
	Title        string
	PageTitle    string
	PageSubtitle string
	ActiveRoute  string
	CSRFToken    string
	CurrentUser  CurrentUserView
}

type ServerEnvironmentPageData struct {
	Title          string
	PageTitle      string
	PageSubtitle   string
	ActiveRoute    string
	CSRFToken      string
	CurrentUser    CurrentUserView
	GoVersion      string
	GOOS           string
	GOARCH         string
	CPUs           int
	GOMAXPROCS     int
	PID            int
	Uptime         string
	AppEnvironment string
	AppDebug       bool
	Database       string
	SessionDriver  string
	ModuleVersion  string
}

type ServerController struct {
	tasks     *appserver.TaskStore
	startedAt time.Time
}

func NewServerController() *ServerController {
	store, err := appserver.NewTaskStore(facades.Config().GetString("server.tasks_root"))
	if err != nil {
		facades.Log().Errorf("no se pudo configurar plantillas de tareas: %v", err)
		return &ServerController{startedAt: time.Now()}
	}
	if err := store.Ensure(); err != nil {
		facades.Log().Errorf("no se pudo preparar storage/server-tasks: %v", err)
		return &ServerController{startedAt: time.Now()}
	}

	return &ServerController{
		tasks:     store,
		startedAt: time.Now(),
	}
}

func (r *ServerController) Overview(ctx http.Context) http.Response {
	meta := pageMeta(
		ctx,
		"Servidor — Axolotl",
		"Servidor",
		"CPU, RAM, red, almacenamiento y estado técnico en una página dedicada.",
		"server",
	)
	return ctx.Response().View().Make("server.tmpl", serverPageData(meta))
}

func (r *ServerController) Tasks(ctx http.Context) http.Response {
	meta := pageMeta(
		ctx,
		"Administrador de tareas — Axolotl",
		"Administrador de tareas",
		"Crea y revisa plantillas systemd sin instalarlas ni ejecutarlas en el host.",
		"server-tasks",
	)
	return ctx.Response().View().Make("server_tasks.tmpl", serverPageData(meta))
}

func (r *ServerController) Environment(ctx http.Context) http.Response {
	meta := pageMeta(
		ctx,
		"Entorno del servidor — Axolotl",
		"Entorno",
		"Información segura del runtime y configuración técnica no sensible.",
		"server-environment",
	)

	moduleVersion := "devel"
	if info, ok := debug.ReadBuildInfo(); ok && info.Main.Version != "" {
		moduleVersion = info.Main.Version
	}

	return ctx.Response().View().Make("server_environment.tmpl", ServerEnvironmentPageData{
		Title:          meta.Title,
		PageTitle:      meta.PageTitle,
		PageSubtitle:   meta.PageSubtitle,
		ActiveRoute:    meta.ActiveRoute,
		CSRFToken:      meta.CSRFToken,
		CurrentUser:    meta.CurrentUser,
		GoVersion:      runtime.Version(),
		GOOS:           runtime.GOOS,
		GOARCH:         runtime.GOARCH,
		CPUs:           runtime.NumCPU(),
		GOMAXPROCS:     runtime.GOMAXPROCS(0),
		PID:            os.Getpid(),
		Uptime:         humanDuration(time.Since(r.startedAt)),
		AppEnvironment: facades.Config().GetString("app.env"),
		AppDebug:       facades.Config().GetBool("app.debug"),
		Database:       facades.Config().GetString("database.default"),
		SessionDriver:  facades.Config().GetString("session.default"),
		ModuleVersion:  moduleVersion,
	})
}

func (r *ServerController) ListTasks(ctx http.Context) http.Response {
	if r.tasks == nil {
		return serverTaskError(ctx, stdhttp.StatusServiceUnavailable, "Administrador de tareas no disponible.")
	}
	tasks, err := r.tasks.List()
	if err != nil {
		facades.Log().Errorf("no se pudieron listar plantillas systemd: %v", err)
		return serverTaskError(ctx, stdhttp.StatusInternalServerError, "No fue posible listar las plantillas.")
	}
	return ctx.Response().Json(stdhttp.StatusOK, map[string]any{"tasks": tasks})
}

func (r *ServerController) ReadTask(ctx http.Context) http.Response {
	if r.tasks == nil {
		return serverTaskError(ctx, stdhttp.StatusServiceUnavailable, "Administrador de tareas no disponible.")
	}
	task, err := r.tasks.Read(ctx.Request().Query("filename"))
	if err != nil {
		return serverTaskResponseFromError(ctx, err)
	}
	return ctx.Response().Json(stdhttp.StatusOK, task)
}

func (r *ServerController) CreateTask(ctx http.Context) http.Response {
	if r.tasks == nil {
		return serverTaskError(ctx, stdhttp.StatusServiceUnavailable, "Administrador de tareas no disponible.")
	}

	task, err := r.tasks.Create(appserver.CreateTaskInput{
		Name:             ctx.Request().Input("name"),
		Kind:             ctx.Request().Input("kind"),
		Description:      ctx.Request().Input("description"),
		ExecStart:        ctx.Request().Input("exec_start"),
		WorkingDirectory: ctx.Request().Input("working_directory"),
		User:             ctx.Request().Input("user"),
		Restart:          ctx.Request().Input("restart"),
		After:            ctx.Request().Input("after"),
		OnCalendar:       ctx.Request().Input("on_calendar"),
		TargetUnit:       ctx.Request().Input("target_unit"),
		Persistent:       strings.EqualFold(ctx.Request().Input("persistent"), "true") || ctx.Request().Input("persistent") == "1",
	})
	if err != nil {
		return serverTaskResponseFromError(ctx, err)
	}

	return ctx.Response().Json(stdhttp.StatusCreated, task)
}

func (r *ServerController) DeleteTask(ctx http.Context) http.Response {
	if r.tasks == nil {
		return serverTaskError(ctx, stdhttp.StatusServiceUnavailable, "Administrador de tareas no disponible.")
	}
	if err := r.tasks.Delete(ctx.Request().Input("filename")); err != nil {
		return serverTaskResponseFromError(ctx, err)
	}
	return ctx.Response().Json(stdhttp.StatusOK, map[string]any{"deleted": true})
}

func serverPageData(meta PageMeta) ServerPageData {
	return ServerPageData{
		Title:        meta.Title,
		PageTitle:    meta.PageTitle,
		PageSubtitle: meta.PageSubtitle,
		ActiveRoute:  meta.ActiveRoute,
		CSRFToken:    meta.CSRFToken,
		CurrentUser:  meta.CurrentUser,
	}
}

func serverTaskResponseFromError(ctx http.Context, err error) http.Response {
	switch {
	case errors.Is(err, fs.ErrExist):
		return serverTaskError(ctx, stdhttp.StatusConflict, "Ya existe una plantilla con ese nombre.")
	case errors.Is(err, fs.ErrNotExist):
		return serverTaskError(ctx, stdhttp.StatusNotFound, "La plantilla no existe.")
	case errors.Is(err, appserver.ErrInvalidTask):
		return serverTaskError(ctx, stdhttp.StatusBadRequest, "Los datos de la plantilla no son válidos.")
	default:
		facades.Log().Errorf("error administrando plantilla systemd: %v", err)
		return serverTaskError(ctx, stdhttp.StatusInternalServerError, "No fue posible completar la operación.")
	}
}

func serverTaskError(ctx http.Context, status int, message string) http.Response {
	return ctx.Response().Json(status, map[string]any{"error": message})
}

func humanDuration(value time.Duration) string {
	if value < time.Minute {
		return "< 1 min"
	}
	hours := int(value.Hours())
	minutes := int(value.Minutes()) % 60
	if hours > 0 {
		return fmt.Sprintf("%dh %02dm", hours, minutes)
	}
	return fmt.Sprintf("%dm", minutes)
}
