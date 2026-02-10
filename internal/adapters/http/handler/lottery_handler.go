package handler

import (
	"net/http"

	"github.com/backend-challenge/user-api/internal/ports"
	"github.com/gin-gonic/gin"
)

type LotteryHandler struct {
	service ports.LotteryService
}

func NewLotteryHandler(service ports.LotteryService) *LotteryHandler {
	return &LotteryHandler{
		service: service,
	}
}

func (h *LotteryHandler) Search(c *gin.Context) {
	pattern := c.Query("pattern")
	if pattern == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "pattern is required"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	tickets, err := h.service.SearchLottery(c.Request.Context(), pattern, userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"pattern": pattern,
		"results": tickets,
		"count":   len(tickets),
	})
}
