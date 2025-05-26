package qrcode

import (
	"context"

	"github.com/alessandra1408/goqrlog/internal/config"
	"github.com/alessandra1408/goqrlog/pkg/log"
)

type App interface {
	QRCodeHandler(ctx context.Context, req *Request, log log.Log) error
}

type app struct {
	cfg config.Config
}

func NewApp(cfg config.Config) App {
	return &app{
		cfg: cfg,
	}
}

func (a *app) QRCodeHandler(ctx context.Context, req *Request, log log.Log) error {
	// Implement the QRCodeHandler logic here
	// For example, you might want to generate a QR code and return it
	log.Info("QRCodeHandler called")
	return nil
}
