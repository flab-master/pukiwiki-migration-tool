# 開発ガイド

## Requirements

- [Go v1.25](https://go.dev/doc/install)
- [Task](https://taskfile.dev/docs/installation)
- [SQLite3](https://zenn.dev/enlog/articles/cc37c08f4b6d3f)
- [goose](https://github.com/pressly/goose)

## 環境変数

| 変数名                 | 説明                           | デフォルト              |
| ---------------------- | ------------------------------ | ----------------------- |
| `DB_PATH`              | SQLite DB ファイルパス         | `pukiwiki-migration.db` |
| `PUKIWIKI_BASE_URL`    | PukiWiki のベース URL          | —                       |
| `PUKIWIKI_USERNAME`    | PukiWiki のログインユーザー名  | —                       |
| `PUKIWIKI_PASSWORD`    | PukiWiki のログインパスワード  | —                       |
| `JWT_SECRET`           | JWT 署名シークレット           | —                       |
| `NOTION_API_TOKEN`     | Notion API トークン            | —                       |
| `CORS_ALLOWED_ORIGINS` | 許可する Origin (カンマ区切り) | —                       |

## 開発コマンド

```bash
# API サーバーを起動する
task dev

# API をビルドする
task build

# API の単体テスト実行
task test

# DB マイグレーション
task migrate:up
task migrate:down

# SQLite3 に接続する
sqlite3 pukiwiki-migration.db
```
