# Connect Four Backend Skeleton

This repository contains the Phase B1 backend scaffold for the realtime Connect Four project. The goal of this phase is to provide a WebSocket gateway, connection management, and message contracts compatible with the existing frontend.

## Prerequisites

- Go 1.22+
- (Optional) Docker + Docker Compose for local infrastructure stubs

## Running the Server

```bash
cd backend
go run ./cmd/server
```

By default the API listens on `http://localhost:8080`. The WebSocket endpoint is exposed at `ws://localhost:8080/ws?username=<name>`.

## Message Contract

The server speaks the same protocol as the frontend:

- **Client → Server**
  - `MAKE_MOVE { col: number, gameId: string }`
  - `RECONNECT { username: string, gameId?: string }`
- **Server → Client**
  - `GAME_START { gameId: string, you: 1|2, opponent: string }`
  - `BOARD_UPDATE { board: number[][], currentTurn: 1|2 }`
  - `GAME_OVER { result: "WIN"|"LOSS"|"DRAW", board: number[][] }`
  - `INFO { message: string }`

Phase B1 responds with `INFO` messages acknowledging valid actions and emits validation errors for malformed payloads.

## Sample WebSocket Client

```js
const socket = new WebSocket('ws://localhost:8080/ws?username=alice');

socket.addEventListener('open', () => {
  console.log('connected');
  socket.send(JSON.stringify({ type: 'MAKE_MOVE', gameId: 'game-1', col: 3 }));
});

socket.addEventListener('message', (event) => {
  const data = JSON.parse(event.data);
  console.log('message from server', data);
});

socket.addEventListener('close', () => {
  console.log('connection closed');
});
```

## Testing

```bash
cd backend
go test ./...
```

The included test spins up an in-memory server, submits a `MAKE_MOVE` request, and expects an `INFO` acknowledgment.

## Docker Compose Stubs

A `docker-compose.yml` is provided with placeholders for PostgreSQL, Redis, and Kafka. They are not required for Phase B1 but document the infrastructure planned for later phases.
