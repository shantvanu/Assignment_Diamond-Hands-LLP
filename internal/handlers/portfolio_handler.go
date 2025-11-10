package handlers

import (
	"net/http"
	"stocky/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type PortfolioHandler struct {
	portfolioService *services.PortfolioService
	logger           *logrus.Logger
}

func NewPortfolioHandler(portfolioService *services.PortfolioService, logger *logrus.Logger) *PortfolioHandler {
	return &PortfolioHandler{
		portfolioService: portfolioService,
		logger:           logger,
	}
}

// GetHistoricalINR handles GET /historical-inr/:userId
func (h *PortfolioHandler) GetHistoricalINR(c *gin.Context) {
	userID := c.Param("userId")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	historical, err := h.portfolioService.GetHistoricalINR(userID)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get historical INR")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get historical INR"})
		return
	}

	c.JSON(http.StatusOK, historical)
}

// GetStats handles GET /stats/:userId
func (h *PortfolioHandler) GetStats(c *gin.Context) {
	userID := c.Param("userId")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	stats, err := h.portfolioService.GetStats(userID)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get stats")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get stats"})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// GetPortfolio handles GET /portfolio/:userId (Bonus endpoint)
func (h *PortfolioHandler) GetPortfolio(c *gin.Context) {
	userID := c.Param("userId")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	portfolio, err := h.portfolioService.GetPortfolio(userID)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get portfolio")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get portfolio"})
		return
	}

	c.JSON(http.StatusOK, portfolio)
}

