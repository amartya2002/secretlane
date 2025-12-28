package config

import (
	"encoding/json"
	"net/http"
	"time"
)

type HealthStatus struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
}

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {

	// --- DB check disabled for now ---
	// dbErr := DB.Ping()
	// if dbErr != nil {
	//     w.WriteHeader(http.StatusServiceUnavailable)
	//     json.NewEncoder(w).Encode(HealthStatus{
	//         Status:    "unhealthy",
	//         Timestamp: time.Now().UTC().Format(time.RFC3339),
	//     })
	//     return
	// }

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(HealthStatus{
		Status:    "healthy",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	})
}
