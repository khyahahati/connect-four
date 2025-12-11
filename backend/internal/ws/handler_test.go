package ws

import (
	"context"
	"encoding/json"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"github.com/example/connect-four/backend/internal/bot"
	"github.com/example/connect-four/backend/internal/game"
	"github.com/example/connect-four/backend/internal/store"
	"github.com/example/connect-four/backend/internal/types"
)

type mockResultStore struct {
	mu         sync.Mutex
	ensures    []string
	increments []string
	saved      []*store.CompletedGame
	wg         sync.WaitGroup
}

func newMockResultStore() *mockResultStore {
	m := &mockResultStore{}
	m.wg.Add(1)
	return m
}

func (m *mockResultStore) EnsurePlayer(username string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.ensures = append(m.ensures, username)
	return nil
}

func (m *mockResultStore) IncrementWin(username string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.increments = append(m.increments, username)
	return nil
}

func (m *mockResultStore) SaveCompletedGame(game *store.CompletedGame) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.saved = append(m.saved, game)
	m.wg.Done()
	return nil
}

func (m *mockResultStore) waitForSave(timeout time.Duration) bool {
	done := make(chan struct{})
	go func() {
		m.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		return true
	case <-time.After(timeout):
		return false
	}
}
func TestWebSocketMakeMoveBoardUpdate(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	manager := NewManager()
	gameManager := game.NewManager()
	handler := NewHandler(manager, gameManager, nil, nil, nil)
	handler.RegisterRoutes(r)

	created := gameManager.CreateGame("tester", "opponent")

	ts := httptest.NewServer(r)
	t.Cleanup(ts.Close)
	t.Cleanup(func() { manager.Shutdown(context.Background()) })

	url := "ws" + strings.TrimPrefix(ts.URL, "http") + "/ws?username=tester"
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		t.Fatalf("dial failed: %v", err)
	}
	t.Cleanup(func() { _ = conn.Close() })

	readWithDeadline := func() types.ServerMessage {
		if err := conn.SetReadDeadline(time.Now().Add(5 * time.Second)); err != nil {
			t.Fatalf("set read deadline: %v", err)
		}
		_, data, err := conn.ReadMessage()
		if err != nil {
			t.Fatalf("read message: %v", err)
		}
		var msg types.ServerMessage
		if err := json.Unmarshal(data, &msg); err != nil {
			t.Fatalf("decode message: %v", err)
		}
		return msg
	}

	welcome := readWithDeadline()
	if welcome.Type != "INFO" {
		t.Fatalf("expected INFO welcome, got %s", welcome.Type)
	}

	payload := map[string]any{"type": "MAKE_MOVE", "gameId": created.ID, "col": 3}
	if err := conn.WriteJSON(payload); err != nil {
		t.Fatalf("write MAKE_MOVE: %v", err)
	}

	update := readWithDeadline()
	if update.Type != "BOARD_UPDATE" {
		t.Fatalf("expected BOARD_UPDATE, got %s", update.Type)
	}
	if update.GameID != created.ID {
		t.Fatalf("unexpected game id %s", update.GameID)
	}
	if update.CurrentTurn != 2 {
		t.Fatalf("expected current turn 2, got %d", update.CurrentTurn)
	}
	if len(update.Board) != game.Rows || len(update.Board[0]) != game.Columns {
		t.Fatalf("unexpected board dimensions")
	}
	if update.Board[game.Rows-1][3] != 1 {
		t.Fatalf("expected disc at bottom row col 3")
	}
}

func TestWebSocketMakeMoveGameOver(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	manager := NewManager()
	gameManager := game.NewManager()
	handler := NewHandler(manager, gameManager, nil, nil, nil)
	handler.RegisterRoutes(r)

	created := gameManager.CreateGame("tester", "opponent")
	if _, _, err := gameManager.ApplyMove(created.ID, "tester", 0); err != nil {
		t.Fatalf("setup move: %v", err)
	}
	if _, _, err := gameManager.ApplyMove(created.ID, "opponent", 1); err != nil {
		t.Fatalf("setup move: %v", err)
	}
	if _, _, err := gameManager.ApplyMove(created.ID, "tester", 0); err != nil {
		t.Fatalf("setup move: %v", err)
	}
	if _, _, err := gameManager.ApplyMove(created.ID, "opponent", 1); err != nil {
		t.Fatalf("setup move: %v", err)
	}
	if _, _, err := gameManager.ApplyMove(created.ID, "tester", 0); err != nil {
		t.Fatalf("setup move: %v", err)
	}
	if _, _, err := gameManager.ApplyMove(created.ID, "opponent", 1); err != nil {
		t.Fatalf("setup move: %v", err)
	}

	ts := httptest.NewServer(r)
	t.Cleanup(ts.Close)
	t.Cleanup(func() { manager.Shutdown(context.Background()) })

	url := "ws" + strings.TrimPrefix(ts.URL, "http") + "/ws?username=tester"
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		t.Fatalf("dial failed: %v", err)
	}
	t.Cleanup(func() { _ = conn.Close() })

	readWithDeadline := func() types.ServerMessage {
		if err := conn.SetReadDeadline(time.Now().Add(5 * time.Second)); err != nil {
			t.Fatalf("set read deadline: %v", err)
		}
		_, data, err := conn.ReadMessage()
		if err != nil {
			t.Fatalf("read message: %v", err)
		}
		var msg types.ServerMessage
		if err := json.Unmarshal(data, &msg); err != nil {
			t.Fatalf("decode message: %v", err)
		}
		return msg
	}

	_ = readWithDeadline() // welcome

	payload := map[string]any{"type": "MAKE_MOVE", "gameId": created.ID, "col": 0}
	if err := conn.WriteJSON(payload); err != nil {
		t.Fatalf("write MAKE_MOVE: %v", err)
	}

	update := readWithDeadline()
	if update.Type != "BOARD_UPDATE" {
		t.Fatalf("expected BOARD_UPDATE, got %s", update.Type)
	}
	state, ok := gameManager.GetGame(created.ID)
	if !ok {
		t.Fatalf("game should exist")
	}
	if state.Winner == nil {
		t.Fatalf("expected winner after winning move")
	}

	gameOver := readWithDeadline()
	if gameOver.Type != "GAME_OVER" {
		t.Fatalf("expected GAME_OVER, got %s", gameOver.Type)
	}
	if gameOver.Result != "WIN" {
		t.Fatalf("expected WIN result, got %s", gameOver.Result)
	}
}

func TestWebSocketBotAutoResponds(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	manager := NewManager()
	gameManager := game.NewManager()
	botEngine := bot.New(gameManager)
	handler := NewHandler(manager, gameManager, nil, botEngine, nil)
	handler.RegisterRoutes(r)

	created := gameManager.CreateGame("tester", "BOT")

	ts := httptest.NewServer(r)
	t.Cleanup(ts.Close)
	t.Cleanup(func() { manager.Shutdown(context.Background()) })

	url := "ws" + strings.TrimPrefix(ts.URL, "http") + "/ws?username=tester"
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		t.Fatalf("dial failed: %v", err)
	}
	t.Cleanup(func() { _ = conn.Close() })

	readWithDeadline := func() types.ServerMessage {
		if err := conn.SetReadDeadline(time.Now().Add(5 * time.Second)); err != nil {
			t.Fatalf("set read deadline: %v", err)
		}
		_, data, err := conn.ReadMessage()
		if err != nil {
			t.Fatalf("read message: %v", err)
		}
		var msg types.ServerMessage
		if err := json.Unmarshal(data, &msg); err != nil {
			t.Fatalf("decode message: %v", err)
		}
		return msg
	}

	_ = readWithDeadline() // welcome

	payload := map[string]any{"type": "MAKE_MOVE", "gameId": created.ID, "col": 0}
	if err := conn.WriteJSON(payload); err != nil {
		t.Fatalf("write MAKE_MOVE: %v", err)
	}

	first := readWithDeadline()
	if first.Type != "BOARD_UPDATE" {
		t.Fatalf("expected BOARD_UPDATE, got %s", first.Type)
	}

	second := readWithDeadline()
	if second.Type != "BOARD_UPDATE" {
		t.Fatalf("expected second BOARD_UPDATE, got %s", second.Type)
	}
	if second.CurrentTurn != 1 {
		t.Fatalf("expected turn to return to player 1, got %d", second.CurrentTurn)
	}
	if second.Board[game.Rows-1][3] != 2 {
		t.Fatalf("expected bot to play center column")
	}
}

func TestGameOverPersistsResults(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	manager := NewManager()
	gameManager := game.NewManager()
	mockStore := newMockResultStore()
	handler := NewHandler(manager, gameManager, nil, nil, mockStore)
	handler.RegisterRoutes(r)

	created := gameManager.CreateGame("tester", "opponent")
	if _, _, err := gameManager.ApplyMove(created.ID, "tester", 0); err != nil {
		t.Fatalf("setup move: %v", err)
	}
	if _, _, err := gameManager.ApplyMove(created.ID, "opponent", 1); err != nil {
		t.Fatalf("setup move: %v", err)
	}
	if _, _, err := gameManager.ApplyMove(created.ID, "tester", 0); err != nil {
		t.Fatalf("setup move: %v", err)
	}
	if _, _, err := gameManager.ApplyMove(created.ID, "opponent", 1); err != nil {
		t.Fatalf("setup move: %v", err)
	}
	if _, _, err := gameManager.ApplyMove(created.ID, "tester", 0); err != nil {
		t.Fatalf("setup move: %v", err)
	}
	if _, _, err := gameManager.ApplyMove(created.ID, "opponent", 1); err != nil {
		t.Fatalf("setup move: %v", err)
	}

	ts := httptest.NewServer(r)
	t.Cleanup(ts.Close)
	t.Cleanup(func() { manager.Shutdown(context.Background()) })

	url := "ws" + strings.TrimPrefix(ts.URL, "http") + "/ws?username=tester"
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		t.Fatalf("dial failed: %v", err)
	}
	t.Cleanup(func() { _ = conn.Close() })

	readWithDeadline := func() types.ServerMessage {
		if err := conn.SetReadDeadline(time.Now().Add(5 * time.Second)); err != nil {
			t.Fatalf("set read deadline: %v", err)
		}
		_, data, err := conn.ReadMessage()
		if err != nil {
			t.Fatalf("read message: %v", err)
		}
		var msg types.ServerMessage
		if err := json.Unmarshal(data, &msg); err != nil {
			t.Fatalf("decode message: %v", err)
		}
		return msg
	}

	_ = readWithDeadline() // welcome

	payload := map[string]any{"type": "MAKE_MOVE", "gameId": created.ID, "col": 0}
	if err := conn.WriteJSON(payload); err != nil {
		t.Fatalf("write MAKE_MOVE: %v", err)
	}

	_ = readWithDeadline() // board update
	gameOver := readWithDeadline()
	if gameOver.Type != "GAME_OVER" {
		t.Fatalf("expected GAME_OVER, got %s", gameOver.Type)
	}

	if !mockStore.waitForSave(2 * time.Second) {
		t.Fatalf("persistence did not complete in time")
	}

	mockStore.mu.Lock()
	ensures := append([]string(nil), mockStore.ensures...)
	increments := append([]string(nil), mockStore.increments...)
	var saved *store.CompletedGame
	if len(mockStore.saved) > 0 {
		saved = mockStore.saved[0]
	}
	mockStore.mu.Unlock()

	if saved == nil {
		t.Fatalf("expected saved game record")
	}
	if saved.Winner == nil || *saved.Winner != "tester" {
		t.Fatalf("expected winner tester, got %+v", saved.Winner)
	}
	if len(saved.Moves) == 0 {
		t.Fatalf("expected moves recorded")
	}

	if len(ensures) < 2 {
		t.Fatalf("expected EnsurePlayer called for both players, got %v", ensures)
	}
	hasTester := false
	hasOpponent := false
	for _, u := range ensures {
		if u == "tester" {
			hasTester = true
		}
		if u == "opponent" {
			hasOpponent = true
		}
	}
	if !hasTester || !hasOpponent {
		t.Fatalf("ensure player did not include both users: %v", ensures)
	}

	if len(increments) == 0 || increments[0] != "tester" {
		t.Fatalf("expected increment for tester, got %v", increments)
	}
}
