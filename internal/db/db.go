package database

import (
	"fmt"
	"net/url"

	"github.com/AZhur771/wg-grpc-api/internal/app"
	"github.com/jmoiron/sqlx"
)

func GetConnectionString(config app.Config) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable&connect_timeout=%d",
		url.QueryEscape(config.DBUsername),
		url.QueryEscape(config.DBPassword),
		config.DBHost,
		config.DBPort,
		config.DBName,
		config.DBTimeout,
	)
}

func New(config app.Config) (*sqlx.DB, error) {
	connectionString := GetConnectionString(config)
	db, err := sqlx.Open("pgx", connectionString)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(config.DBMaxOpenConnections)
	db.SetMaxIdleConns(config.DBMaxIdleConnections)

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("error %s database connection %w", connectionString, err)
	}

	return db, nil
}
