package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"gopkg.in/yaml.v3"
)

// Config holds application configuration loaded from YAML/env.
type Config struct {
	App      AppConfig      `yaml:"app"`
	Database DatabaseConfig `yaml:"database"`
	Postgres PostgresConfig `yaml:"postgres"`
}

type AppConfig struct {
	// Port the HTTP server listens on.
	Port int `yaml:"port"`
	// EnableFrontend toggles serving the frontend (if any).
	EnableFrontend bool `yaml:"enable_frontend"`
	// SeedDefaultUser controls whether migrations create a default admin user.
	SeedDefaultUser bool `yaml:"seed_default_user"`
}

// DatabaseConfig controls which DB driver is used.
type DatabaseConfig struct {
	Driver string `yaml:"driver"`
}

// PostgresConfig holds basic Postgres connection details.
type PostgresConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
	SSLMode  string `yaml:"sslmode"`
}

// App is the runtime application configuration used by the rest of the code.
// Port is stringified here for easy use in http.ListenAndServe.
type AppRuntimeConfig struct {
	Port           string
	EnableFrontend bool
	SeedDefaultUser bool
}

var (
	// App holds app-level runtime configuration.
	App AppRuntimeConfig

	// DBConfig holds loaded Postgres configuration; currently unused but kept
	// for upcoming Postgres integration.
	DBConfig PostgresConfig

	// DBDriver is the selected database driver ("sqlite" or "postgres").
	DBDriver string
)

// LoadAppConfig initialises application configuration from config.yaml and env.
// Precedence:
//   1. Defaults
//   2. config.yaml (if present)
//   3. Environment variables
func LoadAppConfig() error {
	// Defaults
	cfg := Config{
		App: AppConfig{
			Port:            8080,
			EnableFrontend:  true,
			SeedDefaultUser: true,
		},
		Database: DatabaseConfig{
			Driver: "sqlite",
		},
		Postgres: PostgresConfig{
			Host:     "localhost",
			Port:     5432,
			User:     "postgres",
			Password: "",
			DBName:   "secretlane",
			SSLMode:  "disable",
		},
	}

	// Optional YAML config
	path := filepath.Join(".", "config.yaml")
	if data, err := os.ReadFile(path); err == nil {
		var fileCfg Config
		if err := yaml.Unmarshal(data, &fileCfg); err != nil {
			return fmt.Errorf("failed to parse config file: %w", err)
		}
		mergeConfig(&cfg, &fileCfg)
	}

	applyEnvOverrides(&cfg)

	// Expose runtime config
	App = AppRuntimeConfig{
		Port:            strconv.Itoa(cfg.App.Port),
		EnableFrontend:  cfg.App.EnableFrontend,
		SeedDefaultUser: cfg.App.SeedDefaultUser,
	}
	DBConfig = cfg.Postgres
	DBDriver = cfg.Database.Driver

	return nil
}

// mergeConfig overlays non-zero values from src onto dst.
func mergeConfig(dst, src *Config) {
	if src.App.Port != 0 {
		dst.App.Port = src.App.Port
	}
	// bools: only override when explicitly true in src
	if src.App.EnableFrontend {
		dst.App.EnableFrontend = true
	}
	if src.App.SeedDefaultUser {
		dst.App.SeedDefaultUser = true
	}

	if src.Postgres.Host != "" {
		dst.Postgres.Host = src.Postgres.Host
	}
	if src.Postgres.Port != 0 {
		dst.Postgres.Port = src.Postgres.Port
	}
	if src.Postgres.User != "" {
		dst.Postgres.User = src.Postgres.User
	}
	if src.Postgres.Password != "" {
		dst.Postgres.Password = src.Postgres.Password
	}
	if src.Postgres.DBName != "" {
		dst.Postgres.DBName = src.Postgres.DBName
	}
	if src.Postgres.SSLMode != "" {
		dst.Postgres.SSLMode = src.Postgres.SSLMode
	}

	if src.Database.Driver != "" {
		dst.Database.Driver = src.Database.Driver
	}
}

// applyEnvOverrides applies environment variables over the config.
func applyEnvOverrides(c *Config) {
	if v := os.Getenv("PORT"); v != "" {
		if port, err := strconv.Atoi(v); err == nil {
			c.App.Port = port
		}
	}
	if v := os.Getenv("ENABLE_FRONTEND"); v != "" {
		c.App.EnableFrontend = v == "true" || v == "1"
	}
	if v := os.Getenv("SEED_DEFAULT_USER"); v != "" {
		c.App.SeedDefaultUser = v == "true" || v == "1"
	}

	if v := os.Getenv("PGHOST"); v != "" {
		c.Postgres.Host = v
	}
	if v := os.Getenv("PGPORT"); v != "" {
		if port, err := strconv.Atoi(v); err == nil {
			c.Postgres.Port = port
		}
	}
	if v := os.Getenv("PGUSER"); v != "" {
		c.Postgres.User = v
	}
	if v := os.Getenv("PGPASSWORD"); v != "" {
		c.Postgres.Password = v
	}
	if v := os.Getenv("PGDATABASE"); v != "" {
		c.Postgres.DBName = v
	}
	if v := os.Getenv("PGSSLMODE"); v != "" {
		c.Postgres.SSLMode = v
	}

	if v := os.Getenv("DB_DRIVER"); v != "" {
		c.Database.Driver = v
	}
}
