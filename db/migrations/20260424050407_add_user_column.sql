-- +goose Up
DROP TABLE pages;
CREATE TABLE pages (
    user       TEXT NOT NULL,
    page_name  TEXT NOT NULL,
    status     TEXT NOT NULL DEFAULT 'pending',
    notion_id  TEXT NOT NULL DEFAULT '',
    error_msg  TEXT NOT NULL DEFAULT '',
    updated_at TEXT NOT NULL,
    PRIMARY KEY (user, page_name)
);

-- +goose Down
DROP TABLE pages;
CREATE TABLE pages (
    page_name  TEXT PRIMARY KEY,
    status     TEXT NOT NULL DEFAULT 'pending',
    notion_id  TEXT NOT NULL DEFAULT '',
    error_msg  TEXT NOT NULL DEFAULT '',
    updated_at TEXT NOT NULL
);
