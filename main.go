package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	_ "modernc.org/sqlite"
)

func main() {
	db, err := sql.Open("sqlite", "pukiwiki-migration.db")
	if err != nil {
		slog.Error("failed to open database", slog.String("error", err.Error()))
		return
	}
	defer db.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/migration/list", listMigrationGet(db))
	mux.HandleFunc("POST /api/migration/apply", applyMigrationPost(db))
	mux.HandleFunc("POST /api/migration/accept", acceptMigrationPost(db))

	addr := ":8080"
	slog.Info("server starting", slog.String("addr", addr))
	if err := http.ListenAndServe(addr, mux); err != nil {
		slog.Error("server failed", slog.String("error", err.Error()))
	}
}

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

func scanMigration(row *sql.Row) (Migration, error) {
	var m Migration
	var createdAt, updatedAt string
	if err := row.Scan(&m.ID, &m.Title, &m.Markdown, &m.NotionPageID, &m.Status, &createdAt, &updatedAt); err != nil {
		return Migration{}, err
	}
	m.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	m.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)
	return m, nil
}

func listMigrationGet(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT id, title, markdown, notion_page_id, status, created_at, updated_at FROM migrations")
		if err != nil {
			slog.Error("failed to query migrations", slog.String("error", err.Error()))
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		list := []Migration{}
		for rows.Next() {
			var m Migration
			var createdAt, updatedAt string
			if err := rows.Scan(&m.ID, &m.Title, &m.Markdown, &m.NotionPageID, &m.Status, &createdAt, &updatedAt); err != nil {
				slog.Error("failed to scan migration", slog.String("error", err.Error()))
				http.Error(w, "internal server error", http.StatusInternalServerError)
				return
			}
			m.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
			m.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)
			list = append(list, m)
		}

		resp := map[string]any{"migrations": list}
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

func applyMigrationPost(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		var status MigrationStatus
		err := db.QueryRow("SELECT status FROM migrations WHERE id = ?", req.ID).Scan(&status)
		if err == sql.ErrNoRows {
			http.Error(w, `{"error":"migration not found"}`, http.StatusNotFound)
			return
		}
		if err != nil {
			slog.Error("failed to query migration", slog.String("error", err.Error()))
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		if status != StatusPending {
			http.Error(w, `{"error":"migration must be in pending status"}`, http.StatusConflict)
			return
		}

		notionPageID := "notion-page-mock-" + req.ID
		updatedAt := time.Now().UTC().Format(time.RFC3339)
		_, err = db.Exec("UPDATE migrations SET status='applied', notion_page_id=?, updated_at=? WHERE id=?", notionPageID, updatedAt, req.ID)
		if err != nil {
			slog.Error("failed to update migration", slog.String("error", err.Error()))
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		m, err := scanMigration(db.QueryRow("SELECT id, title, markdown, notion_page_id, status, created_at, updated_at FROM migrations WHERE id = ?", req.ID))
		if err != nil {
			slog.Error("failed to fetch updated migration", slog.String("error", err.Error()))
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		var buf bytes.Buffer
		if err := json.NewEncoder(&buf).Encode(map[string]any{"migration": m}); err != nil {
			slog.Error("failed to encode response", slog.String("error", err.Error()))
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(buf.Bytes())
	}
}

func acceptMigrationPost(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		var status MigrationStatus
		err := db.QueryRow("SELECT status FROM migrations WHERE id = ?", req.ID).Scan(&status)
		if err == sql.ErrNoRows {
			http.Error(w, `{"error":"migration not found"}`, http.StatusNotFound)
			return
		}
		if err != nil {
			slog.Error("failed to query migration", slog.String("error", err.Error()))
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		if status != StatusApplied {
			http.Error(w, `{"error":"migration must be in applied status"}`, http.StatusConflict)
			return
		}

		updatedAt := time.Now().UTC().Format(time.RFC3339)
		_, err = db.Exec("UPDATE migrations SET status='accepted', updated_at=? WHERE id=?", updatedAt, req.ID)
		if err != nil {
			slog.Error("failed to update migration", slog.String("error", err.Error()))
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		m, err := scanMigration(db.QueryRow("SELECT id, title, markdown, notion_page_id, status, created_at, updated_at FROM migrations WHERE id = ?", req.ID))
		if err != nil {
			slog.Error("failed to fetch updated migration", slog.String("error", err.Error()))
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		var buf bytes.Buffer
		if err := json.NewEncoder(&buf).Encode(map[string]any{"migration": m}); err != nil {
			slog.Error("failed to encode response", slog.String("error", err.Error()))
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(buf.Bytes())
	}
}
