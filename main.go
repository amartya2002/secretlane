package main

import (
	"log"
	"net/http"

	"github.com/amartya2002/secretlane/internal/auth"
	"github.com/amartya2002/secretlane/internal/config"
	"github.com/amartya2002/secretlane/internal/middleware"
	"github.com/amartya2002/secretlane/internal/routes"
	"github.com/amartya2002/secretlane/internal/workspace"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env first so config can pick up env overrides.
	_ = godotenv.Load()

	if err := config.LoadAppConfig(); err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	config.InitDatabase()
	config.RunMigrations()
	auth.InitJWT()

	authService := auth.NewAuthService()
	wsService := workspace.NewService()

	mux := http.NewServeMux()

	routes.SetupRoutes(mux, authService, wsService)

	handler := middleware.CORS(mux)
	log.Printf("server running :%s", config.App.Port)
	err := http.ListenAndServe(":"+config.App.Port, handler)
	if err != nil {
		return
	}
}
