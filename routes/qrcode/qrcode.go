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
	"github.com/alessandra1408/goqrlog/pkg/util"
	"github.com/labstack/echo/v4"
)

const (
	trunkedTokenString  = "trunked-token"
	remoteAddressString = "remote-address"
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
	ctx := c.Request().Context()
	token := c.Request().Header.Get("Authorization")
	trunkedToken := util.GetMaskedToken(token)

	l := log.
		LogWithUserAgent(h.log, c.Request().UserAgent()).
		With(trunkedTokenString, trunkedToken).
		With(remoteAddressString, c.RealIP())

	req := new(qrcode.Request)
	if err := c.Bind(req); err != nil {
		h.log.Warnf(defaultErrors.ErrBindingMessage, err)
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Message: defaultErrors.BindingError(echo.MIMEApplicationJSON, c.Request()).Error(),
		})
	}

	res, err := h.apps.QRCode.QRCodeHandler(ctx, req, l)
	if err != nil {
		h.log.Errorf("Error processing QR code request: %v", err)
		return c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Message: "Internal server error",
		})
	}

	h.log.Infof("response from QR code handler: %+v", res)

	if res == nil {
		h.log.Warn("No data found for the provided QR code request")
		return c.JSON(http.StatusNotFound, model.ErrorResponse{
			Message: "No data found for the provided QR code request",
		})
	}

	h.log.Infof("received QR code payload: %+v", req)

	return c.JSON(http.StatusOK, res)
}
