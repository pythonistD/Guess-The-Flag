package db

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pythonistD/Guess-The-Flag/internal/config"
)

func NewPostgres(cfg config.DBConfig) (*sqlx.DB, error) {
	connectStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DbName)
	db, err := sqlx.Connect("postgres", connectStr)
	if err != nil {
		return nil, fmt.Errorf("error connecting to postgres: %w", err)
	}
	db.SetMaxIdleConns(cfg.MaxIdleConnections)
	db.SetMaxOpenConns(cfg.MaxOpenConnections)
	db.SetConnMaxLifetime(cfg.MaxConnectionLifeTime)
	return db, nil
}
