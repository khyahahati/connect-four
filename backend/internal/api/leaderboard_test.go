package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/example/connect-four/backend/internal/store"
)

type stubRepo struct {
	entries []store.LeaderboardEntry
	err     error
	limit   int
}

func (s *stubRepo) GetLeaderboard(limit int) ([]store.LeaderboardEntry, error) {
	s.limit = limit
	return s.entries, s.err
}

func TestGetLeaderboardSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := &stubRepo{entries: []store.LeaderboardEntry{{Username: "alice", Wins: 5}}}
	handler := New(repo)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/leaderboard?limit=3", nil)

	handler.GetLeaderboard(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
	if repo.limit != 3 {
		t.Fatalf("expected repository to receive limit 3, got %d", repo.limit)
	}
	if body := w.Body.String(); body == "[]" {
		t.Fatalf("expected body to contain leaderboard entries")
	}
}

func TestGetLeaderboardBadLimit(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := New(&stubRepo{})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/leaderboard?limit=abc", nil)

	handler.GetLeaderboard(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}
