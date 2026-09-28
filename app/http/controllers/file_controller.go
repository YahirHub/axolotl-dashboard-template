package controllers

import (
	"errors"
	"io/fs"
	stdhttp "net/http"
	"os"

	"github.com/goravel/framework/contracts/http"

	"github.com/YahirHub/axolotl-dashboard-template/app/facades"
	appfiles "github.com/YahirHub/axolotl-dashboard-template/app/files"
)

type FilePageData struct {
	Title        string
	PageTitle    string
	PageSubtitle string
	ActiveRoute  string
	CSRFToken    string
	CurrentUser  CurrentUserView
	RootLabel    string
	MaxEditBytes int
}

type FileController struct {
	workspace *appfiles.Workspace
}

func NewFileController() *FileController {
	workspace, err := appfiles.NewConfiguredWorkspace()
	if err != nil {
		facades.Log().Errorf("no se pudo configurar explorador de archivos: %v", err)
		return &FileController{}
	}
	if err := workspace.EnsureDemo(); err != nil {
		facades.Log().Errorf("no se pudo preparar el workspace demo: %v", err)
		return &FileController{}
	}
	return &FileController{workspace: workspace}
}

func (r *FileController) Index(ctx http.Context) http.Response {
	meta := pageMeta(
		ctx,
		"Archivos — Axolotl",
		"Explorador de archivos",
		"Workspace demo controlado por el backend Go.",
		"files",
	)

	return ctx.Response().View().Make("files.tmpl", FilePageData{
		Title:        meta.Title,
		PageTitle:    meta.PageTitle,
		PageSubtitle: meta.PageSubtitle,
		ActiveRoute:  meta.ActiveRoute,
		CSRFToken:    meta.CSRFToken,
		CurrentUser:  meta.CurrentUser,
		RootLabel:    "Demo",
		MaxEditBytes: facades.Config().GetInt("files.max_edit_bytes", 2*1024*1024),
	})
}

func (r *FileController) List(ctx http.Context) http.Response {
	if r.workspace == nil {
		return fileJSONError(ctx, stdhttp.StatusServiceUnavailable, "Explorador no disponible.")
	}
	path := ctx.Request().Query("path")
	entries, err := r.workspace.List(path)
	if err != nil {
		return fileJSONFromError(ctx, err)
	}

	return ctx.Response().Json(stdhttp.StatusOK, map[string]any{
		"path":    path,
		"entries": entries,
	})
}

func (r *FileController) Read(ctx http.Context) http.Response {
	if r.workspace == nil {
		return fileJSONError(ctx, stdhttp.StatusServiceUnavailable, "Explorador no disponible.")
	}
	content, err := r.workspace.Read(ctx.Request().Query("path"))
	if err != nil {
		return fileJSONFromError(ctx, err)
	}
	return ctx.Response().Json(stdhttp.StatusOK, content)
}

func (r *FileController) Create(ctx http.Context) http.Response {
	if r.workspace == nil {
		return fileJSONError(ctx, stdhttp.StatusServiceUnavailable, "Explorador no disponible.")
	}
	entry, err := r.workspace.Create(
		ctx.Request().Input("parent"),
		ctx.Request().Input("name"),
		ctx.Request().Input("type"),
	)
	if err != nil {
		return fileJSONFromError(ctx, err)
	}
	return ctx.Response().Json(stdhttp.StatusCreated, entry)
}

func (r *FileController) Save(ctx http.Context) http.Response {
	if r.workspace == nil {
		return fileJSONError(ctx, stdhttp.StatusServiceUnavailable, "Explorador no disponible.")
	}
	if err := r.workspace.Save(
		ctx.Request().Input("path"),
		ctx.Request().Input("content"),
	); err != nil {
		return fileJSONFromError(ctx, err)
	}
	return ctx.Response().Json(stdhttp.StatusOK, map[string]any{
		"saved": true,
		"path":  ctx.Request().Input("path"),
	})
}

func (r *FileController) Delete(ctx http.Context) http.Response {
	if r.workspace == nil {
		return fileJSONError(ctx, stdhttp.StatusServiceUnavailable, "Explorador no disponible.")
	}
	if err := r.workspace.Delete(ctx.Request().Input("path")); err != nil {
		return fileJSONFromError(ctx, err)
	}
	return ctx.Response().Json(stdhttp.StatusOK, map[string]any{
		"deleted": true,
	})
}

func fileJSONFromError(ctx http.Context, err error) http.Response {
	switch {
	case errors.Is(err, fs.ErrNotExist):
		return fileJSONError(ctx, stdhttp.StatusNotFound, "El archivo o carpeta no existe.")
	case errors.Is(err, fs.ErrExist):
		return fileJSONError(ctx, stdhttp.StatusConflict, "Ya existe un elemento con ese nombre.")
	case errors.Is(err, appfiles.ErrInvalidPath):
		return fileJSONError(ctx, stdhttp.StatusBadRequest, "Ruta o nombre no permitido.")
	case errors.Is(err, appfiles.ErrTooLarge):
		return fileJSONError(ctx, stdhttp.StatusRequestEntityTooLarge, "El archivo supera el tamaño máximo editable.")
	case errors.Is(err, appfiles.ErrBinaryFile):
		return fileJSONError(ctx, stdhttp.StatusUnsupportedMediaType, "Este archivo no es texto UTF-8 editable.")
	case errors.Is(err, os.ErrPermission):
		return fileJSONError(ctx, stdhttp.StatusForbidden, "No hay permisos para completar la operación.")
	default:
		facades.Log().Errorf("error del explorador de archivos: %v", err)
		return fileJSONError(ctx, stdhttp.StatusInternalServerError, "No fue posible completar la operación.")
	}
}

func fileJSONError(ctx http.Context, status int, message string) http.Response {
	return ctx.Response().Json(status, map[string]any{
		"error": message,
	})
}
