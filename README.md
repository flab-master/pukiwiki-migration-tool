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
go run main.go

# Task を使う場合はこっち
task run
```

**API をビルドする**

```bash
go build -o pukiwiki-migration .

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

**[エンドポイントの仕様](./docs/endpoints.md)**

以下は cURL コマンドを使った例 (Postman でも可)

**移行一覧を取得する**

```bash
curl http://localhost:8080/api/migration/list
```

**移行申請を行う**

```bash
curl -X POST http://localhost:8080/api/migration/apply \
    -H "Content-Type: application/json" \
    -d '{"id":"mig-001"}'
```

**移行を承認する**

```bash
curl -X POST http://localhost:8080/api/migration/accept \
    -H "Content-Type: application/json" \
    -d '{"id":"mig-002"}' | jq .
```
