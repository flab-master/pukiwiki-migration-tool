PukiWiki FLab の個人ページを Notion に移行するツール

## Guide

**goose DB マイグレーション実行**

```bash
goose -dir db/migrations sqlite3 pukiwiki-migration.db up
```

**Sqlite3 接続**

```bash
sqlite3 pukiwiki-migration.db
```

**API 起動**

```bash
task run
```

**API ビルド**

```bash
task build
```

**テスト実行**

```bash
task test
```

**API テスト**

```bash
curl -s http://localhost:8080/api/migration/list | jq .
```

```bash
curl -s -X POST http://localhost:8080/api/migration/apply \
 -H "Content-Type: application/json" \
 -d '{"id":"mig-001"}' | jq .
```

```bash
curl -s -X POST http://localhost:8080/api/migration/accept \
 -H "Content-Type: application/json" \
 -d '{"id":"mig-002"}' | jq .
```
