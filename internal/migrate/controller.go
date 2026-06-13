package migrate

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
)

type usecaseInterface interface {
	start(ctx context.Context, user, notionPageID string) (migration, error)
	getStatus(ctx context.Context, id string) (migration, pageSummary, error)
}

func NewMigrationController(uc usecaseInterface) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /migrate", handleMigrate(uc))
	mux.HandleFunc("GET /migrate/{id}/status", handleStatus(uc))
	return mux
}

type migrateRequest struct {
	User         string `json:"user"`
	NotionPageID string `json:"notionPageId"`
}

type migrateResponse struct {
	ID string `json:"id"`
}

func handleMigrate(uc usecaseInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req migrateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.User == "" {
			http.Error(w, `{"error":"user is required"}`, http.StatusBadRequest)
			return
		}

		m, err := uc.start(r.Context(), req.User, req.NotionPageID)
		if err != nil {
			slog.Error("failed to start migration", "error", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(migrateResponse{ID: m.id})
	}
}

type statusResponse struct {
	ID      string      `json:"id"`
	User    string      `json:"user"`
	Status  string      `json:"status"`
	Summary pageSummary `json:"summary"`
}

func handleStatus(uc usecaseInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")

		m, summary, err := uc.getStatus(r.Context(), id)
		if err != nil {
			slog.Error("failed to get status", "error", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		resp := statusResponse{
			ID:      m.id,
			User:    m.pukiwikiUser,
			Status:  m.status,
			Summary: summary,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}
