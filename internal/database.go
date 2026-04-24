package internal

import (
	"database/sql"
	"time"
)

type statusSummary struct {
	Total   int `json:"total"`
	Done    int `json:"done"`
	Failed  int `json:"failed"`
	Pending int `json:"pending"`
}

func upsertPages(db *sql.DB, user string, pageNames []string) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare("INSERT OR IGNORE INTO pages (user, page_name, status, updated_at) VALUES (?, ?, 'pending', ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	now := time.Now().UTC().Format(time.RFC3339)
	for _, name := range pageNames {
		if _, err := stmt.Exec(user, name, now); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func updatePageStatus(db *sql.DB, user, pageName, status, notionID, errMsg string) error {
	now := time.Now().UTC().Format(time.RFC3339)
	_, err := db.Exec(
		"UPDATE pages SET status=?, notion_id=?, error_msg=?, updated_at=? WHERE user=? AND page_name=?",
		status, notionID, errMsg, now, user, pageName,
	)
	return err
}

func getStatusSummary(db *sql.DB, user string) (statusSummary, error) {
	rows, err := db.Query("SELECT status, COUNT(*) FROM pages WHERE user=? GROUP BY status", user)
	if err != nil {
		return statusSummary{}, err
	}
	defer rows.Close()

	var s statusSummary
	for rows.Next() {
		var status string
		var count int
		if err := rows.Scan(&status, &count); err != nil {
			return statusSummary{}, err
		}
		s.Total += count
		switch status {
		case "done":
			s.Done = count
		case "failed":
			s.Failed = count
		case "pending":
			s.Pending = count
		}
	}
	return s, nil
}
