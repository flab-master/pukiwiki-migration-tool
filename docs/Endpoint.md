# API エンドポイント仕様

Base URL: `http://localhost:8080`

## POST /api/migrate

移行を開始する。処理はバックグラウンドで実行され、即座に 202 を返す。  
再実行した場合、`done` 済みのページはスキップし `pending` / `failed` のページのみ再処理する。

### レスポンス

**202 Accepted** — 移行開始

```json
{}
```

**409 Conflict** — すでに移行が実行中

```json
{ "error": "migration is already running" }
```

### cURL 例

```bash
curl -X POST http://localhost:8080/api/migrate
```

---

## GET /api/migrate/status

移行の進捗を取得する。

### レスポンス

**200 OK**

```json
{
  "running": true,
  "total": 100,
  "done": 42,
  "failed": 3,
  "pending": 55
}
```

| フィールド | 説明                                 |
| ---------- | ------------------------------------ |
| `running`  | バックグラウンド処理が実行中かどうか |
| `total`    | 対象ページ総数                       |
| `done`     | 移行完了数                           |
| `failed`   | 失敗数（再実行で再処理される）       |
| `pending`  | 未処理数                             |

### cURL 例

```bash
curl http://localhost:8080/api/migrate/status
```
