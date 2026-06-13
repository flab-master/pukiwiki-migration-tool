package router

import (
	"database/sql"
	"net/http"

	"pukiwiki-migration/internal/auth"
	"pukiwiki-migration/internal/infra"
	"pukiwiki-migration/internal/infra/router/middlewares"
	"pukiwiki-migration/internal/migrate"
)

func New(
	db *sql.DB,
	cfg *infra.Config,
) (http.Handler, error) {
	uc, err := migrate.NewMigrationUsecase(db, cfg.PukiBaseURL, cfg.PukiUsername, cfg.PukiPassword, cfg.NotionToken)
	if err != nil {
		return nil, err
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", handleHealth())
	mux.HandleFunc("POST /api/auth/login", auth.HandleLogin(cfg.PukiUsername, cfg.PukiPassword, cfg.JWTSecret))

	authMW := middlewares.JWTMiddleware(cfg.JWTSecret)
	authedMigrate := authMW(http.StripPrefix("/api", migrate.NewMigrationController(uc)))
	mux.Handle("/api/migrate", authedMigrate)
	mux.Handle("/api/migrate/", authedMigrate)

	return middlewares.CORSMiddleware(cfg.AllowedOrigins, mux), nil
}
