# API JSON スキーマ

**スキーマは仮なので後から色々変更になる可能性があります**

---

Base URL: `http://localhost:8080`

## 1. 移行一覧取得

エンドポイント: `GET /api/migration/list`

### レスポンス形式

ステータス: `200 OK`

```json
{
  "migrations": [
    {
      "id": "mig-001",
      "title": "ページタイトル",
      "markdown": "## 本文...",
      "notion_page_id": "notion-page-mock-mig-001",
      "status": "pending",
      "created_at": "2026-01-01T00:00:00Z",
      "updated_at": "2026-01-01T00:00:00Z"
    }
  ]
}
```

## 2. 移行申請

エンドポイント: `POST /api/migration/apply`

### リクエスト

```json
{
  "id": "mig-001"
}
```

### レスポンス

**201 Created**

```json
{
  "migration": {
    "id": "mig-001",
    "title": "ページタイトル",
    "markdown": "## 本文...",
    "notion_page_id": "notion-page-mock-mig-001",
    "status": "applied",
    "created_at": "2026-01-01T00:00:00Z",
    "updated_at": "2026-01-01T12:00:00Z"
  }
}
```

**400 Bad Request** (リクエストボディの形式が不正)

```json
{ "error": "invalid request body" }
{ "error": "id is required" }
```

**404 Not Found** — 指定 ID が存在しない

```json
{ "error": "migration not found" }
```

**409 Conflict** — ステータスが `pending` 以外

```json
{ "error": "migration must be in pending status" }
```

## 3. 移行承認

エンドポイント: `POST /api/migration/accept`

### リクエスト

```json
{
  "id": "mig-002"
}
```

### レスポンス

**200 OK**

```json
{
  "migration": {
    "id": "mig-002",
    "title": "ページタイトル",
    "markdown": "## 本文...",
    "notion_page_id": "notion-page-mock-mig-002",
    "status": "accepted",
    "created_at": "2026-01-01T00:00:00Z",
    "updated_at": "2026-01-01T13:00:00Z"
  }
}
```

**400 Bad Request** — リクエストボディ不正 / `id` が空

```json
{ "error": "invalid request body" }
{ "error": "id is required" }
```

**404 Not Found** — 指定IDが存在しない

```json
{ "error": "migration not found" }
```

**409 Conflict** — ステータスが `applied` 以外

```json
{ "error": "migration must be in applied status" }
```
