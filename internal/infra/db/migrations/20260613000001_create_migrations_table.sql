-- +goose Up
CREATE TABLE migrations (
    id              TEXT PRIMARY KEY,
    pukiwiki_user   TEXT NOT NULL,
    notion_page_id  TEXT NOT NULL,
    status          TEXT NOT NULL DEFAULT 'pending',
    created_at      TEXT NOT NULL,
    updated_at      TEXT NOT NULL
);

-- +goose Down
DROP TABLE migrations;
