package auth

import (
	"context"
	"database/sql"

	pgx "github.com/jackc/pgx/v5"

	"github.com/amartya2002/secretlane/internal/config"
)

// Repository encapsulates all DB operations for auth.
// It works with either sqlite (*sql.DB) or postgres (*pgx.Conn) based on config.DBDriver.
type Repository struct {
	sqlDB  *sql.DB
	pgxConn *pgx.Conn
}

func NewRepository(sqlDB *sql.DB, pgxConn *pgx.Conn) *Repository {
	return &Repository{sqlDB: sqlDB, pgxConn: pgxConn}
}

func NewDefaultRepository() *Repository {
	return &Repository{sqlDB: config.DB, pgxConn: config.PGXConn}
}

func (r *Repository) FindByUsername(username string) (*User, error) {
	u := &User{}

	if config.DBDriver == "postgres" {
		row := r.pgxConn.QueryRow(context.Background(),
			`SELECT id, username, password FROM users WHERE username = $1`, username)
		if err := row.Scan(&u.ID, &u.Username, &u.Password); err != nil {
			return nil, err
		}
		return u, nil
	}

	row := r.sqlDB.QueryRow(`SELECT id, username, password FROM users WHERE username = ?`, username)
	if err := row.Scan(&u.ID, &u.Username, &u.Password); err != nil {
		return nil, err
	}
	return u, nil
}

func (r *Repository) UserExists(username string) (bool, error) {
	if config.DBDriver == "postgres" {
		var id int
		row := r.pgxConn.QueryRow(context.Background(),
			`SELECT id FROM users WHERE username = $1`, username)
		err := row.Scan(&id)
		if err == pgx.ErrNoRows {
			return false, nil
		}
		if err != nil {
			return false, err
		}
		return true, nil
	}

	var id int
	err := r.sqlDB.QueryRow(`SELECT id FROM users WHERE username = ?`, username).Scan(&id)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *Repository) CreateUser(username, password string) (*User, error) {
	u := &User{
		Username: username,
		Password: password,
	}

	if config.DBDriver == "postgres" {
		row := r.pgxConn.QueryRow(context.Background(), `
			INSERT INTO users (username, password)
			VALUES ($1, $2)
			RETURNING id
		`, username, password)
		if err := row.Scan(&u.ID); err != nil {
			return nil, err
		}
		return u, nil
	}

	res, err := r.sqlDB.Exec(`
		INSERT INTO users (username, password)
		VALUES (?, ?)
	`, username, password)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	u.ID = int(id)
	return u, nil
}
