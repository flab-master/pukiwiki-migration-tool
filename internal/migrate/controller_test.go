package migrate

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type usecaseStub struct {
	startErr      error
	statusResult  migration
	statusSummary pageSummary
	statusErr     error
}

func (m *usecaseStub) start(_ context.Context, _, _ string) (migration, error) {
	return migration{}, m.startErr
}

func (m *usecaseStub) getStatus(_ context.Context, _ string) (migration, pageSummary, error) {
	return m.statusResult, m.statusSummary, m.statusErr
}

func TestHandleMigrate(t *testing.T) {
	t.Run("正常なリクエストは 202 を返す", func(t *testing.T) {
		handler := NewMigrationController(&usecaseStub{})
		body := `{"user":"testuser","notionPageId":"page-1"}`
		req := httptest.NewRequest(http.MethodPost, "/migrate", strings.NewReader(body))
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		if w.Code != http.StatusAccepted {
			t.Errorf("expected 202, got %d", w.Code)
		}
	})

	t.Run("リクエストの user フィールドが空の場合は 400 を返す", func(t *testing.T) {
		handler := NewMigrationController(&usecaseStub{})
		req := httptest.NewRequest(http.MethodPost, "/migrate", strings.NewReader(`{}`))
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", w.Code)
		}
	})

	t.Run("Usecase クラスがエラーを返した場合は 500 を返す", func(t *testing.T) {
		handler := NewMigrationController(&usecaseStub{startErr: errors.New("db error")})
		body := `{"user":"testuser","notionPageId":"page-1"}`
		req := httptest.NewRequest(http.MethodPost, "/migrate", strings.NewReader(body))
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected 500, got %d", w.Code)
		}
	})
}

func TestHandleStatus(t *testing.T) {
	t.Run("パスパラメータで 移行 ID を受け取り、移行ステータスの JSON を返す", func(t *testing.T) {
		mock := &usecaseStub{
			statusResult:  migration{id: "mig-1", pukiwikiUser: "testuser", status: "done"},
			statusSummary: pageSummary{Total: 3, Done: 3},
		}
		handler := NewMigrationController(mock)
		req := httptest.NewRequest(http.MethodGet, "/migrate/mig-1/status", nil)
		req.SetPathValue("id", "mig-1")
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", w.Code)
		}
		var resp statusResponse
		if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if resp.ID != "mig-1" || resp.User != "testuser" || resp.Status != "done" {
			t.Errorf("unexpected response: %+v", resp)
		}
		if resp.Summary.Total != 3 || resp.Summary.Done != 3 {
			t.Errorf("unexpected summary: %+v", resp.Summary)
		}
	})

	t.Run("Usecase クラスがエラーを返した場合は 500 を返す", func(t *testing.T) {
		handler := NewMigrationController(&usecaseStub{statusErr: errors.New("not found")})
		req := httptest.NewRequest(http.MethodGet, "/migrate/bad-id/status", nil)
		req.SetPathValue("id", "bad-id")
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected 500, got %d", w.Code)
		}
	})
}
