package controllers

import (
	stdhttp "net/http"
	"time"

	"github.com/goravel/framework/contracts/http"

	"github.com/YahirHub/axolotl-dashboard-template/app/facades"
)

type HealthController struct{}

func NewHealthController() *HealthController {
	return &HealthController{}
}

func (r *HealthController) Show(ctx http.Context) http.Response {
	var probe struct {
		OK int `gorm:"column:ok"`
	}

	if err := facades.Orm().Query().Raw("SELECT 1 AS ok").Scan(&probe); err != nil || probe.OK != 1 {
		return ctx.Response().Json(stdhttp.StatusServiceUnavailable, http.Json{
			"status":   "unavailable",
			"database": "error",
			"time":     time.Now().UTC().Format(time.RFC3339),
		})
	}

	return ctx.Response().Json(stdhttp.StatusOK, http.Json{
		"status":   "ok",
		"database": "ok",
		"time":     time.Now().UTC().Format(time.RFC3339),
	})
}
