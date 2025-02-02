package server

import (
	"errors"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/generative-ai-go/genai"
	"github.com/joho/godotenv"
)

type Server struct {
	Router            *gin.Engine
	GeminiClient      *genai.Client
	GeminiModel       *genai.GenerativeModel
	nytKey            string
	polygonKeys       []string
	currentPolygonKey int
	geminiKey         string
}

func (server *Server) GetPolygonKey() string {
	server.currentPolygonKey = (server.currentPolygonKey + 1) % len(server.polygonKeys)
	return server.polygonKeys[server.currentPolygonKey]
}

func GetNewServer() (*Server, error) {
	// Load env vars
	godotenv.Load()
	nytKey := os.Getenv("NYT_API_KEY")
	if nytKey == "" {
		return nil, errors.New("new york times api key not found")
	}

	polygonKeys := []string{}
	newPolygonKey := os.Getenv("POLYGON_API_KEY")
	if newPolygonKey == "" {
		return nil, errors.New("polygon api key not found")
	}
	polygonKeys = append(polygonKeys, newPolygonKey)
	newPolygonKey = os.Getenv("POLYGON_API_KEY_1")
	if newPolygonKey == "" {
		return nil, errors.New("polygon api key 1 not found")
	}
	polygonKeys = append(polygonKeys, newPolygonKey)
	newPolygonKey = os.Getenv("POLYGON_API_KEY_2")
	if newPolygonKey == "" {
		return nil, errors.New("polygon api key 2 not found")
	}
	polygonKeys = append(polygonKeys, newPolygonKey)
	newPolygonKey = os.Getenv("POLYGON_API_KEY_3")
	if newPolygonKey == "" {
		return nil, errors.New("polygon api key 3 not found")
	}
	polygonKeys = append(polygonKeys, newPolygonKey)
	newPolygonKey = os.Getenv("POLYGON_API_KEY_4")
	if newPolygonKey == "" {
		return nil, errors.New("polygon api key 4 not found")
	}
	polygonKeys = append(polygonKeys, newPolygonKey)

	geminiKey := os.Getenv("GOOGLE_GEMINI_API_KEY")
	if geminiKey == "" {
		return nil, errors.New("gemini api key not found")
	}

	// Initialize router
	router := gin.Default()

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	router.Use(cors.New(config))

	server := &Server{
		Router:            router,
		nytKey:            nytKey,
		polygonKeys:       polygonKeys,
		currentPolygonKey: 0,
		geminiKey:         geminiKey,
	}

	server.InitializeModel()

	// Mount routes
	api := router.Group("/api")
	{
		// Gives information about the API in general, particularly about how to switch between versions
		api.GET("", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "This is the API root. To access a specific version, add /v<version_number> to the end of the URL. Available versions: v1"})
		})

		// Contains all routes relating to version 1 of the API
		v1 := api.Group("/v1")
		{
			// Gives information about this particular version
			v1.GET("", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"status": "This is version 1 of the API. It is the current version."})
			})

			// Contains all routes relating to stocks
			stocks := v1.Group("/stocks")
			{

				// Contains all routes relating to tickers
				tickers := stocks.Group("/tickers")
				{

					// Contains all routes relating to a specific ticker
					searchTicker := tickers.Group("/:symbol")
					{
						// Returns information about a ticker
						searchTicker.GET("", server.GetTickerInfo)

						// Returns the historical prices of a ticker
						searchTicker.GET("/history", server.GetTickerHistory)

						// Returns the news sentiment of a ticker
						searchTicker.GET("/news", server.GetTickerNews)
					}
				}

				// Contains all routes relating to holdings
				holdings := stocks.Group("/holdings")
				{
					// Returns all the holdings of a user
					holdings.GET("", server.GetHoldings)

					// Returns historical data about a holding
					holdings.GET("/:symbol", server.GetHoldingInfo)
				}
			}

			// Contains all routes relating to the AI chat
			chat := v1.Group("/chat")
			{
				// Returns a response from the AI chat
				chat.POST("", server.GenerateContent)
			}
		}
	}

	// load router
	// load token maker
	return server, nil
}

func (server *Server) NotImplemented(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"status": "This resource is not yet implemented, but will be in the future"})
}
