# API エンドポイント仕様

Base URL: `http://localhost:8080`

## POST /api/migrate

ユーザーの移行をキューに積む。処理はバックグラウンドで順番に実行され、即座に 202 を返す。
再実行した場合、`done` 済みのページはスキップし `pending` / `failed` のページのみ再処理する。

### リクエスト

```json
{
  "user": "morita2023",
  "notionPageId": "djalfja...uuid"
}
```

| フィールド     | 説明                                               |
| -------------- | -------------------------------------------------- |
| `user`         | PukiWiki のユーザー名（ページパスの第2セグメント） |
| `notionPageId` | Notion の PageId                                   |

### レスポンス

**202 Accepted** — キューに積んだ

```json
{}
```

**400 Bad Request** — `user` が空またはリクエストボディ不正

```json
{ "error": "user is required" }
```

### cURL 例

```bash
curl -X POST http://localhost:8080/api/migrate \
    -H "Content-Type: application/json" \
    -d '{"user":"morita2023"}'
```

---

## GET /api/migrate/{user}/status

指定ユーザーの移行進捗を取得する。

### レスポンス

**200 OK**

```json
{
  "user": "morita2023",
  "running": true,
  "total": 30,
  "done": 12,
  "failed": 1,
  "pending": 17
}
```

| フィールド | 説明                                      |
| ---------- | ----------------------------------------- |
| `user`     | ユーザー名                                |
| `running`  | worker が現在このユーザーを処理中かどうか |
| `total`    | 対象ページ総数                            |
| `done`     | 移行完了数                                |
| `failed`   | 失敗数（再実行で再処理される）            |
| `pending`  | 未処理数                                  |

**404 Not Found** — 指定ユーザーの移行レコードが存在しない

```json
{ "error": "user not found" }
```

### cURL 例

```bash
curl http://localhost:8080/api/migrate/morita2023/status
```
