package matchmaking

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/example/connect-four/backend/internal/game"
	"github.com/example/connect-four/backend/internal/types"
)

const (
	matchInterval        = time.Second
	botFallbackThreshold = 10 * time.Second
)

type socketSender interface {
	SendToUsername(ctx context.Context, username string, message types.ServerMessage) error
}

// Matchmaker coordinates players waiting for a game session.
type Matchmaker struct {
	mu      sync.Mutex
	waiting []waitingPlayer
	gameMgr *game.GameManager
	wsMgr   socketSender
	botName string
}

type waitingPlayer struct {
	username   string
	enqueuedAt time.Time
}

// NewMatchmaker builds a Matchmaker.
func NewMatchmaker(gameMgr *game.GameManager, wsMgr socketSender, botName string) *Matchmaker {
	return &Matchmaker{
		gameMgr: gameMgr,
		wsMgr:   wsMgr,
		botName: botName,
		waiting: make([]waitingPlayer, 0),
	}
}

// Enqueue adds a player to the waiting list.
func (m *Matchmaker) Enqueue(username string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, player := range m.waiting {
		if player.username == username {
			return
		}
	}

	m.waiting = append(m.waiting, waitingPlayer{username: username, enqueuedAt: time.Now().UTC()})
	log.Printf("matchmaker: queued username=%s", username)
}

// Start launches the matchmaking loop in the provided context.
func (m *Matchmaker) Start(ctx context.Context) {
	ticker := time.NewTicker(matchInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			m.tick(ctx)
		}
	}
}

func (m *Matchmaker) tick(ctx context.Context) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for len(m.waiting) >= 2 {
		p1 := m.waiting[0]
		p2 := m.waiting[1]
		m.waiting = m.waiting[2:]

		game := m.gameMgr.CreateGame(p1.username, p2.username)
		log.Printf("matchmaker: created game id=%s p1=%s p2=%s", game.ID, game.Player1, game.Player2)

		m.notifyPlayers(ctx, game)
	}

	if len(m.waiting) == 1 {
		player := m.waiting[0]
		if time.Since(player.enqueuedAt) >= botFallbackThreshold {
			m.waiting = m.waiting[1:]

			game := m.gameMgr.CreateGame(player.username, m.botName)
			log.Printf("matchmaker: created bot game id=%s player=%s bot=%s", game.ID, player.username, m.botName)

			m.notifyBotGame(ctx, game)
		}
	}
}

// WaitingCount returns the number of players currently queued.
func (m *Matchmaker) WaitingCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()

	return len(m.waiting)
}

func (m *Matchmaker) notifyPlayers(ctx context.Context, game *game.Game) {
	msgP1 := types.ServerMessage{
		Type:     "GAME_START",
		GameID:   game.ID,
		You:      1,
		Opponent: game.Player2,
	}

	msgP2 := types.ServerMessage{
		Type:     "GAME_START",
		GameID:   game.ID,
		You:      2,
		Opponent: game.Player1,
	}

	for _, entry := range []struct {
		username string
		msg      types.ServerMessage
	}{
		{game.Player1, msgP1},
		{game.Player2, msgP2},
	} {
		if err := m.wsMgr.SendToUsername(ctx, entry.username, entry.msg); err != nil {
			log.Printf("matchmaker: send GAME_START failed username=%s err=%v", entry.username, err)
		}
	}
}

func (m *Matchmaker) notifyBotGame(ctx context.Context, game *game.Game) {
	msg := types.ServerMessage{
		Type:     "GAME_START",
		GameID:   game.ID,
		You:      1,
		Opponent: game.Player2,
	}

	if err := m.wsMgr.SendToUsername(ctx, game.Player1, msg); err != nil {
		log.Printf("matchmaker: send bot GAME_START failed username=%s err=%v", game.Player1, err)
	}
}
