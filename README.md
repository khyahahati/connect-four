# Connect Four

## Project Overview
This is my real-time Connect Four game that you can play either online against other people or hop in against my built-in bot. I wired it up with a Go backend and a React/Vite frontend, and the live matches run through WebSockets so turns update instantly. I stash finished games and the shared leaderboard in Postgres so stats stick around between sessions.

Live App: https://frontend-production-b8a6.up.railway.app


## Features
- Real-time multiplayer
- Matchmaking + bot fallback
- Game reconnect support
- Local mode for testing
- Persistent leaderboard
- Postgres database
- WebSocket based gameplay

## Tech Stack
- Go (Gin, WebSockets)
- React + Vite
- Postgres
- Docker (for local DB)

## How to Run the Project
### Frontend Setup
```
cd frontend
npm install
npm run dev
```
Then open http://localhost:5173

### Backend Setup
```
cd backend
go mod tidy
go run cmd/server/main.go
```
The backend listens on port 8080.

## How to Start Postgres with Docker
```
docker-compose up -d postgres
```
To hop into the container:
```
docker exec -it backend-postgres-1 bash
psql -U connectfour -d connectfour
```
Check the tables with:
```
\dt
```

## How to Apply Migrations
Copy migrations.sql into the container, then run:
```
psql -U connectfour -d connectfour -f /migrations.sql
```

## Environment Variables
Backend .env:
```
DB_HOST=localhost
DB_PORT=5433
DB_USER=connectfour
DB_PASS=connectfour
DB_NAME=connectfour
```

Frontend .env:
```
VITE_API_WS_URL=ws://localhost:8080/ws
```

## Folder Structure Overview
```
/frontend
/backend
  /cmd
  /internal
  docker-compose.yml
  .env
```

## Deployment
Both the backend and frontend run on Railway.
- Frontend: https://frontend-production-b8a6.up.railway.app
- Backend (HTTP): https://backend-production-3edf2.up.railway.app
- WebSocket endpoint: wss://backend-production-3edf2.up.railway.app/ws
