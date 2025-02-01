package server

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetTickerInfo returns information about a stock
//
// GET /api/v1/stocks/tickers/{symbol}
//
// Input:
//   - symbol: the ticker's symbol
//
// Output:
//   - TickerInfo: the ticker information struct
func (server *Server) GetTickerInfo(c *gin.Context) {
	symbol := c.Param("symbol")
	url := fmt.Sprintf("https://api.polygon.io/v3/reference/tickers?ticker=%s&active=true&limit=100&apiKey=%s", symbol, server.polygonKey)
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println("Error generating request for Polygon.io", err)
		return
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending/receiving request to Polygon.io", err)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error reading Polygon.io response", err)
		return
	}
	fmt.Println(string(body))

	var result map[string]interface{}
	if err := c.ShouldBindJSON(&result); err != nil {
		fmt.Println("Error unmarshalling Polygon.io response", err)
		return
	}
	results, ok := result["results"].([]interface{})
	if !ok || len(results) == 0 {
		fmt.Println("Error: results not found or empty")
		return
	}
	firstResult, ok := results[0].(map[string]interface{})
	if !ok {
		fmt.Println("Error: first result is not a map")
		return
	}
	fmt.Println(firstResult["name"])

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

// GetTickerHistory returns the historical prices of a stock
//
// GET /api/v1/stocks/tickers/{symbol}/history
//
// Input:
//   - symbol: the ticker's symbol
//
// Output:
//   - TickerHistory: the ticker history struct
func (server *Server) GetTickerHistory(c *gin.Context) {
	c.JSON(http.StatusOK, testTickerHistory)
}
