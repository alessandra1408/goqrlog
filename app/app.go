package app

import "github.com/alessandra1408/goqrlog/app/qrcode"

type App struct {
	QRCode qrcode.App
	Sheets qrcode.App
}

func New() *App {
	return &App{
		QRCode: qrcode.NewApp(),
		Sheets: qrcode.NewApp(),
	}
}
