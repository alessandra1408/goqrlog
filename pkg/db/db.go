package db

import (
	"context"
	"fmt"
	"time"

	"github.com/alessandra1408/goqrlog/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Database struct {
	Pool *pgxpool.Pool
}

func NewDatabase(cfg *config.Database) (*Database, error) {
	connStr := fmt.Sprintf("%s://%s:%s@%s:%d/%s?sslmode=%s&channel_binding=%s",
		cfg.Scheme, cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Name,
		cfg.SSLMode, cfg.ChannelBinding,
	)

	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, fmt.Errorf("não foi possível parsear configuração do pool: %w", err)
	}

	// 🔧 Ajuste de timeouts e healthcheck
	config.MaxConnIdleTime = 30 * time.Second   // conexões inativas serão fechadas após 30s
	config.MaxConnLifetime = 30 * time.Minute   // vida máxima da conexão
	config.HealthCheckPeriod = 30 * time.Second // verificação de saúde do pool

	// Opcional: limitar o número de conexões simultâneas, dependendo do Neon
	// config.MaxConns = 10

	// Use um context com timeout curto só para criar o pool
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("não foi possível conectar: %w", err)
	}

	return &Database{Pool: pool}, nil
}

func (db *Database) Get(ctx context.Context, matricula int, dest ...any) error {
	query := fmt.Sprintf("SELECT * FROM domjaimedb WHERE matricula = %d", matricula)

	return db.Pool.QueryRow(ctx, query).Scan(dest...)
}

func (db *Database) Exec(ctx context.Context, query string, args ...any) (int64, error) {
	ct, err := db.Pool.Exec(ctx, query, args...)
	if err != nil {
		return 0, err
	}
	return ct.RowsAffected(), nil
}

func (db *Database) Query(ctx context.Context, query string, args ...any) ([][]any, error) {
	rows, err := db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result [][]any

	for rows.Next() {
		values, err := rows.Values()
		if err != nil {
			return nil, err
		}
		result = append(result, values)
	}

	return result, rows.Err()
}
