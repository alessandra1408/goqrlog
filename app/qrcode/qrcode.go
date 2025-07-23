package qrcode

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/alessandra1408/goqrlog/internal/config"
	"github.com/alessandra1408/goqrlog/pkg/db"
	"github.com/alessandra1408/goqrlog/pkg/log"
)

var lastRequests sync.Map

type App interface {
	QRCodeHandler(ctx context.Context, req *Request, log log.Log) (Response, error)
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

func (a *app) QRCodeHandler(ctx context.Context, req *Request, log log.Log) (Response, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	key := fmt.Sprintf("%d", req.MatriculaAluno)
	now := time.Now().Unix()

	if val, ok := lastRequests.Load(key); ok {
		if last, ok := val.(int64); ok && now-last < 3 {
			log.Warnf("Ignored duplicate request for matricula %d", req.MatriculaAluno)
			return Response{}, nil // ou retornar erro de Too Many Requests
		}
	}
	lastRequests.Store(key, now)

	log.Info("QRCodeHandler called")

	query := "SELECT * FROM estudantes WHERE matricula = $1;"
	rows, err := a.db.Pool.Query(ctx, query, req.MatriculaAluno)
	if err != nil {
		log.Errorf("Error fetching data for matricula %d: %v", req.MatriculaAluno, err)
		return Response{}, err
	}

	defer rows.Close()

	var estudante Response
	if rows.Next() {
		err := rows.Scan(&estudante.Turma,
			&estudante.Matricula,
			&estudante.Estudante)
		if err != nil {
			log.Errorf("Error scanning row: %v", err)
			return Response{}, err
		}
	} else {
		return Response{}, fmt.Errorf("no estudante found for matricula %d", req.MatriculaAluno)
	}

	// Você pode retornar o estudante como JSON na camada de rota
	// Exemplo: na rota, faça json.NewEncoder(w).Encode(estudante)

	// Se quiser retornar o estudante diretamente aqui:
	// return nil, estudante

	// Como a assinatura exige (pgx.Rows, error), apenas retorne rows e nil

	return estudante, nil
}

//"relation \"domjaimedb.estudantes\" does not exist"
