package ws

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	"github.com/example/connect-four/backend/internal/types"
)

// Connection represents an active websocket client.
type Connection struct {
	ID       string
	Username string
	Socket   *websocket.Conn
	sendMu   sync.Mutex
}

// Manager coordinates active websocket connections.
type Manager struct {
	mu          sync.RWMutex
	connections map[string]*Connection
	byUsername  map[string]*Connection
}

// NewManager builds a Manager instance.
func NewManager() *Manager {
	return &Manager{
		connections: make(map[string]*Connection),
		byUsername:  make(map[string]*Connection),
	}
}

// Register adds a websocket connection to the manager.
func (m *Manager) Register(username string, socket *websocket.Conn) *Connection {
	conn := &Connection{
		ID:       uuid.NewString(),
		Username: username,
		Socket:   socket,
	}

	m.mu.Lock()
	m.connections[conn.ID] = conn
	m.byUsername[username] = conn
	m.mu.Unlock()

	log.Printf("ws: connected id=%s username=%s", conn.ID, username)
	return conn
}

// Unregister removes a connection from the manager.
func (m *Manager) Unregister(conn *Connection) {
	if conn == nil {
		return
	}

	m.mu.Lock()
	delete(m.connections, conn.ID)
	delete(m.byUsername, conn.Username)
	m.mu.Unlock()

	log.Printf("ws: disconnected id=%s username=%s", conn.ID, conn.Username)
}

// Send transmits a server message to a specific connection.
func (m *Manager) Send(ctx context.Context, conn *Connection, message types.ServerMessage) error {
	if conn == nil {
		return nil
	}

	payload, err := json.Marshal(message)
	if err != nil {
		return err
	}

	conn.sendMu.Lock()
	defer conn.sendMu.Unlock()

	if deadline, ok := ctx.Deadline(); ok {
		_ = conn.Socket.SetWriteDeadline(deadline)
	} else {
		_ = conn.Socket.SetWriteDeadline(time.Now().Add(10 * time.Second))
	}

	return conn.Socket.WriteMessage(websocket.TextMessage, payload)
}

// SendToUsername finds a connection by username and sends a message if present.
func (m *Manager) SendToUsername(ctx context.Context, username string, message types.ServerMessage) error {
	return m.Send(ctx, m.FindByUsername(username), message)
}

// FindByUsername looks up a connection by username.
func (m *Manager) FindByUsername(username string) *Connection {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.byUsername[username]
}

// Broadcast sends a server message to all active connections.
func (m *Manager) Broadcast(ctx context.Context, message types.ServerMessage) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, conn := range m.connections {
		if err := m.Send(ctx, conn, message); err != nil {
			log.Printf("ws: broadcast error id=%s err=%v", conn.ID, err)
		}
	}
}

// Shutdown closes all active sockets gracefully.
func (m *Manager) Shutdown(ctx context.Context) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for id, conn := range m.connections {
		closeErr := conn.Socket.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""), time.Now().Add(2*time.Second))
		if closeErr != nil {
			log.Printf("ws: close control error id=%s err=%v", id, closeErr)
		}
		_ = conn.Socket.Close()
	}

	m.connections = make(map[string]*Connection)
	m.byUsername = make(map[string]*Connection)
}
