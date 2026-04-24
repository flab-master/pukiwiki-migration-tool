# PukiWiki Migration Tool

PukiWiki FLab の個人ページを Notion に移行するツール

## 必要なものをインストールする

- [Go v1.25](https://go.dev/doc/install)
- [Task](https://taskfile.dev/docs/installation)
- [SQLite3](https://zenn.dev/enlog/articles/cc37c08f4b6d3f)
- [goose](https://github.com/pressly/goose)

## 開発コマンド

**goose を使って DB マイグレーションを実行する**

```bash
goose -dir db/migrations sqlite3 pukiwiki-migration.db up
```

**Sqlite3 に接続する**

```bash
sqlite3 pukiwiki-migration.db
```

**API を起動する**

```bash
go run ./cmd/

# Task を使う場合はこっち
task run
```

**API をビルドする**

```bash
go build -o pukiwiki-migration ./cmd/

# Task を使う場合はこっち
task build
```

**テストを実行する**

```bash
go test -v ./...

# Task を使う場合はこっち
task test
```

## API の動作確認

**[API設計](./docs/API-Design.md)**
**[エンドポイントの仕様](./docs/Endpoint.md)**

以下は cURL コマンドを使った例 (Postman でも可)

**移行を開始する**

```bash
curl -X POST http://localhost:8080/api/migrate \
    -H "Content-Type: application/json" \
    -d '{"user":"morita2023"}'
```

**移行の進捗を確認する**

```bash
curl http://localhost:8080/api/migrate/morita2023/status
```
