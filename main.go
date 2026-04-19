package main

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"
)

type MigrationStatus string

const (
	StatusPending  MigrationStatus = "pending"
	StatusApplied  MigrationStatus = "applied"
	StatusAccepted MigrationStatus = "accepted"
)

type Migration struct {
	ID           string          `json:"id"`
	Title        string          `json:"title"`
	Markdown     string          `json:"markdown"`
	NotionPageID string          `json:"notion_page_id,omitempty"`
	Status       MigrationStatus `json:"status"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
}

var mockMigrations = []Migration{
	{
		ID:        "mig-001",
		Title:     "FrontPage",
		Markdown:  "# FrontPage\n\nWelcome to the wiki.",
		Status:    StatusPending,
		CreatedAt: time.Date(2026, 4, 1, 9, 0, 0, 0, time.UTC),
		UpdatedAt: time.Date(2026, 4, 1, 9, 0, 0, 0, time.UTC),
	},
	{
		ID:           "mig-002",
		Title:        "技術メモ",
		Markdown:     "# 技術メモ\n\n- Go言語のメモ\n- Notionの使い方",
		NotionPageID: "notion-page-abc123",
		Status:       StatusApplied,
		CreatedAt:    time.Date(2026, 4, 2, 10, 0, 0, 0, time.UTC),
		UpdatedAt:    time.Date(2026, 4, 10, 14, 30, 0, 0, time.UTC),
	},
	{
		ID:           "mig-003",
		Title:        "日記/2026-01",
		Markdown:     "# 日記 2026-01\n\n今月の振り返り。",
		NotionPageID: "notion-page-def456",
		Status:       StatusAccepted,
		CreatedAt:    time.Date(2026, 4, 3, 11, 0, 0, 0, time.UTC),
		UpdatedAt:    time.Date(2026, 4, 12, 16, 0, 0, 0, time.UTC),
	},
}

func listMigrationGet(w http.ResponseWriter, r *http.Request) {
	resp := map[string]any{"migrations": mockMigrations}

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

func applyMigrationPost(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}
	if req.ID == "" {
		http.Error(w, `{"error":"id is required"}`, http.StatusBadRequest)
		return
	}

	for i := range mockMigrations {
		if mockMigrations[i].ID != req.ID {
			continue
		}
		if mockMigrations[i].Status != StatusPending {
			http.Error(w, `{"error":"migration must be in pending status"}`, http.StatusConflict)
			return
		}
		mockMigrations[i].Status = StatusApplied
		mockMigrations[i].NotionPageID = "notion-page-mock-" + req.ID
		mockMigrations[i].UpdatedAt = time.Now().UTC()

		resp := map[string]any{"migration": mockMigrations[i]}
		var buf bytes.Buffer
		if err := json.NewEncoder(&buf).Encode(resp); err != nil {
			slog.Error("failed to encode response", slog.String("error", err.Error()))
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(buf.Bytes())
		return
	}

	http.Error(w, `{"error":"migration not found"}`, http.StatusNotFound)
}

func acceptMigrationPost(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}
	if req.ID == "" {
		http.Error(w, `{"error":"id is required"}`, http.StatusBadRequest)
		return
	}

	for i := range mockMigrations {
		if mockMigrations[i].ID != req.ID {
			continue
		}
		if mockMigrations[i].Status != StatusApplied {
			http.Error(w, `{"error":"migration must be in applied status"}`, http.StatusConflict)
			return
		}
		mockMigrations[i].Status = StatusAccepted
		mockMigrations[i].UpdatedAt = time.Now().UTC()

		resp := map[string]any{"migration": mockMigrations[i]}
		var buf bytes.Buffer
		if err := json.NewEncoder(&buf).Encode(resp); err != nil {
			slog.Error("failed to encode response", slog.String("error", err.Error()))
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(buf.Bytes())
		return
	}

	http.Error(w, `{"error":"migration not found"}`, http.StatusNotFound)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/migration/list", listMigrationGet)
	mux.HandleFunc("POST /api/migration/apply", applyMigrationPost)
	mux.HandleFunc("POST /api/migration/accept", acceptMigrationPost)

	addr := ":8080"
	slog.Info("server starting", slog.String("addr", addr))
	if err := http.ListenAndServe(addr, mux); err != nil {
		slog.Error("server failed", slog.String("error", err.Error()))
	}
}
