-- +goose Up
-- +goose StatementBegin
DROP TABLE migrations;

CREATE TABLE pages (
     page_name  TEXT PRIMARY KEY,
     status     TEXT NOT NULL DEFAULT 'pending',
     notion_id  TEXT NOT NULL DEFAULT '',
     error_msg  TEXT NOT NULL DEFAULT '',
     updated_at TEXT NOT NULL
 );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE pages;
-- +goose StatementEnd
