package server

import (
	"errors"
	"financial-helper/mongodb"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/generative-ai-go/genai"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
)

var THROTTLE_TIME time.Duration = 2

type Server struct {
	Router            *gin.Engine
	GeminiClient      *genai.Client
	GeminiModel       *genai.GenerativeModel
	nytKey            string
	polygonKeys       []string
	currentPolygonKey int
	geminiKey         string
	mongoClient       *mongo.Client
}

func (server *Server) GetPolygonKey() string {
	server.currentPolygonKey = (server.currentPolygonKey + 1) % len(server.polygonKeys)
	return server.polygonKeys[server.currentPolygonKey]
}

func GetNewServer() (*Server, error) {
	// Load env vars, if not in Docker container
	log.Println("Checking if Docker is running...")
	if dockerRunning := os.Getenv("DOCKER_RUNNING"); dockerRunning != "true" {
		if err := LoadEnvironment(); err != nil {
			return nil, errors.Join(errors.New("couldn't load environment"), err)
		}
	}

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
	newPolygonKey = os.Getenv("POLYGON_API_KEY_5")
	if newPolygonKey == "" {
		return nil, errors.New("polygon api key 5 not found")
	}
	polygonKeys = append(polygonKeys, newPolygonKey)

	geminiKey := os.Getenv("GOOGLE_GEMINI_API_KEY")
	if geminiKey == "" {
		return nil, errors.New("gemini api key not found")
	}

	mongoUsername := os.Getenv("MONGO_INITDB_ROOT_USERNAME")
	if mongoUsername == "" {
		return nil, errors.New("mongo username not found")
	}

	mongoPassword := os.Getenv("MONGO_INITDB_ROOT_PASSWORD")
	if mongoPassword == "" {
		return nil, errors.New("mongo password not found")
	}

	mongoPortString := os.Getenv("MONGO_PORT")
	if mongoPortString == "" {
		return nil, errors.New("mongo port not found")
	}

	mongoPort, err := strconv.Atoi(mongoPortString)
	if err != nil {
		return nil, errors.Join(errors.New("failed to convert mongo port to int"), err)
	}

	// Initialize router
	router := gin.Default()

	// Set CORS rules
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	router.Use(cors.New(config))

	// Initialize MongoDB connection
	mongoClient, err := mongodb.GetMongoDBInstance(mongoUsername, mongoPassword, mongoPort)
	if err != nil {
		return nil, errors.Join(errors.New("failed to initialize mongodb connection"), err)
	}

	server := &Server{
		Router:            router,
		nytKey:            nytKey,
		polygonKeys:       polygonKeys,
		currentPolygonKey: 0,
		geminiKey:         geminiKey,
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

// LoadEnvironment gets environment variables from a .env file in the project directory and allows us to use them
// Falls through if in Docker environment, to protect against infinite loops
func LoadEnvironment() error {
	log.Println("Double-checking if Docker is running...")
	// Safeguard to prevent infinite loops
	if dockerRunning := os.Getenv("DOCKER_RUNNING"); dockerRunning == "true" {
		return nil
	}

	log.Println("Loading environment variables...")
	// Find root directory
	for {
		workingDirectoryPath, err := os.Getwd()
		if err != nil {
			return errors.Join(errors.New("could not get current working directory for environment file loading"), err)
		}
		_, workingDirectoryName := path.Split(workingDirectoryPath)
		log.Println("Current directory:", workingDirectoryName)
		if workingDirectoryName == "backend" {
			break
		}
		os.Chdir("..")
	}
	err := godotenv.Load(".env", "mongo.env")
	if err != nil {
		return errors.Join(errors.New("could not load environment variables"), err)
	}
	log.Println("Environment variables loaded.")

	return nil
}
