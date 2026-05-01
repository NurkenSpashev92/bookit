package store

import (
	"context"
	"fmt"
	"runtime"
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

	cfg, err := pgxpool.ParseConfig(conf.DatabaseURL())
	if err != nil {
		return nil, fmt.Errorf("failed to parse db config: %w", err)
	}

	// Pool sizing: 6× CPU is a sane starting point for I/O-bound workload.
	// Capped low-end at 25 so dev with NumCPU<5 still has enough headroom.
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
