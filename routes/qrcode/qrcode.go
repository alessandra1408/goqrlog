package qrcode

import (
	"net/http"

	"github.com/alessandra1408/goqrlog/app"
	"github.com/alessandra1408/goqrlog/app/qrcode"
	"github.com/alessandra1408/goqrlog/internal/config"
	defaultErrors "github.com/alessandra1408/goqrlog/internal/errors"
	"github.com/alessandra1408/goqrlog/internal/middleware"
	"github.com/alessandra1408/goqrlog/pkg/log"
	"github.com/alessandra1408/goqrlog/pkg/model"
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
	req := new(qrcode.Request)

	if err := c.Bind(req); err != nil {
		h.log.Warnf(defaultErrors.ErrBindingMessage, err)
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Message: defaultErrors.BindingError(echo.MIMEApplicationJSON, c.Request()).Error(),
		})
	}

	h.log.Infof("received QR code payload: %+v", req)

	return c.JSON(http.StatusOK, req)
}
