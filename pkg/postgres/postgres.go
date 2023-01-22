package postgres

import (
	"fmt"
	"github.com/jackc/pgx"
)

func NewConnect() (*pgx.ConnPool, error) {
	conn, err := pgx.ParseConnectionString(fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", "forum", "forum", "localhost", "5432", "forum"))

	if err != nil {
		return nil, err
	}

	db, err := pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig:     conn,
		MaxConnections: 100,
		AfterConnect:   nil,
		AcquireTimeout: 0,
	})

	if err != nil {
		return nil, err
	}

	return db, nil
}
