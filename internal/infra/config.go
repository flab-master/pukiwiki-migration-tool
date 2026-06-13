package infra

import (
	"os"
	"strings"
)

type Config struct {
	DBPath         string
	PukiBaseURL    string
	PukiUsername   string
	PukiPassword   string
	JWTSecret      string
	NotionToken    string
	AllowedOrigins []string
}

func NewConfig() *Config {
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "pukiwiki-migration.db"
	}

	var allowedOrigins []string
	if raw := os.Getenv("CORS_ALLOWED_ORIGINS"); raw != "" {
		for o := range strings.SplitSeq(raw, ",") {
			allowedOrigins = append(allowedOrigins, strings.TrimSpace(o))
		}
	}

	return &Config{
		DBPath:         dbPath,
		PukiBaseURL:    os.Getenv("PUKIWIKI_BASE_URL"),
		PukiUsername:   os.Getenv("PUKIWIKI_USERNAME"),
		PukiPassword:   os.Getenv("PUKIWIKI_PASSWORD"),
		JWTSecret:      os.Getenv("JWT_SECRET"),
		NotionToken:    os.Getenv("NOTION_API_TOKEN"),
		AllowedOrigins: allowedOrigins,
	}
}
