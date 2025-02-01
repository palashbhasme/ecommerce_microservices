package common

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresConfig struct {
	Host     string
	Port     int
	Password string
	User     string
	DbName   string
	SSLMode  string
}

func ensureDBExists(config *PostgresConfig) error {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s", config.Host, config.User, config.Password, config.DbName, config.Port, config.SSLMode)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to open connection to database: %w", err)
	}

	defer db.Close()
	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", config.DbName))
	if err != nil && err.Error() != fmt.Sprintf(`pq: database "%s" already exists`, config.DbName) {
		return fmt.Errorf("failed to create database: %w", err)
	}
	return nil

}
func ConnectToDb(config *PostgresConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s", config.Host, config.User, config.Password, config.DbName, config.Port, config.SSLMode)

	err := ensureDBExists(config)

	if err != nil {
		return nil, err
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {

		return nil, err
	}
	return db, nil

}
