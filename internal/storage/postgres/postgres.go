package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/north-fy/golang-restapi-todolist/internal/config"
)

const op = "storage/postgres"

type Storage struct {
	conn *pgx.Conn
}

func NewStorage(cfg config.StorageConfig) *Storage {
	dataSourceName := fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DBname)

	conn, err := pgx.Connect(context.TODO(), dataSourceName)
	if err != nil {
		panic(err)
	}

	if err = conn.Ping(context.TODO()); err != nil {
		panic(err)
	}

	return &Storage{
		conn: conn,
	}
}

func (s *Storage) Close(ctx context.Context) {
	_ = s.conn.Close(ctx)
}
