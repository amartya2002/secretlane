package workspace

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/amartya2002/secretlane/internal/auth"
)

type Handler struct {
	service *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{service: s}
}

// /workspaces -> POST (create), GET (list)
func (h *Handler) Workspaces(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserID(r)

	switch r.Method {

	case http.MethodPost:
		var body struct {
			Name        string `json:"name"`
			Description string `json:"description"`
		}
		json.NewDecoder(r.Body).Decode(&body)

		id, err := h.service.Create(body.Name, body.Description, userID)
		if err != nil {
			log.Println("Workspace create error:", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error": err.Error(),
			})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(struct {
			ID int `json:"id"`
		}{ID: id})

	case http.MethodGet:
		list, err := h.service.ListForUser(userID)
		if err != nil {
			http.Error(w, "Failed to load workspaces", 500)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(list)

	default:
		http.Error(w, "Method not allowed", 405)
	}
}

// /workspaces/{id} -> PUT (update), DELETE (delete)
func (h *Handler) WorkspaceByID(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserID(r)

	idStr := r.URL.Path[len("/workspaces/"):]
	wsID, _ := strconv.Atoi(idStr)

	switch r.Method {

	case http.MethodPut:
		var body struct {
			Name        string `json:"name"`
			Description string `json:"description"`
		}
		json.NewDecoder(r.Body).Decode(&body)

		err := h.service.Update(wsID, body.Name, body.Description, userID)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error": err.Error(),
			})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message": "workspace updated",
		})

	case http.MethodDelete:
		err := h.service.Delete(wsID, userID)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error": err.Error(),
			})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message": "workspace deleted",
		})

	default:
		http.Error(w, "Method not allowed", 405)
	}
}
