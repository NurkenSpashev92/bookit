package configs

import (
	"fmt"
)

type DBConfig struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	AppPort    string
}

func (c *DBConfig) DatabaseURL() string {
	return fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		c.DBUser,
		c.DBPassword,
		c.DBHost,
		c.DBPort,
		c.DBName,
	)
}

func NewDBConfig() *DBConfig {
	return &DBConfig{
		DBHost:     getEnv("POSTGRES_HOST", "localhost"),
		DBPort:     getEnv("POSTGRES_PORT", "5432"),
		DBUser:     getEnv("POSTGRES_USER", "bookit"),
		DBPassword: getEnv("POSTGRES_PASSWORD", "bookit"),
		DBName:     getEnv("POSTGRES_DB", "bookit"),
		AppPort:    getEnv("APP_PORT", "8080"),
	}
}
