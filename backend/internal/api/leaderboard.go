package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/example/connect-four/backend/internal/store"
)

// LeaderboardRepository describes the storage dependency required by the API.
type LeaderboardRepository interface {
	GetLeaderboard(limit int) ([]store.LeaderboardEntry, error)
}

// API bundles HTTP handlers that depend on the leaderboard repository.
type API struct {
	repo LeaderboardRepository
}

// New constructs a new API surface.
func New(repo LeaderboardRepository) *API {
	return &API{repo: repo}
}

// GetLeaderboard responds with the leaderboard in JSON.
func (a *API) GetLeaderboard(c *gin.Context) {
	if a == nil || a.repo == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "repository unavailable"})
		return
	}

	limit := 10
	if raw := c.Query("limit"); raw != "" {
		parsed, err := strconv.Atoi(raw)
		if err != nil || parsed <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "limit must be a positive integer"})
			return
		}
		limit = parsed
	}

	entries, err := a.repo.GetLeaderboard(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, entries)
}
