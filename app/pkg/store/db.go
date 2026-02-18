package store

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/nurkenspashev92/bookit/configs"
)

type Database struct {
	Conn *pgxpool.Pool
}

var dbInstance *Database

func NewPostgresDb(conf *configs.DBConfig) (*Database, error) {
	if dbInstance != nil {
		return dbInstance, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := pgxpool.New(ctx, conf.DatabaseURL())
	if err != nil {
		return nil, fmt.Errorf("failed to connect db: %w", err)
	}

	if err := conn.Ping(ctx); err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to ping db: %w", err)
	}

	dbInstance = &Database{
		Conn: conn,
	}

	return dbInstance, nil
}

func (db *Database) Close() {
	if db != nil && db.Conn != nil {
		db.Conn.Close()
	}
}
