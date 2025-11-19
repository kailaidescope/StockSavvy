package server

import (
	"errors"
	"financial-helper/environment"
	"financial-helper/mongodb"
	"financial-helper/polygon"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/generative-ai-go/genai"
	"go.mongodb.org/mongo-driver/mongo"
)

type Server struct {
	Router            *gin.Engine
	GeminiClient      *genai.Client
	GeminiModel       *genai.GenerativeModel
	nytKey            string
	geminiKey         string
	polygonConnection *polygon.PolygonConnection
	mongoClient       *mongo.Client
}

func GetNewServer() (*Server, error) {
	if err := environment.LoadEnvironment(); err != nil {
		return nil, errors.Join(errors.New("couldn't load environment for server"), err)
	}

	vars, polygonKeys, err := environment.LoadVars()
	if err != nil {
		return nil, errors.Join(errors.New("couldn't load env vars"), err)
	}

	// Initialize router
	router := gin.Default()

	// Set CORS rules
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	router.Use(cors.New(config))

	// Initialize MongoDB connection
	mongoPort, _ := strconv.Atoi(vars["MONGO_PORT"]) // Don't need to check that this works because LoadVars() already did
	mongoClient, err := mongodb.GetMongoDBInstance(vars["MONGO_INITDB_ROOT_USERNAME"], vars["MONGO_INITDB_ROOT_PASSWORD"], mongoPort)
	if err != nil {
		return nil, errors.Join(errors.New("failed to initialize mongodb connection"), err)
	}

	// Initialize Polygon connection
	throttleTimeInt, _ := strconv.Atoi(vars["THROTTLE_TIME"]) // Don't need to check that this works because LoadVars() already did
	polygonConnection := polygon.GetPolygonConnection(polygonKeys, time.Duration(throttleTimeInt)*time.Second)

	server := &Server{
		Router:            router,
		nytKey:            vars["NYT_API_KEY"],
		geminiKey:         vars["GOOGLE_GEMINI_API_KEY"],
		polygonConnection: polygonConnection,
		mongoClient:       mongoClient,
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
