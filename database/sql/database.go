package sql

import (
	"database/sql"
	"log"
	"time"
)

// WithConnMaxLifeTime to call func SetConnMaxLifetime on db, default is 10 seconds
func WithConnMaxLifeTime(i time.Duration) func(*sql.DB) {
	return func(db *sql.DB) {
		db.SetConnMaxLifetime(i)
	}
}

// WithMaxIdleConns to call func SetMaxIdleConns on db, default is 50
func WithMaxIdleConns(i int) func(*sql.DB) {
	return func(db *sql.DB) {
		db.SetMaxIdleConns(i)
	}
}

// WithMaxOpenConns to call func SetMaxOpenConns on db, default is 50
func WithMaxOpenConns(i int) func(*sql.DB) {
	return func(db *sql.DB) {
		db.SetMaxOpenConns(i)
	}
}

func New(dsn string, configFuncs ...func(*sql.DB)) *sql.DB {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalln("error while trying to create connection pool", dsn, err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalln("error while trying to ping", dsn, err.Error())
	}

	db.SetConnMaxLifetime(10 * time.Second)
	db.SetMaxIdleConns(50)
	db.SetMaxOpenConns(50)

	for _, f := range configFuncs {
		f(db)
	}

	return db
}
