package migrate

import (
	"context"
	"database/sql"
	"log/slog"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/moriT958/libpukiwiki"
)

const pukiScopePrefix = "seminar-personal"

type migrationUsecase struct {
	repo         *sqliteRepository
	pukiBaseURL  string
	pukiUsername string
	pukiPassword string
	notionToken  string
}

func NewMigrationUsecase(db *sql.DB, pukiBaseURL, pukiUsername, pukiPassword, notionToken string) (*migrationUsecase, error) {
	// validate puki config by attempting client creation
	if _, err := libpukiwiki.NewClient(pukiBaseURL, libpukiwiki.WithAuth(pukiUsername, pukiPassword)); err != nil {
		return nil, err
	}
	return &migrationUsecase{
		repo:         newSQLiteRepository(db),
		pukiBaseURL:  pukiBaseURL,
		pukiUsername: pukiUsername,
		pukiPassword: pukiPassword,
		notionToken:  notionToken,
	}, nil
}

func (u *migrationUsecase) start(ctx context.Context, user, notionPageID string) (migration, error) {
	m := migration{
		id:           uuid.New().String(),
		pukiwikiUser: user,
		notionPageID: notionPageID,
		status:       "pending",
	}
	if err := u.repo.createMigration(ctx, m); err != nil {
		return migration{}, err
	}
	go u.startJob(m.id, user, notionPageID)
	return m, nil
}

func (u *migrationUsecase) getStatus(ctx context.Context, id string) (migration, pageSummary, error) {
	m, err := u.repo.getMigrationByID(ctx, id)
	if err != nil {
		return migration{}, pageSummary{}, err
	}
	summary, err := u.repo.getPageSummary(ctx, m.pukiwikiUser)
	if err != nil {
		return migration{}, pageSummary{}, err
	}
	return m, summary, nil
}

func (u *migrationUsecase) startJob(migID, user, notionPageID string) {
	ctx := context.Background()
	slog.Info("migration started", "migrationID", migID, "user", user)

	if err := u.repo.updateMigrationStatus(ctx, migID, "running"); err != nil {
		slog.Error("failed to update migration status", "error", err)
		return
	}

	puki, err := libpukiwiki.NewClient(u.pukiBaseURL, libpukiwiki.WithAuth(u.pukiUsername, u.pukiPassword), libpukiwiki.WithScope(pukiScopePrefix+"/"+user))
	if err != nil {
		slog.Error("failed to create pukiwiki client", "error", err)
		u.repo.updateMigrationStatus(ctx, migID, "failed")
		return
	}

	if err := puki.Login(); err != nil {
		slog.Error("failed to login to pukiwiki", "error", err)
		u.repo.updateMigrationStatus(ctx, migID, "failed")
		return
	}

	pages, err := puki.ListPages()
	if err != nil {
		slog.Error("failed to list pukiwiki pages", "error", err)
		u.repo.updateMigrationStatus(ctx, migID, "failed")
		return
	}
	slog.Info("pages fetched", "count", len(pages))

	if err := u.repo.upsertPages(ctx, user, pages); err != nil {
		slog.Error("failed to upsert pages", "error", err)
		u.repo.updateMigrationStatus(ctx, migID, "failed")
		return
	}

	pending, err := u.repo.getPendingPages(ctx, user)
	if err != nil {
		slog.Error("failed to get pending pages", "error", err)
		u.repo.updateMigrationStatus(ctx, migID, "failed")
		return
	}
	slog.Info("migration processing", "total", len(pending))

	jobs := make(chan string, len(pending))
	for _, p := range pending {
		jobs <- p
	}
	close(jobs)

	// Notion API rate limit: 3 req/sec
	ticker := time.NewTicker(time.Second / 3)
	defer ticker.Stop()

	var wg sync.WaitGroup
	for range 4 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for pageName := range jobs {
				slog.Info("migrating page", "page", pageName)

				content, err := puki.GetPageSource(pageName)
				if err != nil {
					slog.Warn("failed to get page content", "page", pageName, "error", err)
					if dbErr := u.repo.updatePageStatus(ctx, user, pageName, "failed", "", err.Error()); dbErr != nil {
						slog.Error("failed to update page status to failed", "page", pageName, "error", dbErr)
					}
					continue
				}

				<-ticker.C

				notionID, err := createPage(u.notionToken, notionPageID, pageName, content)
				if err != nil {
					slog.Warn("failed to create notion page", "page", pageName, "error", err)
					if dbErr := u.repo.updatePageStatus(ctx, user, pageName, "failed", "", err.Error()); dbErr != nil {
						slog.Error("failed to update page status to failed", "page", pageName, "error", dbErr)
					}
					continue
				}

				if dbErr := u.repo.updatePageStatus(ctx, user, pageName, "done", notionID, ""); dbErr != nil {
					slog.Error("failed to update page status to done", "page", pageName, "error", dbErr)
					continue
				}
				slog.Info("page migrated", "page", pageName)
			}
		}()
	}

	wg.Wait()
	u.repo.updateMigrationStatus(ctx, migID, "done")
	slog.Info("migration completed", "migrationID", migID, "user", user)
}
