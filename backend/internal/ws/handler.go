package ws

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"github.com/example/connect-four/backend/internal/bot"
	"github.com/example/connect-four/backend/internal/game"
	"github.com/example/connect-four/backend/internal/matchmaking"
	"github.com/example/connect-four/backend/internal/types"
)

var upgrader = websocket.Upgrader{
	CheckOrigin:     func(r *http.Request) bool { return true },
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Handler provides websocket HTTP handlers.
type Handler struct {
	Manager    *Manager
	GameMgr    *game.GameManager
	Matchmaker *matchmaking.Matchmaker
	Bot        *bot.Bot
}

// NewHandler constructs a Handler.
func NewHandler(manager *Manager, gameMgr *game.GameManager, matchmaker *matchmaking.Matchmaker, botEngine *bot.Bot) *Handler {
	return &Handler{Manager: manager, GameMgr: gameMgr, Matchmaker: matchmaker, Bot: botEngine}
}

// RegisterRoutes wires the websocket endpoint.
func (h *Handler) RegisterRoutes(router *gin.Engine) {
	router.GET("/ws", h.accept)
}

func (h *Handler) accept(c *gin.Context) {
	username := c.Query("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username query parameter is required"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("ws: upgrade failed: %v", err)
		return
	}

	client := h.Manager.Register(username, conn)
	ctx, cancel := context.WithCancel(context.Background())

	if err := h.Manager.Send(ctx, client, types.ServerMessage{Type: "INFO", Message: "Welcome to Connect Four."}); err != nil {
		log.Printf("ws: send welcome failed: %v", err)
	}

	if gameID := c.Query("gameId"); gameID != "" {
		if err := h.Manager.Send(ctx, client, types.ServerMessage{Type: "INFO", Message: "Reconnect ack for game " + gameID}); err != nil {
			log.Printf("ws: send reconnect ack failed: %v", err)
		}
	}

	if h.Matchmaker != nil {
		h.Matchmaker.Enqueue(username)
	}

	go h.listen(ctx, cancel, client)
}

func (h *Handler) listen(ctx context.Context, cancel context.CancelFunc, conn *Connection) {
	defer func() {
		cancel()
		h.Manager.Unregister(conn)
		_ = conn.Socket.Close()
	}()

	conn.Socket.SetReadDeadline(time.Now().Add(60 * time.Second))
	conn.Socket.SetPongHandler(func(string) error {
		conn.Socket.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		typeCode, payload, err := conn.Socket.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("ws: read error id=%s err=%v", conn.ID, err)
			}
			return
		}

		if typeCode != websocket.TextMessage {
			continue
		}

		var msg types.ClientMessage
		if err := json.Unmarshal(payload, &msg); err != nil {
			h.sendInfo(ctx, conn, "Invalid message format")
			continue
		}

		if err := h.handleMessage(ctx, conn, msg); err != nil {
			h.sendInfo(ctx, conn, err.Error())
		}
	}
}

func (h *Handler) handleMessage(ctx context.Context, conn *Connection, msg types.ClientMessage) error {
	switch msg.Type {
	case "MAKE_MOVE":
		return h.handleMakeMove(ctx, conn, msg)
	case "RECONNECT":
		return h.handleReconnect(ctx, conn, msg)
	default:
		return errors.New("unsupported message type")
	}
}

func (h *Handler) handleMakeMove(ctx context.Context, conn *Connection, msg types.ClientMessage) error {
	if msg.Col == nil {
		return errors.New("MAKE_MOVE missing col")
	}
	if msg.GameID == "" {
		return errors.New("MAKE_MOVE missing gameId")
	}

	log.Printf("ws: MAKE_MOVE id=%s username=%s gameId=%s col=%d", conn.ID, conn.Username, msg.GameID, *msg.Col)

	if h.GameMgr == nil {
		return errors.New("game manager unavailable")
	}

	if _, ok := h.GameMgr.GetGame(msg.GameID); !ok {
		return errors.New("unknown game")
	}

	updatedGame, result, err := h.GameMgr.ApplyMove(msg.GameID, conn.Username, *msg.Col)
	if err != nil {
		return h.sendInfo(ctx, conn, err.Error())
	}

	h.sendBoardUpdate(ctx, updatedGame)
	h.handleGameOutcome(ctx, updatedGame, conn.Username, result)

	if result == game.CONTINUE && updatedGame.Player2 == "BOT" && h.Bot != nil && updatedGame.CurrentTurn == 2 {
		h.handleBotTurn(ctx, updatedGame)
	}

	return nil
}

func (h *Handler) handleReconnect(ctx context.Context, conn *Connection, msg types.ClientMessage) error {
	if msg.Username == "" {
		return errors.New("RECONNECT missing username")
	}

	log.Printf("ws: RECONNECT request username=%s gameId=%s", msg.Username, msg.GameID)
	return h.sendInfo(ctx, conn, "Reconnect acknowledged")
}

func (h *Handler) sendInfo(ctx context.Context, conn *Connection, message string) error {
	sendCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return h.Manager.Send(sendCtx, conn, types.ServerMessage{Type: "INFO", Message: message})
}

func (h *Handler) handleGameOutcome(ctx context.Context, gameState *game.Game, mover string, result game.MoveResult) {
	sendCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	switch result {
	case game.WIN:
		h.sendGameOver(sendCtx, gameState, mover, false)
	case game.DRAW:
		h.sendGameOver(sendCtx, gameState, "", true)
	}
}

func (h *Handler) sendBoardUpdate(ctx context.Context, gameState *game.Game) {
	sendCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	msg := types.ServerMessage{
		Type:        "BOARD_UPDATE",
		GameID:      gameState.ID,
		Board:       gameState.Board,
		CurrentTurn: gameState.CurrentTurn,
	}

	h.sendToPlayers(sendCtx, gameState, msg, &msg)
}

func (h *Handler) sendGameOver(ctx context.Context, gameState *game.Game, winner string, draw bool) {
	if draw {
		msg := types.ServerMessage{Type: "GAME_OVER", GameID: gameState.ID, Board: gameState.Board, Result: "DRAW"}
		h.sendToPlayers(ctx, gameState, msg, &msg)
		return
	}

	msgP1 := types.ServerMessage{Type: "GAME_OVER", GameID: gameState.ID, Board: gameState.Board, Result: "LOSS"}
	msgP2 := types.ServerMessage{Type: "GAME_OVER", GameID: gameState.ID, Board: gameState.Board, Result: "LOSS"}

	if winner == gameState.Player1 {
		msgP1.Result = "WIN"
	} else if winner == gameState.Player2 {
		msgP2.Result = "WIN"
	}

	if gameState.Player2 == "BOT" {
		// Only notify the human player.
		h.sendToPlayers(ctx, gameState, msgP1, nil)
		return
	}

	h.sendToPlayers(ctx, gameState, msgP1, &msgP2)
}

func (h *Handler) sendToPlayers(ctx context.Context, gameState *game.Game, msgP1 types.ServerMessage, msgP2 *types.ServerMessage) {
	if err := h.Manager.SendToUsername(ctx, gameState.Player1, msgP1); err != nil {
		log.Printf("ws: failed to send to %s: %v", gameState.Player1, err)
	}

	if gameState.Player2 == "BOT" || msgP2 == nil {
		return
	}

	if err := h.Manager.SendToUsername(ctx, gameState.Player2, *msgP2); err != nil {
		log.Printf("ws: failed to send to %s: %v", gameState.Player2, err)
	}
}

func (h *Handler) handleBotTurn(ctx context.Context, currentGame *game.Game) {
	if h.Bot == nil {
		return
	}

	botGame, result, _, err := h.Bot.TakeTurn(currentGame.ID)
	if err != nil {
		log.Printf("ws: bot move error gameID=%s err=%v", currentGame.ID, err)
		return
	}

	sendCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	h.sendBoardUpdate(sendCtx, botGame)

	switch result {
	case game.WIN:
		h.sendGameOver(sendCtx, botGame, botGame.Player2, false)
	case game.DRAW:
		h.sendGameOver(sendCtx, botGame, "", true)
	}
}
