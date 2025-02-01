package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

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

	defaultErrMsg := "Error receiving ticker info"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println("Error generating request for Polygon.io", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": defaultErrMsg})
		return
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending/receiving request to Polygon.io", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": defaultErrMsg})
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error reading Polygon.io response", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": defaultErrMsg})
		return
	}
	//fmt.Println(string(body))

	// Unmarshall the unmarshalledBodyà
	var unmarshalledBody map[string]interface{}
	if err = json.Unmarshal(body, &unmarshalledBody); err != nil {
		fmt.Println("Error unmarshalling response", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": defaultErrMsg})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": defaultErrMsg})
		return
	}

	// Convert each element to map[string]interface{}
	var convertedResults []map[string]interface{}
	for _, result := range results {
		resultMap, ok := result.(map[string]interface{})
		if !ok {
			fmt.Println("Error: result element is not a map")
			c.JSON(http.StatusInternalServerError, gin.H{"error": defaultErrMsg})
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
	time.Sleep(12 * time.Second)
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
	symbol := c.Param("symbol")
	url := fmt.Sprintf("https://api.polygon.io/v1/indicators/sma/%s?timespan=day&adjusted=true&window=20&series_type=close&order=asc&limit=100&apiKey=%s", symbol, server.polygonKey)
	method := "GET"

	defaultErrMsg := "Error receiving ticker history"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println("Error generating request", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": defaultErrMsg})
		return
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": defaultErrMsg})
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error reading response body", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": defaultErrMsg})
		return
	}
	//fmt.Println(string(body))

	// Unmarshall the unmarshalledBody
	var unmarshalledBody map[string]interface{}
	if err = json.Unmarshal(body, &unmarshalledBody); err != nil {
		fmt.Println("Error unmarshalling response", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": defaultErrMsg})
		return
	}

	// Check if results exist
	if _, exists := unmarshalledBody["results"]; !exists {
		fmt.Println("Error: results not found")
		c.JSON(http.StatusInternalServerError, gin.H{"error": defaultErrMsg})
		return
	}

	// Convert to correct data types
	results, ok := unmarshalledBody["results"].(map[string]interface{})
	if !ok {
		fmt.Println("Error: results not converted")
		c.JSON(http.StatusInternalServerError, gin.H{"error": defaultErrMsg})
		return
	}

	// Check if values exist
	if _, exists := results["values"]; !exists {
		fmt.Println("Error: values not found")
		c.JSON(http.StatusInternalServerError, gin.H{"error": defaultErrMsg})
		return
	}

	// Convert to a history list
	values, ok := results["values"].([]interface{})
	if !ok {
		fmt.Println("Error: values not converted")
		c.JSON(http.StatusInternalServerError, gin.H{"error": defaultErrMsg})
		return
	}

	// Convert each element to map[string]interface{}
	var convertedValues []map[string]interface{}
	for _, result := range values {
		resultMap, ok := result.(map[string]interface{})
		if !ok {
			fmt.Println("Error: result element is not a map")
			c.JSON(http.StatusInternalServerError, gin.H{"error": defaultErrMsg})
			return
		}
		convertedValues = append(convertedValues, resultMap)
	}

	for _, value := range convertedValues {
		if timestamp, exists := value["timestamp"]; exists {
			value["time"] = timestamp
			delete(value, "timestamp")
		}
	}

	history := TickerHistory{
		History: convertedValues,
	}

	c.JSON(http.StatusOK, history)
	time.Sleep(12 * time.Second)
}

// GetTickerNews returns the news sentiment of a stock
//
// GET /api/v1/stocks/tickers/{symbol}/news
//
// Input:
//   - symbol: the ticker's symbol
//
// Output:
//   - TickerNews: the ticker news struct
func (server *Server) GetTickerNews(c *gin.Context) {
	symbol := c.Param("symbol")
	url := fmt.Sprintf("https://api.polygon.io/v2/reference/news?ticker=%s&order=desc&limit=350&sort=published_utc&apiKey=%s&published_utc.gte=2024-10-11T19:01:33Z", symbol, server.polygonKey)
	method := "GET"

	defaultErrMsg := "Error receiving ticker news"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println("Error generating request", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": defaultErrMsg})
		return
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": defaultErrMsg})
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error reading response body", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": defaultErrMsg})
		return
	}
	//fmt.Println(string(body))

	// Unmarshall the unmarshalledBody
	var unmarshalledBody map[string]interface{}
	if err = json.Unmarshal(body, &unmarshalledBody); err != nil {
		fmt.Println("Error unmarshalling response", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": defaultErrMsg})
		return
	}

	// Convert to correct data types
	results, ok := unmarshalledBody["results"].([]interface{})
	if !ok {
		fmt.Println("Error: results not converted")
		c.JSON(http.StatusInternalServerError, gin.H{"error": defaultErrMsg})
		return
	}

	// Get number of articles
	numArticles, ok := unmarshalledBody["count"].(float64)
	if !ok {
		fmt.Println("Error: count not converted")
		c.JSON(http.StatusInternalServerError, gin.H{"error": defaultErrMsg})
		return
	}
	numArticlesInt := int(numArticles)

	// Convert each element to map[string]interface{}
	var articles []map[string]interface{}
	for _, result := range results {
		resultMap, ok := result.(map[string]interface{})
		if !ok {
			fmt.Println("Error: result element is not a map")
			c.JSON(http.StatusInternalServerError, gin.H{"error": defaultErrMsg})
			return
		}
		articles = append(articles, resultMap)
	}

	// Calculate average sentiment

	var sentiments []float64
	for _, article := range articles {
		if insights, exists := article["insights"]; exists {
			insightsList, ok := insights.([]interface{})
			if !ok {
				fmt.Println("Error: insights not converted")
				c.JSON(http.StatusInternalServerError, gin.H{"error": defaultErrMsg})
				return
			}

			for _, singleTickerInsight := range insightsList {
				convertedSingleTickerInsight, ok := singleTickerInsight.(map[string]interface{})
				if !ok {
					fmt.Println("Error: result element is not a map")
					c.JSON(http.StatusInternalServerError, gin.H{"error": defaultErrMsg})
					return
				}
				if ticker, exists := convertedSingleTickerInsight["ticker"]; exists {
					if ticker == symbol {
						if sentiment, exists := convertedSingleTickerInsight["sentiment"]; exists {
							sentimentString, ok := sentiment.(string)
							if !ok {
								fmt.Println("Error: sentiment not converted")
								c.JSON(http.StatusInternalServerError, gin.H{"error": defaultErrMsg})
								return
							}

							if sentimentString == "positive" {
								sentiments = append(sentiments, 1)
							} else if sentimentString == "negative" {
								sentiments = append(sentiments, -1)
							} else {
								sentiments = append(sentiments, 0)
							}
						}
					}
				}
			}
		}
	}

	// Calculate average sentiment
	var sumSentiment float64
	for _, sentiment := range sentiments {
		sumSentiment += sentiment
	}
	avgSentiment := sumSentiment / float64(len(sentiments))

	// Calculate standard deviation
	var sumSquaredDifferences float64
	for _, sentiment := range sentiments {
		sumSquaredDifferences += (sentiment - avgSentiment) * (sentiment - avgSentiment)
	}
	stdDevSentiment := sumSquaredDifferences / float64(len(sentiments))

	news := TickerNews{
		AverageSentiment: float32(avgSentiment),
		StdDevSentiment:  float32(stdDevSentiment),
		NumArticles:      numArticlesInt,
	}

	c.JSON(http.StatusOK, news)
	time.Sleep(12 * time.Second)
}

// GetTickerHoldings returns the holdings of a stock
//
// GET /api/v1/stocks/holdings
//
// Output:
//   - TickerHoldings: the ticker holdings struct
func (server *Server) GetHoldings(c *gin.Context) {
	c.JSON(http.StatusOK, testHoldings)
}
