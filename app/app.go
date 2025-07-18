package app

import (
	"net/http"

	"github.com/alessandra1408/goqrlog/app/qrcode"
	"github.com/alessandra1408/goqrlog/internal/config"
	"github.com/alessandra1408/goqrlog/pkg/db"
)

type App struct {
	QRCode qrcode.App
	Sheets qrcode.App
}

type Options struct {
	Cfg        config.Config
	HttpClient *http.Client
	DB         *db.Database
}

func New(opts Options) *App {
	return &App{
		QRCode: qrcode.NewApp(opts.Cfg, opts.DB),
	}
}
