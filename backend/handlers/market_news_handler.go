package handlers

import (
	"isxportfolio-backend/scraper"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MarketNewsHandler struct {
	scraper *scraper.MarketNewsScraper
}

func NewMarketNewsHandler() *MarketNewsHandler {
	return &MarketNewsHandler{
		scraper: scraper.NewMarketNewsScraper(
			"/app/data/market_news.csv", // Updated path for Docker
			"/app/data/pdfs",            // Updated path for Docker
		),
	}
}

// GetMarketNews handles GET /api/market/news
func (h *MarketNewsHandler) GetMarketNews(c *gin.Context) {
	// Run the scraper to get latest news
	err := h.scraper.Run()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch market news",
		})
		return
	}

	// Return all items from the scraper
	c.JSON(http.StatusOK, h.scraper.AllItems)
}

// RefreshMarketNews handles POST /api/market/news/refresh
func (h *MarketNewsHandler) RefreshMarketNews(c *gin.Context) {
	err := h.scraper.Run()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to refresh market news",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Market news refreshed successfully",
	})
}
