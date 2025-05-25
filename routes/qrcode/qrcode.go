package qrcode

import (
	"github.com/alessandra1408/goqrlog/app"
	"github.com/alessandra1408/goqrlog/internal/config"
	"github.com/alessandra1408/goqrlog/middleware"
	"github.com/alessandra1408/goqrlog/pkg/log"
	"github.com/labstack/echo/v4"
)

type handler struct {
	apps *app.App
	log  log.Log
}

func Routes(g *echo.Group, apps *app.App, cfg config.Config, log log.Log) {
	handler := &handler{
		apps: apps,
		log:  log,
	}

	mgtmGroup := g.Group("/mgtm")

	mgtmGroup.POST("/qrcode/ingest", handler.QRCodeHandler, middleware.AuthMiddleware(cfg.Auth.Key))

}

func (h *handler) QRCodeHandler(c echo.Context) error {
	// Implement the QRCodeHandler logic here
	// For example, you might want to generate a QR code and return it
	h.log.Info("QRCodeHandler called")
	return nil
}
