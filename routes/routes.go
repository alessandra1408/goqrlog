package routes

import (
	"github.com/alessandra1408/goqrlog/app"
	"github.com/alessandra1408/goqrlog/internal/config"
	"github.com/alessandra1408/goqrlog/pkg/log"
	"github.com/alessandra1408/goqrlog/routes/qrcode"
	"github.com/alessandra1408/goqrlog/routes/sheets"
	"github.com/labstack/echo/v4"
)

type Options struct {
	Group *echo.Group
	Apps  *app.App
	Cfg   config.Config
	Log   log.Log
}

func RegisterRoutes(opts Options) {
	qrcode.Routes(opts.Group, opts.Apps, opts.Cfg, opts.Log)
	sheets.Routes(opts.Group, opts.Apps, opts.Log)

	opts.Log.Info("Registered API routes successfully")
}
