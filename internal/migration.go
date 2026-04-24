package internal

import (
	"database/sql"
	"fmt"
	"log/slog"
	"sync/atomic"

	"github.com/moriT958/libpukiwiki"
)

type migrationJob struct {
	user         string
	notionPageId string
}

type pageMigrator struct {
	db           *sql.DB
	queue        chan migrationJob
	errCh        chan error
	currentUser  atomic.Pointer[string]
	pukiBaseURL  string
	pukiUsername string
	pukiPassword string
}

func NewPageMigrator(db *sql.DB, pukiBaseURL, pukiUsername, pukiPassword string) *pageMigrator {
	u := &pageMigrator{
		db:           db,
		queue:        make(chan migrationJob, 100),
		errCh:        make(chan error, 10),
		pukiBaseURL:  pukiBaseURL,
		pukiUsername: pukiUsername,
		pukiPassword: pukiPassword,
	}

	go u.worker()
	go u.handleErrors()

	return u
}

const pukiwikiPersonalScope = "seminar-personal/"

func (pm *pageMigrator) worker() {
	for job := range pm.queue {
		pm.currentUser.Store(&job.user)
		slog.Info("migration starting", slog.String("user", job.user))

		if err := pm.migrate(job); err != nil {
			pm.errCh <- fmt.Errorf("user %s: %w", job.user, err)
			pm.currentUser.Store(nil)
			return
		}

		pm.currentUser.Store(nil)

		slog.Info("migration finished", slog.String("user", job.user))
	}
}

func (pm *pageMigrator) migrate(job migrationJob) error {
	client, err := libpukiwiki.NewClient(pm.pukiBaseURL,
		libpukiwiki.WithAuth(pm.pukiUsername, pm.pukiPassword),
		libpukiwiki.WithScope(pukiwikiPersonalScope+job.user),
	)
	if err != nil {
		return fmt.Errorf("create client: %w", err)
	}

	pages, err := client.ListPages()
	if err != nil {
		return fmt.Errorf("list pages: %w", err)
	}
	slog.Debug("fetched pages", slog.String("user", job.user), slog.Int("count", len(pages)))

	// TODO: 2. upsertPages(db, user, pageNames) で SQLite に登録

	// TODO: 3. pending/failed ページを取得

	// TODO: 4. ticker でレート制限しながら各ページを処理
	//          a. PukiWiki からコンテンツ取得
	//          b. Notion API でページ作成（notionPageId 配下に作成）
	//          c. updatePageStatus(db, user, pageName, ...) で結果を記録

	return nil
}

func (pm *pageMigrator) handleErrors() {
	for err := range pm.errCh {
		slog.Error("migration failed", slog.String("error", err.Error()))
	}
}

func (pm *pageMigrator) enqueue(user, notionPageId string) {
	pm.queue <- migrationJob{user: user, notionPageId: notionPageId}
}

func (pm *pageMigrator) isMigrating(user string) bool {
	p := pm.currentUser.Load()
	return p != nil && *p == user
}
