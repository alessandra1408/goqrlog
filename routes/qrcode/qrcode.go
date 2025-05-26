package qrcode

import (
	"net/http"

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
	body, err := c.Request().Body, error(nil)
	defer body.Close()
	data := make([]byte, c.Request().ContentLength)
	_, err = body.Read(data)
	if err != nil {
		h.log.Error("failed to read body: %v", err)
		return c.String(http.StatusBadRequest, "failed to read body")
	}

	h.log.Info("QRCodeHandler called")
	return c.Blob(http.StatusOK, "application/octet-stream", data)
}
