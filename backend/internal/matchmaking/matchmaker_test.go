package matchmaking

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/example/connect-four/backend/internal/game"
	"github.com/example/connect-four/backend/internal/types"
)

type stubSocketManager struct {
	mu          sync.Mutex
	connections map[string]struct{}
	messages    map[string][]types.ServerMessage
}

func newStubSocketManager() *stubSocketManager {
	return &stubSocketManager{
		connections: make(map[string]struct{}),
		messages:    make(map[string][]types.ServerMessage),
	}
}

func (s *stubSocketManager) add(username string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.connections[username] = struct{}{}
}

func (s *stubSocketManager) messagesFor(username string) []types.ServerMessage {
	s.mu.Lock()
	defer s.mu.Unlock()

	return append([]types.ServerMessage(nil), s.messages[username]...)
}

func (s *stubSocketManager) SendToUsername(ctx context.Context, username string, message types.ServerMessage) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.connections[username]; !ok {
		return nil
	}

	s.messages[username] = append(s.messages[username], message)
	return nil
}

func TestMatchmakerPairsPlayers(t *testing.T) {
	ctx := context.Background()
	gm := game.NewManager()
	sockets := newStubSocketManager()
	sockets.add("alice")
	sockets.add("bob")

	matcher := NewMatchmaker(gm, sockets, "BOT")
	matcher.Enqueue("alice")
	matcher.Enqueue("bob")

	matcher.tick(ctx)

	if matcher.WaitingCount() != 0 {
		t.Fatalf("expected queue to be empty, got %d", matcher.WaitingCount())
	}

	game, ok := gm.FindGameByPlayers("alice", "bob")
	if !ok {
		t.Fatalf("expected game between alice and bob to exist")
	}

	msgsAlice := sockets.messagesFor("alice")
	msgsBob := sockets.messagesFor("bob")

	if len(msgsAlice) != 1 || len(msgsBob) != 1 {
		t.Fatalf("expected exactly one GAME_START per player, got %d and %d", len(msgsAlice), len(msgsBob))
	}

	msgAlice := msgsAlice[0]
	msgBob := msgsBob[0]

	if msgAlice.Type != "GAME_START" || msgBob.Type != "GAME_START" {
		t.Fatalf("expected GAME_START messages, got %q and %q", msgAlice.Type, msgBob.Type)
	}

	if msgAlice.GameID == "" || msgBob.GameID == "" {
		t.Fatalf("expected non-empty game ids")
	}

	if msgAlice.GameID != msgBob.GameID || msgAlice.GameID != game.ID {
		t.Fatalf("expected shared game id %s, got %s and %s", game.ID, msgAlice.GameID, msgBob.GameID)
	}

	if msgAlice.You != 1 || msgBob.You != 2 {
		t.Fatalf("expected turn assignments 1/2, got %d/%d", msgAlice.You, msgBob.You)
	}

	if msgAlice.Opponent != "bob" || msgBob.Opponent != "alice" {
		t.Fatalf("unexpected opponents %q/%q", msgAlice.Opponent, msgBob.Opponent)
	}
}

func TestMatchmakerBotFallback(t *testing.T) {
	ctx := context.Background()
	gm := game.NewManager()
	sockets := newStubSocketManager()
	sockets.add("carol")

	matcher := NewMatchmaker(gm, sockets, "BOT")
	matcher.Enqueue("carol")

	matcher.waiting[0].enqueuedAt = time.Now().Add(-botFallbackThreshold - time.Second)

	matcher.tick(ctx)

	if matcher.WaitingCount() != 0 {
		t.Fatalf("expected queue to be empty after bot fallback, got %d", matcher.WaitingCount())
	}

	_, ok := gm.FindGameByPlayers("carol", "BOT")
	if !ok {
		t.Fatalf("expected game between carol and BOT")
	}

	msgs := sockets.messagesFor("carol")
	if len(msgs) != 1 {
		t.Fatalf("expected one GAME_START message, got %d", len(msgs))
	}

	msg := msgs[0]
	if msg.Type != "GAME_START" || msg.Opponent != "BOT" {
		t.Fatalf("unexpected bot fallback message: %+v", msg)
	}
	if msg.You != 1 {
		t.Fatalf("expected player to be assigned turn 1, got %d", msg.You)
	}
}

func TestMatchmakerRetainsUnpairedPlayer(t *testing.T) {
	ctx := context.Background()
	gm := game.NewManager()
	sockets := newStubSocketManager()
	sockets.add("alice")
	sockets.add("bob")
	sockets.add("charlie")

	matcher := NewMatchmaker(gm, sockets, "BOT")
	matcher.Enqueue("alice")
	matcher.Enqueue("bob")
	matcher.Enqueue("charlie")

	matcher.tick(ctx)

	if matcher.WaitingCount() != 1 {
		t.Fatalf("expected exactly one player waiting, got %d", matcher.WaitingCount())
	}

	if _, ok := gm.FindGameByPlayers("alice", "bob"); !ok {
		t.Fatalf("expected game between alice and bob")
	}

	if len(sockets.messagesFor("charlie")) != 0 {
		t.Fatalf("expected no messages for charlie")
	}
}
