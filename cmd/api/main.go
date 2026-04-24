package main

import (
	"database/sql"
	"log/slog"
	"net/http"
	"os"

	"pukiwiki-migration/internal"

	_ "modernc.org/sqlite"
)

func init() {
	// ロガーを初期化
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, nil)))
}

func main() {

	// 環境変数から接続情報を取得する
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "pukiwiki-migration.db"
	}
	pukiBaseURL := os.Getenv("PUKIWIKI_BASE_URL")
	pukiUsername := os.Getenv("PUKIWIKI_USERNAME")
	pukiPassword := os.Getenv("PUKIWIKI_PASSWORD")
	// TODO: NOTION_API_TOKEN := os.Getenv("NOTION_API_TOKEN")

	// DB へ接続する
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		slog.Error("failed to open database", slog.String("error", err.Error()))
		return
	}
	defer db.Close()

	pmu := internal.NewPageMigrator(db, pukiBaseURL, pukiUsername, pukiPassword)

	// HTTP ハンドラーのルーティング設定
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/migrate", internal.HandleMigrate(pmu))
	mux.HandleFunc("GET /api/migrate/{user}/status", internal.HandleStatus(pmu))

	// API サーバーを起動する
	addr := ":8080"
	slog.Info("server starting", slog.String("addr", addr))
	if err := http.ListenAndServe(addr, mux); err != nil {
		slog.Error("server failed", slog.String("error", err.Error()))
	}
}
