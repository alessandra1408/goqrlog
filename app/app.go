package app

import (
	"github.com/alessandra1408/goqrlog/app/qrcode"
	"github.com/alessandra1408/goqrlog/internal/config"
)

type App struct {
	QRCode qrcode.App
	Sheets qrcode.App
}

type Options struct {
	Cfg config.Config
}

func New(opts Options) *App {
	return &App{
		QRCode: qrcode.NewApp(),
		Sheets: qrcode.NewApp(),
	}
}
