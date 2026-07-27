package store

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/nurkenspashev92/bookit/configs"
)

type Database struct {
	Conn *pgxpool.Pool
}

var (
	dbOnce     sync.Once
	dbInstance *Database
	dbErr      error
)

func NewPostgresDb(conf *configs.DBConfig) (*Database, error) {
	dbOnce.Do(func() {
		dbInstance, dbErr = newPostgresDb(conf)
	})

	return dbInstance, dbErr
}

func newPostgresDb(conf *configs.DBConfig) (*Database, error) {
	cfg, err := pgxpool.ParseConfig(conf.DatabaseURL())
	if err != nil {
		return nil, fmt.Errorf("failed to parse db config: %w", err)
	}

	maxConns := int32(runtime.NumCPU() * 6)
	if maxConns < 25 {
		maxConns = 25
	}
	cfg.MaxConns = maxConns
	cfg.MinConns = 5
	cfg.MaxConnLifetime = time.Hour
	cfg.MaxConnIdleTime = 30 * time.Minute
	cfg.HealthCheckPeriod = time.Minute

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to connect db: %w", err)
	}

	if err := conn.Ping(ctx); err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to ping db: %w", err)
	}

	dbInstance = &Database{Conn: conn}
	return dbInstance, nil
}

func (db *Database) Close() {
	if db != nil && db.Conn != nil {
		db.Conn.Close()
	}
}
