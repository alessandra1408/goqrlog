package db

import (
	"context"
	"fmt"

	"github.com/alessandra1408/goqrlog/internal/config"
	"github.com/jackc/pgx/v5"
)

type Database struct {
	Conn *pgx.Conn
}

func NewDatabase(cfg *config.Database) (*Database, error) {
	connStr := fmt.Sprintf("%s://%s:%s@%s:%d/%s?sslmode=%s&channel_binding=%s",
		cfg.Scheme, cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Name,
		cfg.SSLMode, cfg.ChannelBinding,
	)

	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		return nil, fmt.Errorf("não foi possível conectar: %w", err)
	}

	return &Database{Conn: conn}, nil
}

func (db *Database) Get(ctx context.Context, matricula int, dest ...any) error {
	query := fmt.Sprintf("SELECT * FROM domjaimedb WHERE matricula = %d", matricula)

	return db.Conn.QueryRow(ctx, query).Scan(dest...)
}

func (db *Database) Exec(ctx context.Context, query string, args ...any) (int64, error) {
	ct, err := db.Conn.Exec(ctx, query, args...)
	if err != nil {
		return 0, err
	}
	return ct.RowsAffected(), nil
}

func (db *Database) Query(ctx context.Context, query string, args ...any) ([][]any, error) {
	rows, err := db.Conn.Query(ctx, query, args...)
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
