package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"     // ← correct position here

	"github.com/example/connect-four/backend/internal/api"
	"github.com/example/connect-four/backend/internal/bot"
	"github.com/example/connect-four/backend/internal/game"
	"github.com/example/connect-four/backend/internal/matchmaking"
	"github.com/example/connect-four/backend/internal/store"
	"github.com/example/connect-four/backend/internal/ws"
)

func main() {
	// Load environment variables from backend/.env automatically
	godotenv.Load()

	logger := log.New(os.Stdout, "", log.LstdFlags|log.Lmicroseconds)
	log.SetOutput(logger.Writer())

	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	db, err := store.NewPostgres()
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("database close error: %v", err)
		}
	}()

	repo := store.NewRepository(db)
	apiHandlers := api.New(repo)

	manager := ws.NewManager()
	gameManager := game.NewManager()
	matchmaker := matchmaking.NewMatchmaker(gameManager, manager, "BOT")
	botEngine := bot.New(gameManager)
	handler := ws.NewHandler(manager, gameManager, matchmaker, botEngine, repo)
	handler.RegisterRoutes(r)
	r.GET("/leaderboard", apiHandlers.GetLeaderboard)

	srv := newHTTPServer(r)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go matchmaker.Start(ctx)

	go func() {
		if err := srv.run(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("shutting down...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.shutdown(shutdownCtx); err != nil {
		log.Printf("http shutdown error: %v", err)
	}

	manager.Shutdown(shutdownCtx)
}

type httpServer struct {
	engine *gin.Engine
	server *http.Server
}

func newHTTPServer(engine *gin.Engine) *httpServer {
	return &httpServer{
		engine: engine,
		server: &http.Server{
			Addr:    ":8080",
			Handler: engine,
		},
	}
}

func (s *httpServer) run() error {
	log.Println("server listening on :8080")
	return s.server.ListenAndServe()
}

func (s *httpServer) shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
