# API エンドポイント仕様

Base URL: `http://localhost:8080`

## 認証

- `/api/migrate` 以下のエンドポイントは JWT 認証が必要
- `Authorization: Bearer <token>` ヘッダーを付与する

## GET /health

API プロセスが起動していることを確認する。

### レスポンス

**200 OK**

```json
{ "status": "ok" }
```

### cURL 例

```bash
curl http://localhost:8080/health
```

---

## POST /api/auth/login

PukiWiki の認証情報で JWT を取得する。

### リクエスト

```json
{
  "username": "admin",
  "password": "secret"
}
```

### レスポンス

**200 OK**

```json
{ "token": "<jwt>" }
```

**401 Unauthorized** — 認証情報が不正

```json
{ "error": "invalid credentials" }
```

### cURL 例

```bash
curl -X POST http://localhost:8080/api/auth/login \
    -H "Content-Type: application/json" \
    -d '{"username":"admin", "password":"secret"}'
```

---

## POST /api/migrate

ユーザーの移行ジョブを実行する

### リクエスト

```json
{
  "user": "morit958",
  "notionPageId": "8723b73c-..."
}
```

| フィールド     | 説明                    |
| -------------- | ----------------------- |
| `user`         | PukiWiki のユーザー名   |
| `notionPageId` | 移行先 Notion の PageId |

### レスポンス

**202 Accepted**

```json
{ "id": "550e8400-e29b-41d4-a716-446655440000" }
```

| フィールド | 説明                                       |
| ---------- | ------------------------------------------ |
| `id`       | 移行ジョブ ID。進捗確認 (`/status`) で使う |

**400 Bad Request** — `user` が空またはリクエストボディ不正

```json
{ "error": "user is required" }
```

### cURL 例

```bash
curl -X POST http://localhost:8080/api/migrate \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer <JWT>" \
    -d '{"user":"morit958", "notionPageId":"8723b73c-487b-7427-b497-9f5bd58ff974"}'
```

---

## GET /api/migrate/{id}/status

移行ジョブの進捗を取得する。`{id}` は `POST /api/migrate` で返ってきた `id`。

### レスポンス

**200 OK**

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "user": "morit958",
  "status": "running",
  "summary": {
    "total": 30,
    "done": 12,
    "failed": 1,
    "pending": 17
  }
}
```

| フィールド        | 説明                                                           |
| ----------------- | -------------------------------------------------------------- |
| `id`              | 移行ジョブ ID                                                  |
| `user`            | PukiWiki ユーザー名                                            |
| `status`          | ジョブのステータス (`pending` / `running` / `done` / `failed`) |
| `summary.total`   | 対象ページ総数                                                 |
| `summary.done`    | 移行完了数                                                     |
| `summary.failed`  | 失敗数                                                         |
| `summary.pending` | 未処理数                                                       |

### cURL 例

```bash
curl -H "Authorization: Bearer $TOKEN" \
    http://localhost:8080/api/migrate/550e8400-e29b-41d4-a716-446655440000/status
```
