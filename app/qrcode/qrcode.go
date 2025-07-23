package qrcode

import (
	"context"

	"github.com/alessandra1408/goqrlog/internal/config"
	"github.com/alessandra1408/goqrlog/pkg/db"
	"github.com/alessandra1408/goqrlog/pkg/log"
	"github.com/jackc/pgx/v5"
)

type App interface {
	QRCodeHandler(ctx context.Context, req *Request, log log.Log) (pgx.Rows, error)
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

func (a *app) QRCodeHandler(ctx context.Context, req *Request, log log.Log) (pgx.Rows, error) {
	log.Info("QRCodeHandler called")

	query := "SELECT * FROM estudantes WHERE matricula = $1;"
	rows, err := a.db.Pool.Query(ctx, query, req.MatriculaAluno)
	if err != nil {
		log.Errorf("Error fetching data for matricula %d: %v", req.MatriculaAluno, err)
		return nil, err
	}

	return rows, nil
}

//"relation \"domjaimedb.estudantes\" does not exist"
