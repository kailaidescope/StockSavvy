package server

import (
	"encoding/json"
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error receiving ticker data"})
		return
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending/receiving request to Polygon.io", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error receiving ticker data"})
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error reading Polygon.io response", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error receiving ticker data"})
		return
	}
	fmt.Println(string(body))

	// Unmarshall the unmarshalledBodyà
	var unmarshalledBody map[string]interface{}
	if err = json.Unmarshal(body, &unmarshalledBody); err != nil {
		fmt.Println("Error unmarshalling response", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error receiving ticker data"})
		return
	}

	// Convert to correct data types
	results, ok := unmarshalledBody["results"].([]interface{})
	if !ok || len(results) == 0 {
		if !ok {
			fmt.Println("Error: results not found")
		} else {
			fmt.Println("Error: results empty")
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error receiving ticker data"})
		return
	}

	// Convert each element to map[string]interface{}
	var convertedResults []map[string]interface{}
	for _, result := range results {
		resultMap, ok := result.(map[string]interface{})
		if !ok {
			fmt.Println("Error: result element is not a map")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error receiving ticker data"})
			return
		}
		convertedResults = append(convertedResults, resultMap)
	}

	// Use convertedResults for further processing
	info := TickerInfo{
		Symbol:          fmt.Sprintf("%v", convertedResults[0]["ticker"]),
		Name:            fmt.Sprintf("%v", convertedResults[0]["name"]),
		Industry:        "Not yet set",
		Locale:          fmt.Sprintf("%v", convertedResults[0]["locale"]),
		PrimaryExchange: fmt.Sprintf("%v", convertedResults[0]["primary_exchange"]),
	}

	c.JSON(http.StatusOK, info)
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
