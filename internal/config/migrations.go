package config

import "log"

// RunMigrations creates all required tables.
// Call this AFTER InitDatabase().
func RunMigrations() {
	if DBDriver == "postgres" {
		runPostgresMigrations()
	} else {
		runSQLiteMigrations()
	}
}

func runSQLiteMigrations() {
	// DEV ONLY: Drop everything before recreating.
	_, err := DB.Exec(`
        DROP TABLE IF EXISTS workspaces;
        DROP TABLE IF EXISTS users;
    `)
	if err != nil {
		log.Fatalf("[MIGRATION] failed dropping tables (sqlite): %v", err)
	}

	// USERS
	_, err = DB.Exec(`
        CREATE TABLE IF NOT EXISTS users (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            username TEXT UNIQUE NOT NULL,
            password TEXT NOT NULL
        );
    `)
	if err != nil {
		log.Fatalf("[MIGRATION] failed creating users table (sqlite): %v", err)
	}

	if App.SeedDefaultUser {
		_, err = DB.Exec(`
    INSERT INTO users (username, password)
        VALUES ('admin@local', 'ChangeMe123!')
    ON CONFLICT(username) DO NOTHING;
`)
		if err != nil {
			log.Fatalf("[MIGRATION] failed inserting default user (sqlite): %v", err)
		}
	}

	// WORKSPACES
	_, err = DB.Exec(`
        CREATE TABLE IF NOT EXISTS workspaces (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT NOT NULL,
            description TEXT,
            created_by INTEGER NOT NULL,
            created_at TEXT DEFAULT (datetime('now')),
            FOREIGN KEY (created_by) REFERENCES users(id)
        );
    `)
	if err != nil {
		log.Fatalf("[MIGRATION] failed creating workspaces table (sqlite): %v", err)
	}

	log.Println("[MIGRATION] Users and workspaces tables created successfully (sqlite)")
}

func runPostgresMigrations() {
	// DEV ONLY: drop and recreate just what we need (users + workspaces).
	_, err := DB.Exec(`
        DROP TABLE IF EXISTS workspaces;
        DROP TABLE IF EXISTS users;
    `)
	if err != nil {
		log.Fatalf("[MIGRATION] failed dropping tables (postgres): %v", err)
	}

	_, err = DB.Exec(`
        CREATE TABLE IF NOT EXISTS users (
            id SERIAL PRIMARY KEY,
            username TEXT UNIQUE NOT NULL,
            password TEXT NOT NULL
        );
    `)
	if err != nil {
		log.Fatalf("[MIGRATION] failed creating users table (postgres): %v", err)
	}

	if App.SeedDefaultUser {
		_, err = DB.Exec(`
        INSERT INTO users (username, password)
        VALUES ('admin@local', 'ChangeMe123!')
        ON CONFLICT (username) DO NOTHING;
    `)
		if err != nil {
			log.Fatalf("[MIGRATION] failed inserting default user (postgres): %v", err)
		}
	}

	_, err = DB.Exec(`
        CREATE TABLE IF NOT EXISTS workspaces (
            id SERIAL PRIMARY KEY,
            name TEXT NOT NULL,
            description TEXT,
            created_by INTEGER NOT NULL REFERENCES users(id),
            created_at TIMESTAMPTZ DEFAULT now()
        );
    `)
	if err != nil {
		log.Fatalf("[MIGRATION] failed creating workspaces table (postgres): %v", err)
	}

	log.Println("[MIGRATION] Users and workspaces tables created successfully (postgres)")
}

