package postgres

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

type DbStorage struct {
	db *pgxpool.Pool
}

// Конструктор БД
func New(connstr string) (*DbStorage, error) {
	db, err := pgxpool.Connect(context.Background(), connstr)
	if err != nil {
		return nil, err
	}
	st := DbStorage{
		db,
	}
	return &st, nil
}
