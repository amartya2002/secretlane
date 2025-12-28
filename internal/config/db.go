package config

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	pgx "github.com/jackc/pgx/v5"
	_ "github.com/mattn/go-sqlite3"
)

var (
	// DB is used in sqlite mode.
	DB *sql.DB
	// PGXConn is used in postgres mode.
	PGXConn *pgx.Conn
)

// InitDatabase connects to SQLite (db-less mode) or Postgres depending on config.
// Call this once in main() before starting the server.
func InitDatabase() {
	driver := DBDriver
	if driver == "" {
		driver = "sqlite"
		DBDriver = driver
	}

	switch driver {
	case "postgres":
		conn, err := initPostgres()
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to init postgres: %v\n", err)
			os.Exit(1)
		}
		if err := conn.Ping(context.Background()); err != nil {
			fmt.Fprintf(os.Stderr, "failed to ping postgres: %v\n", err)
			os.Exit(1)
		}
		PGXConn = conn
	default:
		db, err := initSQLite()
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to init sqlite: %v\n", err)
			os.Exit(1)
		}
		if err := db.Ping(); err != nil {
			fmt.Fprintf(os.Stderr, "failed to ping sqlite: %v\n", err)
			os.Exit(1)
		}
		DB = db
	}

	log.Printf("[DB] Connected using driver=%s\n", DBDriver)
}

func initSQLite() (*sql.DB, error) {
	// Local file-based SQLite: DB-less mode.
	dsn := "./sqlite-secretlane.db"
	return sql.Open("sqlite3", dsn)
}

func initPostgres() (*pgx.Conn, error) {
	// Build DSN from config.DBConfig (set by LoadAppConfig).
	host := DBConfig.Host
	port := DBConfig.Port
	user := DBConfig.User
	password := DBConfig.Password
	dbname := DBConfig.DBName
	sslmode := DBConfig.SSLMode

	if host == "" {
		host = "localhost"
	}
	if port == 0 {
		port = 5432
	}
	if user == "" {
		user = "postgres"
	}
	if dbname == "" {
		dbname = "secretlane"
	}
	if sslmode == "" {
		sslmode = "disable"
	}

	connString := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode,
	)

	// Use pgx directly as in the official docs example.
	return pgx.Connect(context.Background(), connString)
}
