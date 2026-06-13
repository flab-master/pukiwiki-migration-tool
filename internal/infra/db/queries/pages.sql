-- name: UpsertPage :exec
INSERT OR IGNORE INTO pages (id, pukiwiki_user, pukiwiki_page, status, updated_at)
VALUES (?, ?, ?, 'pending', ?);

-- name: UpdatePageStatus :exec
UPDATE pages
SET status = ?, notion_id = ?, error_msg = ?, updated_at = ?
WHERE pukiwiki_user = ? AND pukiwiki_page = ?;

-- name: GetPendingPages :many
SELECT pukiwiki_page FROM pages
WHERE pukiwiki_user = ? AND status IN ('pending', 'failed');

-- name: GetStatusSummary :many
SELECT status, COUNT(*) AS count FROM pages
WHERE pukiwiki_user = ? GROUP BY status;
