package workspace

import (
	"context"
	"database/sql"

	pgx "github.com/jackc/pgx/v5"

	"github.com/amartya2002/secretlane/internal/config"
)

// Repository encapsulates all DB operations for workspaces.
// It works with either sqlite (*sql.DB) or postgres (*pgx.Conn) based on config.DBDriver.
type Repository struct {
	sqlDB   *sql.DB
	pgxConn *pgx.Conn
}

func NewRepository(sqlDB *sql.DB, pgxConn *pgx.Conn) *Repository {
	return &Repository{sqlDB: sqlDB, pgxConn: pgxConn}
}

func NewDefaultRepository() *Repository {
	return &Repository{sqlDB: config.DB, pgxConn: config.PGXConn}
}

func (r *Repository) CountByNameForUser(name string, userID int) (int, error) {
	var count int

	if config.DBDriver == "postgres" {
		row := r.pgxConn.QueryRow(context.Background(), `
		SELECT COUNT(*) FROM workspaces WHERE name = $1 AND created_by = $2
		`, name, userID)
		if err := row.Scan(&count); err != nil {
			return 0, err
		}
		return count, nil
	}

	row := r.sqlDB.QueryRow(`
		SELECT COUNT(*) FROM workspaces WHERE name = ? AND created_by = ?
	`, name, userID)
	if err := row.Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

func (r *Repository) CreateWorkspace(name, description string, userID int) (int, error) {
	if config.DBDriver == "postgres" {
		row := r.pgxConn.QueryRow(context.Background(), `
			INSERT INTO workspaces (name, description, created_by)
			VALUES ($1, $2, $3)
			RETURNING id
		`, name, description, userID)
		var id int
		if err := row.Scan(&id); err != nil {
			return 0, err
		}
		return id, nil
	}

	res, err := r.sqlDB.Exec(`
		INSERT INTO workspaces (name, description, created_by)
		VALUES (?, ?, ?)
	`, name, description, userID)
	if err != nil {
		return 0, err
	}
	lastID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(lastID), nil
}

func (r *Repository) ListForUser(userID int) ([]Workspace, error) {
	if config.DBDriver == "postgres" {
		rows, err := r.pgxConn.Query(context.Background(), `
		SELECT id, name, description, created_by, created_at
		FROM workspaces WHERE created_by = $1
		ORDER BY created_at DESC
		`, userID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		var list []Workspace
		for rows.Next() {
			var ws Workspace
			if err := rows.Scan(&ws.ID, &ws.Name, &ws.Description, &ws.CreatedBy, &ws.CreatedAt); err != nil {
				return nil, err
			}
			list = append(list, ws)
		}
		return list, nil
	}

	rows, err := r.sqlDB.Query(`
		SELECT id, name, description, created_by, created_at
		FROM workspaces WHERE created_by = ?
		ORDER BY created_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []Workspace
	for rows.Next() {
		var ws Workspace
		if err := rows.Scan(&ws.ID, &ws.Name, &ws.Description, &ws.CreatedBy, &ws.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, ws)
	}
	return list, nil
}

func (r *Repository) Update(id int, name, description string, userID int) error {
	if config.DBDriver == "postgres" {
		_, err := r.pgxConn.Exec(context.Background(), `
		UPDATE workspaces
		SET name = $1, description = $2
		WHERE id = $3 AND created_by = $4
		`, name, description, id, userID)
		return err
	}

	_, err := r.sqlDB.Exec(`
		UPDATE workspaces
		SET name = ?, description = ?
		WHERE id = ? AND created_by = ?
	`, name, description, id, userID)
	return err
}

func (r *Repository) Delete(id int, userID int) error {
	if config.DBDriver == "postgres" {
		_, err := r.pgxConn.Exec(context.Background(), `
		DELETE FROM workspaces
		WHERE id = $1 AND created_by = $2
		`, id, userID)
		return err
	}

	_, err := r.sqlDB.Exec(`
		DELETE FROM workspaces
		WHERE id = ? AND created_by = ?
	`, id, userID)
	return err
}
