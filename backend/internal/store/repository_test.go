package store

import (
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestEnsurePlayer(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock new: %v", err)
	}
	defer db.Close()

	repo := NewRepository(db)

	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO players (username) VALUES ($1) ON CONFLICT (username) DO NOTHING")).
		WithArgs("alice").
		WillReturnResult(sqlmock.NewResult(0, 1))

	if err := repo.EnsurePlayer("alice"); err != nil {
		t.Fatalf("EnsurePlayer failed: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestIncrementWinEnsuresPlayerOnMissingRow(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock new: %v", err)
	}
	defer db.Close()

	repo := NewRepository(db)

	mock.ExpectExec(regexp.QuoteMeta("UPDATE players SET wins = wins + 1 WHERE username = $1")).
		WithArgs("bob").
		WillReturnResult(sqlmock.NewResult(0, 0))

	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO players (username) VALUES ($1) ON CONFLICT (username) DO NOTHING")).
		WithArgs("bob").
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectExec(regexp.QuoteMeta("UPDATE players SET wins = wins + 1 WHERE username = $1")).
		WithArgs("bob").
		WillReturnResult(sqlmock.NewResult(0, 1))

	if err := repo.IncrementWin("bob"); err != nil {
		t.Fatalf("IncrementWin failed: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestSaveCompletedGame(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock new: %v", err)
	}
	defer db.Close()

	repo := NewRepository(db)

	finished := CompletedGame{
		ID:      "game-123",
		Player1: "alice",
		Player2: "bob",
		Winner:  strPtr("alice"),
		IsDraw:  false,
		Moves: []CompletedMove{
			{Player: "alice", Column: 0, MoveNumber: 1},
			{Player: "bob", Column: 1, MoveNumber: 2},
		},
		StartedAt: time.Now().UTC().Add(-time.Hour),
		EndedAt:   time.Now().UTC(),
	}

	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO games (id, player1, player2, winner, is_draw, moves, started_at, ended_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)")).
		WithArgs(
			finished.ID,
			finished.Player1,
			finished.Player2,
			sqlmock.AnyArg(),
			finished.IsDraw,
			sqlmock.AnyArg(),
			finished.StartedAt,
			finished.EndedAt,
		).
		WillReturnResult(sqlmock.NewResult(0, 1))

	if err := repo.SaveCompletedGame(&finished); err != nil {
		t.Fatalf("SaveCompletedGame failed: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestGetLeaderboard(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock new: %v", err)
	}
	defer db.Close()

	repo := NewRepository(db)

	rows := sqlmock.NewRows([]string{"username", "wins"}).
		AddRow("alice", 5).
		AddRow("bob", 3)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT username, wins FROM players ORDER BY wins DESC, username ASC LIMIT $1")).
		WithArgs(5).
		WillReturnRows(rows)

	result, err := repo.GetLeaderboard(5)
	if err != nil {
		t.Fatalf("GetLeaderboard failed: %v", err)
	}

	if len(result) != 2 || result[0].Username != "alice" || result[0].Wins != 5 {
		t.Fatalf("unexpected leaderboard result: %+v", result)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func strPtr(v string) *string {
	return &v
}
