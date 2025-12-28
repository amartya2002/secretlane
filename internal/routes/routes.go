package routes

import (
	"net/http"

	"github.com/amartya2002/secretlane/internal/auth"
	"github.com/amartya2002/secretlane/internal/config"
	"github.com/amartya2002/secretlane/internal/workspace"
)

func SetupRoutes(mux *http.ServeMux, authService *auth.AuthService, wsService *workspace.Service) {
	authHandler := auth.NewLoginHandler(authService)
	wsHandler := workspace.NewHandler(wsService)

	const apiV1 = "/api/v1"

	// Auth
	mux.HandleFunc(apiV1+"/signup", authHandler.Signup)
	mux.HandleFunc(apiV1+"/login", authHandler.Login)
	mux.Handle(apiV1+"/logout", auth.RequireAuth(http.HandlerFunc(auth.Logout)))

	// Health
	mux.HandleFunc(apiV1+"/healthz", config.HealthCheckHandler)

	// Workspaces (authenticated)
	mux.Handle(apiV1+"/workspaces", auth.RequireAuth(http.HandlerFunc(wsHandler.Workspaces)))
	mux.Handle(apiV1+"/workspaces/", auth.RequireAuth(http.HandlerFunc(wsHandler.WorkspaceByID)))
}
