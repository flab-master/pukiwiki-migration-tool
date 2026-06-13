package main

import (
	"database/sql"
	"log/slog"
	"net/http"
	"os"

	"pukiwiki-migration/internal/infra"
	"pukiwiki-migration/internal/infra/router"

	_ "modernc.org/sqlite"
)

func init() {
	level := slog.LevelInfo
	if v := os.Getenv("LOG_LEVEL"); v != "" {
		if err := level.UnmarshalText([]byte(v)); err != nil {
			slog.Warn("invalid LOG_LEVEL, using INFO", slog.String("value", v))
		}
	}
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level})))
}

func main() {
	cfg := infra.NewConfig()

	db, err := sql.Open("sqlite", cfg.DBPath)
	if err != nil {
		slog.Error("failed to open database", slog.String("error", err.Error()))
		return
	}
	defer db.Close()

	router, err := router.New(db, cfg)
	if err != nil {
		slog.Error("failed to create router", slog.String("error", err.Error()))
		return
	}

	const addr = ":8080"
	slog.Info("server starting", slog.String("addr", addr))
	if err := http.ListenAndServe(addr, router); err != nil {
		slog.Error("server failed", slog.String("error", err.Error()))
	}
}
