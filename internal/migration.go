package internal

import (
	"database/sql"
	"log/slog"
	"sync/atomic"
)

type PageMigrationUsecase struct {
	db          *sql.DB
	queue       chan string
	currentUser atomic.Pointer[string]
}

func NewPageMigrationUsecase(db *sql.DB) *PageMigrationUsecase {
	u := &PageMigrationUsecase{
		db:    db,
		queue: make(chan string, 100),
	}

	go u.worker()

	return u
}

func (u *PageMigrationUsecase) worker() {
	for user := range u.queue {
		u.currentUser.Store(&user)

		slog.Info("migration starting", slog.String("user", user))

		// TODO: 1. PukiWiki から seminar-personal/{user}/ 配下のページ一覧を取得

		// TODO: 2. upsertPages(db, user, pageNames) で SQLite に登録

		// TODO: 3. pending/failed ページを取得

		// TODO: 4. ticker でレート制限しながら各ページを処理
		//          a. PukiWiki からコンテンツ取得
		//          b. Notion API でページ作成
		//          c. updatePageStatus(db, user, pageName, ...) で結果を記録

		u.currentUser.Store(nil)
	}
}

func (u *PageMigrationUsecase) enqueue(user string) {
	u.queue <- user
}

func (u *PageMigrationUsecase) isMigrating(user string) bool {
	p := u.currentUser.Load()
	return p != nil && *p == user
}
