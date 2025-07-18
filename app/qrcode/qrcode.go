package qrcode

import (
	"context"
	"fmt"

	"github.com/alessandra1408/goqrlog/internal/config"
	"github.com/alessandra1408/goqrlog/pkg/db"
	"github.com/alessandra1408/goqrlog/pkg/log"
)

type App interface {
	QRCodeHandler(ctx context.Context, req *Request, log log.Log) ([][]any, error)
}

type app struct {
	cfg config.Config
	db  *db.Database
}

func NewApp(cfg config.Config, db *db.Database) App {
	return &app{
		cfg: cfg,
		db:  db,
	}
}

func (a *app) QRCodeHandler(ctx context.Context, req *Request, log log.Log) ([][]any, error) {
	log.Info("QRCodeHandler called")

	query := fmt.Sprintf("SELECT * FROM estudantes e WHERE e.matricula = '%d';", req.MatriculaAluno)

	result, err := a.db.Query(ctx, query)
	if err != nil {
		log.Errorf("Error fetching data for matricula %d: %v", req.MatriculaAluno, err)
		return nil, err
	}

	return result, nil
}

//"relation \"domjaimedb.estudantes\" does not exist"
