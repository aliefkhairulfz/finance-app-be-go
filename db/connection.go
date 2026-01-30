package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

type Connection struct{}

func (c *Connection) Connect(ctx context.Context) (*pgxpool.Pool, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	dsn := os.Getenv("DB_CONNECTION_STRING")
	if dsn == "" {
		return nil, fmt.Errorf("DB_CONNECTION_STRING is not set")
	}

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}

	fmt.Println("Connected to PostgreSQL (pgx)")
	return pool, nil
}
