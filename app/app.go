package app

import (
	"net/http"

	"github.com/alessandra1408/goqrlog/app/qrcode"
	"github.com/alessandra1408/goqrlog/internal/config"
)

type App struct {
	QRCode qrcode.App
	Sheets qrcode.App
}

type Options struct {
	Cfg        config.Config
	HttpClient *http.Client
}

func New(opts Options) *App {
	return &App{
		QRCode: qrcode.NewApp(opts.Cfg),
	}
}
