# REST API 設計ドキュメント

PukiWiki の個人ページを Notion に移行するツールの API 設計

## 概要

PukiWiki からページデータを取得し、Notion API 経由で Notion ページとして移行する

## エンドポイント設計

エンドポイントは以下の２本です。

| メソッド | パス                  | 説明                                 |
| -------- | --------------------- | ------------------------------------ |
| `POST`   | `/api/migrate`        | 移行処理を開始し 202 Accepted を返す |
| `GET`    | `/api/migrate/status` | 移行処理の進捗状況を確認する         |

詳細は[こちら](./Endpoint.md)

## インフラ設計

- DB:
  - SQLite3
  - API と同一サーバー内にデータを保存する
  - 移行失敗時に再開できるように進捗状況を保存する

- MQ (非同期ジョブキュー):
  - 移行処理はリクエストを受け付けた後、非同期で実行される
  - スケールする必要は無し
  - Redis や RabbitMQ などの外部ジョブキューは利用しない
  - Go の channel + goroutine で捌く
  - PukiWiki からデータを取り出す際は、負荷を考慮する
    - `time.Ticker` を使用して少しづつ取り出す

## 処理フロー

### POST /api/migrate を受け取ったとき

```text
HTTP handler で以下の処理を実行
  ├─ goroutine 起動 (バックグラウンドで実行)
  └─ 202 Accepted を即返す

goroutine が以下の処理を捌く
  1. PukiWiki からページ一覧を取得
  2. 各ページ名を SQLite に INSERT OR IGNORE (done なページは触らない)
  3. SQLite から pending / failed のページ一覧を取得
  4. 各ページに対してループ:
       a. <-ticker.C             ← レート制限（PukiWiki 側の負荷軽減）
       b. PukiWiki からコンテンツ取得 (libpukiwiki の HTTP クライアントを使用)
       c. Notion API でページ作成
       d. SQLite のステータスを done または failed に更新
```

### GET /api/migrate/status を受け取った時

DB に保存されているジョブの進捗状況を取り出して返す

## DB 設計

**テーブルスキーマ**

```sql
CREATE TABLE pages (
    page_name  TEXT PRIMARY KEY,                 -- PukiWiki 側のページ名 (例: seminar-personal/morita2023/20240422)
    status     TEXT NOT NULL DEFAULT 'pending',  -- pending / done / failed
    notion_id  TEXT,                             -- Notion の PageID (このページ配下に PukiWiki と同様の構造のページ一覧が生成される)
    error_msg  TEXT,                             -- エラーメッセージ (記録用)
    updated_at TEXT NOT NULL                     -- レコードの更新日時 (記録用)
);
```

### ステータス遷移

```text
(pending) --(成功)--> (done)
        --(失敗)--> (failed)

(failed)  --(再実行)--> (pending) -> (done / failed)
(done)    --(再実行)--> (スキップ)（INSERT OR IGNORE）
```

### 再開のタイミング

`POST /api/migrate` を再実行したとき:

- `INSERT OR IGNORE` で新規ページのみ追加
- `done` のページはそのまま残る
- `pending` / `failed` のページだけ処理対象になる
