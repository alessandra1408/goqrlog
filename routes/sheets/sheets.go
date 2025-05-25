package sheets

import (
	"github.com/alessandra1408/goqrlog/app"
	"github.com/alessandra1408/goqrlog/pkg/log"
	"github.com/labstack/echo/v4"
)

type handler struct {
	apps *app.App
	log  log.Log
}

func Routes(g *echo.Group, apps *app.App, log log.Log) {
	handler := &handler{
		apps: apps,
		log:  log,
	}

	mgtmGroup := g.Group("/mgtm")

	mgtmGroup.POST("/sheets/send", handler.SheetsHandler)

}

func (h *handler) SheetsHandler(c echo.Context) error {
	// Implement the QRCodeHandler logic here
	// For example, you might want to generate a QR code and return it
	h.log.Info("QRCodeHandler called")
	return nil
}
