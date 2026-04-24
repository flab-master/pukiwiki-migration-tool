# REST API 設計ドキュメント

PukiWiki の個人ページを Notion に移行するツールの API 設計

## 概要

PukiWiki からページデータを取得し、Notion API 経由で Notion ページとして移行する。
移行はユーザー単位で行う。PukiWiki のページパス `seminar-personal/{user}/{page}` の `{user}` 配下を全て移行する。

## エンドポイント設計

エンドポイントは以下の3本です。

| メソッド | パス                         | 説明                                             |
| -------- | ---------------------------- | ------------------------------------------------ |
| `GET`    | `/health`                    | API プロセスのヘルスチェック                     |
| `POST`   | `/api/migrate`               | ユーザーの移行をキューに積み 202 Accepted を返す |
| `GET`    | `/api/migrate/{user}/status` | 指定ユーザーの移行進捗を確認する                 |

詳細は[こちら](./Endpoint.md)

## インフラ設計

- DB:
  - SQLite3
  - API と同一サーバー内にデータを保存する
  - 移行失敗時に再開できるように進捗状況を保存する

- 非同期ジョブキュー:
  - 移行処理はリクエストを受け付けた後、非同期で実行される
  - スケールする必要は無し
  - Redis や RabbitMQ などの外部ジョブキューは利用しない
  - Go の buffered channel + goroutine で捌く
  - worker goroutine 1本がキューを順番に処理する (ユーザー単位で直列)
  - PukiWiki からデータを取り出す際は、負荷を考慮する
    - `time.Ticker` を使用して少しづつ取り出す

## 処理フロー

### POST /api/migrate を受け取ったとき

```text
HTTP handler
  ├─ リクエストボディから user を取得
  ├─ queue (buffered channel) に user を積む
  └─ 202 Accepted を即返す

worker goroutine (サーバー起動時から常駐)
  <- queue から user を取り出す
  1. PukiWiki から seminar-personal/{user}/ 配下のページ一覧を取得
  2. 各ページ名を SQLite に INSERT OR IGNORE (done なページは触らない)
  3. SQLite から pending / failed のページ一覧を取得
  4. 各ページに対してループ:
       a. <-ticker.C             ← レート制限（PukiWiki 側の負荷軽減）
       b. PukiWiki からコンテンツ取得
       c. Notion API でページ作成
       d. SQLite のステータスを done または failed に更新
```

### GET /api/migrate/{user}/status を受け取った時

DB に保存されている指定ユーザーの進捗状況を集計して返す

## DB 設計

**テーブルスキーマ**

```sql
CREATE TABLE pages (
    user       TEXT NOT NULL,                      -- PukiWiki 側のユーザー名 (例: morita2023)
    page_name  TEXT NOT NULL,                      -- PukiWiki 側のページ名 (例: seminar-personal/morita2023/20240422)
    status     TEXT NOT NULL DEFAULT 'pending',    -- pending / done / failed
    notion_id  TEXT NOT NULL DEFAULT '',           -- Notion の PageID
    error_msg  TEXT NOT NULL DEFAULT '',           -- エラーメッセージ (記録用)
    updated_at TEXT NOT NULL,                      -- レコードの更新日時 (記録用)
    PRIMARY KEY (user, page_name)
);
```

### ステータス遷移

```text
(pending) --(成功)--> (done)
          --(失敗)--> (failed)

(failed)  --(再実行)--> pending -> (done / failed)
(done)    --(再実行)--> スキップ（INSERT OR IGNORE）
```

### 再開のタイミング

`POST /api/migrate` を再実行したとき:

- `INSERT OR IGNORE` で新規ページのみ追加
- `done` のページはそのまま残る
- `pending` / `failed` のページだけ処理対象になる
