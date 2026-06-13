-- name: CreateMigration :exec
INSERT INTO migrations (id, pukiwiki_user, notion_page_id, status, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?);

-- name: GetMigrationByID :one
SELECT id, pukiwiki_user, notion_page_id, status, created_at, updated_at
FROM migrations WHERE id = ?;

-- name: UpdateMigrationStatus :exec
UPDATE migrations
SET status = ?, updated_at = ?
WHERE id = ?;
