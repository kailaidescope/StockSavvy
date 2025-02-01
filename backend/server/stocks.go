package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type StockInfo struct {
	StockSymbol     string
	CompanyName     string
	Industry        string
	Locale          string
	PrimaryExchange string
}

// GetStockInfo returns information about a stock
//
// GET /api/v1/stocks/:stockSymbol
//
// Input:
//   - stockSymbol: the stock symbol
//
// Output:
//   - StockInfo: the stock information struct
func (server *Server) GetStockInfo(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
