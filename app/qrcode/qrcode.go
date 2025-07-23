package qrcode

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/alessandra1408/goqrlog/internal/config"
	"github.com/alessandra1408/goqrlog/pkg/db"
	"github.com/alessandra1408/goqrlog/pkg/log"
	"github.com/jackc/pgx/v5"
)

var lastRequests sync.Map

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
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	key := fmt.Sprintf("%d", req.MatriculaAluno)
	now := time.Now().Unix()

	if val, ok := lastRequests.Load(key); ok {
		if last, ok := val.(int64); ok && now-last < 3 {
			log.Warnf("Ignored duplicate request for matricula %d", req.MatriculaAluno)
			return nil, nil // ou retornar erro de Too Many Requests
		}
	}
	lastRequests.Store(key, now)

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
