package internal

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
)

type migrationRequest struct {
	User         string `json:"user"`
	NotionPageId string `json:"notionPageId"`
}

func HandleMigrate(u *pageMigrator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req migrationRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.User == "" {
			http.Error(w, `{"error":"user is required"}`, http.StatusBadRequest)
			return
		}

		u.enqueue(req.User, req.NotionPageId)

		w.WriteHeader(http.StatusAccepted)
	}
}

type statusResponse struct {
	User    string `json:"user"`
	Running bool   `json:"running"`
	statusSummary
}

func HandleStatus(u *pageMigrator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.PathValue("user")

		summary, err := getStatusSummary(u.db, user)
		if err != nil {
			slog.Error("failed to get status summary", slog.String("error", err.Error()))
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		if summary.Total == 0 {
			http.Error(w, `{"error":"user not found"}`, http.StatusNotFound)
			return
		}

		resp := statusResponse{
			User:          user,
			Running:       u.isMigrating(user),
			statusSummary: summary,
		}

		var buf bytes.Buffer
		if err := json.NewEncoder(&buf).Encode(resp); err != nil {
			slog.Error("failed to encode response", slog.String("error", err.Error()))
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(buf.Bytes())
	}
}
