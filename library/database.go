package library

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

type Database struct {
	Name         string
	Host         string
	Port         string
	Username     string
	Password     string
	Params       string
	MaxIdleTime  time.Duration
	MaxLifeTime  time.Duration
	MaxIdleConns int
	MaxOpenConns int
}

func NewDatabaseConnection(dbCfg Database) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn(dbCfg))
	if err != nil {
		return db, fmt.Errorf("open con : %w", err)
	}

	err = db.Ping()
	if err != nil {
		return db, fmt.Errorf("db ping : %w", err)
	}

	db.SetConnMaxIdleTime(time.Minute * dbCfg.MaxIdleTime)
	db.SetConnMaxLifetime(time.Minute * dbCfg.MaxLifeTime)
	db.SetMaxIdleConns(dbCfg.MaxIdleConns)
	db.SetMaxOpenConns(dbCfg.MaxOpenConns)

	return db, nil
}

/**
func dsn(dbCfg Database) string {
	if dbCfg.Env == "production" {
		return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=verify-full&sslrootcert=ap-southeast-1-bundle.pem", dbCfg.Username, dbCfg.Password, dbCfg.Host, dbCfg.Port, dbCfg.Name)
	}
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbCfg.Username, dbCfg.Password, dbCfg.Host, dbCfg.Port, dbCfg.Name)
}
**/

func dsn(dbCfg Database) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?%s", dbCfg.Username, dbCfg.Password, dbCfg.Host, dbCfg.Port, dbCfg.Name, dbCfg.Params)
}
