-- +goose Up
CREATE TABLE migrations (
    id             TEXT PRIMARY KEY,
    title          TEXT NOT NULL,
    markdown       TEXT NOT NULL,
    notion_page_id TEXT NOT NULL DEFAULT '',
    status         TEXT NOT NULL DEFAULT 'pending',
    created_at     TEXT NOT NULL,
    updated_at     TEXT NOT NULL
);

INSERT INTO migrations VALUES
    ('mig-001', 'FrontPage', '# FrontPage\n\nWelcome to the wiki.', '', 'pending', '2026-04-01T09:00:00Z', '2026-04-01T09:00:00Z'),
    ('mig-002', '技術メモ', '# 技術メモ\n\n- Go言語のメモ\n- Notionの使い方', 'notion-page-abc123', 'applied', '2026-04-02T10:00:00Z', '2026-04-10T14:30:00Z'),
    ('mig-003', '日記/2026-01', '# 日記 2026-01\n\n今月の振り返り。', 'notion-page-def456', 'accepted', '2026-04-03T11:00:00Z', '2026-04-12T16:00:00Z');

-- +goose Down
DROP TABLE migrations;
