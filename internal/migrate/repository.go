package migrate

import (
	"context"
	"database/sql"
	"pukiwiki-migration/internal/infra/db"
	"time"

	"github.com/google/uuid"
)

type sqliteRepository struct {
	db *sql.DB
	q  *db.Queries
}

func newSQLiteRepository(sqlDB *sql.DB) *sqliteRepository {
	return &sqliteRepository{db: sqlDB, q: db.New(sqlDB)}
}

type migration struct {
	id           string
	pukiwikiUser string
	notionPageID string
	status       string
}

func (s *sqliteRepository) createMigration(ctx context.Context, m migration) error {
	now := time.Now().UTC().Format(time.RFC3339)
	return s.q.CreateMigration(ctx, db.CreateMigrationParams{
		ID:           m.id,
		PukiwikiUser: m.pukiwikiUser,
		NotionPageID: m.notionPageID,
		Status:       m.status,
		CreatedAt:    now,
		UpdatedAt:    now,
	})
}

func (s *sqliteRepository) getMigrationByID(ctx context.Context, id string) (migration, error) {
	row, err := s.q.GetMigrationByID(ctx, id)
	if err != nil {
		return migration{}, err
	}
	return migration{
		id:           row.ID,
		pukiwikiUser: row.PukiwikiUser,
		notionPageID: row.NotionPageID,
		status:       row.Status,
	}, nil
}

func (s *sqliteRepository) updateMigrationStatus(ctx context.Context, id, status string) error {
	return s.q.UpdateMigrationStatus(ctx, db.UpdateMigrationStatusParams{
		ID:        id,
		Status:    status,
		UpdatedAt: time.Now().UTC().Format(time.RFC3339),
	})
}

func (s *sqliteRepository) upsertPages(ctx context.Context, user string, pages []string) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	now := time.Now().UTC().Format(time.RFC3339)
	txQ := s.q.WithTx(tx)
	for _, name := range pages {
		if err := txQ.UpsertPage(ctx, db.UpsertPageParams{
			ID:           uuid.New().String(),
			PukiwikiUser: user,
			PukiwikiPage: name,
			UpdatedAt:    now,
		}); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (s *sqliteRepository) getPendingPages(ctx context.Context, user string) ([]string, error) {
	return s.q.GetPendingPages(ctx, user)
}

func (s *sqliteRepository) updatePageStatus(ctx context.Context, user, page, status, notionID, errMsg string) error {
	return s.q.UpdatePageStatus(ctx, db.UpdatePageStatusParams{
		PukiwikiUser: user,
		PukiwikiPage: page,
		Status:       status,
		NotionID:     notionID,
		ErrorMsg:     errMsg,
		UpdatedAt:    time.Now().UTC().Format(time.RFC3339),
	})
}

type pageSummary struct {
	Total   int
	Done    int
	Failed  int
	Pending int
}

func (s *sqliteRepository) getPageSummary(ctx context.Context, user string) (pageSummary, error) {
	rows, err := s.q.GetStatusSummary(ctx, user)
	if err != nil {
		return pageSummary{}, err
	}
	var summary pageSummary
	for _, row := range rows {
		summary.Total += int(row.Count)
		switch row.Status {
		case "done":
			summary.Done = int(row.Count)
		case "failed":
			summary.Failed = int(row.Count)
		case "pending":
			summary.Pending = int(row.Count)
		}
	}
	return summary, nil
}
