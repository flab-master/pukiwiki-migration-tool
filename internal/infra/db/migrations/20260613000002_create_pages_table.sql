-- +goose Up
CREATE TABLE pages (
    id              TEXT PRIMARY KEY,
    pukiwiki_user   TEXT NOT NULL,
    pukiwiki_page   TEXT NOT NULL,
    status          TEXT NOT NULL DEFAULT 'pending',
    notion_id       TEXT NOT NULL DEFAULT '',
    error_msg       TEXT NOT NULL DEFAULT '',
    updated_at      TEXT NOT NULL,
    UNIQUE (pukiwiki_user, pukiwiki_page)
);

-- +goose Down
DROP TABLE pages;
