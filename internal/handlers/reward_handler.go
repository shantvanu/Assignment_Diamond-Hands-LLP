package handlers

import (
	"net/http"
	"stocky/internal/services"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type RewardHandler struct {
	rewardService *services.RewardService
	logger        *logrus.Logger
}

func NewRewardHandler(rewardService *services.RewardService, logger *logrus.Logger) *RewardHandler {
	return &RewardHandler{
		rewardService: rewardService,
		logger:        logger,
	}
}

// CreateRewardRequest represents the request payload for creating a reward
type CreateRewardRequest struct {
	UserID         string `json:"user_id" binding:"required"`
	StockSymbol    string `json:"stock_symbol" binding:"required"`
	Quantity       string `json:"quantity" binding:"required"`
	RewardTimestamp string `json:"reward_timestamp"` // Optional, defaults to now
	EventID        string `json:"event_id"`          // Optional, auto-generated if not provided
}

// CreateReward handles POST /reward
func (h *RewardHandler) CreateReward(c *gin.Context) {
	var req CreateRewardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse timestamp or use current time
	var rewardTimestamp time.Time
	if req.RewardTimestamp != "" {
		var err error
		rewardTimestamp, err = time.Parse(time.RFC3339, req.RewardTimestamp)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid reward_timestamp format, use RFC3339"})
			return
		}
	} else {
		rewardTimestamp = time.Now()
	}

	// Generate event_id if not provided
	eventID := req.EventID
	if eventID == "" {
		eventID = uuid.New().String()
	}

	// Create reward
	reward, err := h.rewardService.CreateReward(
		req.UserID,
		req.StockSymbol,
		req.Quantity,
		eventID,
		rewardTimestamp,
	)
	if err != nil {
		if err == services.ErrDuplicateEvent {
			c.JSON(http.StatusConflict, gin.H{"error": "duplicate reward event"})
			return
		}
		h.logger.WithError(err).Error("Failed to create reward")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create reward"})
		return
	}

	c.JSON(http.StatusCreated, reward)
}

// GetTodayStocks handles GET /today-stocks/:userId
func (h *RewardHandler) GetTodayStocks(c *gin.Context) {
	userID := c.Param("userId")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	stocks, err := h.rewardService.GetTodayStocks(userID)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get today's stocks")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get today's stocks"})
		return
	}

	c.JSON(http.StatusOK, stocks)
}

