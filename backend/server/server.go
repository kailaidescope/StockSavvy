package server

import (
	"errors"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type Server struct {
	Router     *gin.Engine
	nytKey     string
	polygonKey string
}

func GetNewServer() (*Server, error) {
	// Load env vars
	godotenv.Load()
	nytKey := os.Getenv("NYT_API_KEY")
	if nytKey == "" {
		return nil, errors.New("new york times api key not found")
	}

	polygonKey := os.Getenv("POLYGON_API_KEY")
	if polygonKey == "" {
		return nil, errors.New("polygon api key not found")
	}

	// Initialize router
	router := gin.Default()

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	router.Use(cors.New(config))

	server := &Server{
		Router:     router,
		nytKey:     nytKey,
		polygonKey: polygonKey,
	}

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

				// Returns the holdings of a user
				stocks.GET("/holdings", server.GetHoldings)
			}

			/* // Contains all routes relating to logging in and authenticating a user
			auth := v1.Group("/auth")
			{
				// Takes login info and returns an auth token as a cookie
				auth.POST("/login", server.LogIn)
				// Takes profile info and adds it to the database if valid
				auth.POST("/signup", server.SignUp)
				//Deletes auth token on the client side
				auth.POST("/logout", server.LogOut)
			}

			// All these routes require an auth token, i.e. you need to be logged in
			// Contains routes relating to the current user
			user := v1.Group("/user")
			{
				// Gets general information about you as a user
				user.GET("", server.NotImplemented)
				// Changes your user information
				user.PUT("", server.NotImplemented)
				// Deletes your account
				user.DELETE("", server.NotImplemented)

				// Gets your list of languages
				user.GET("/languages", server.NotImplemented)

				// Relates to users that follow you
				followers := user.Group("/followers")
				{
					// Gets all users that follow you
					followers.GET("", server.NotImplemented)
					// Removes a current follower from your followers list
					followers.DELETE("", server.NotImplemented)

					// Relates to inbound follow requests
					requests := followers.Group("/requests")
					{
						// Gets all active requests
						requests.GET("", server.NotImplemented)
						// Takes "resolution": "Accepted" or "Not accepted" and handles the request accordingly
						requests.POST("", server.resolveFollowRequest)
					}
				}

				following := user.Group("/following")
				{
					// Gets all the users that you follow
					following.GET("", server.NotImplemented)
					// Takes a user as a parameter and sends a follow request to that user
					following.POST("", server.followUser)
					// Takes a users as a parameter and removes them from your following list
					following.DELETE("", server.unfollowUser)
				}
			}

			// Contains routes related to searching information about another user
			profiles := v1.Group("/profiles")
			{
				// searchProfile = the username of the user in question
				searchProfile := profiles.Group("/:searchProfile")
				{
					// Gets general information about this user
					searchProfile.GET("", server.GetProfile)

					// Gets a list of the user's languages
					searchProfile.GET("/languages", server.NotImplemented)
					// Gets a list of the user's followers
					searchProfile.GET("/followers", server.shouldAuthChecker, server.privacyValidator, server.getFollowers)
					// Gets a list of people who follow this user
					searchProfile.GET("/following", server.shouldAuthChecker, server.privacyValidator, server.getFollowing)
				}
			}

			// Test function to see if auth works
			v1.GET("/test", server.shouldAuthChecker, server.Testing) */
		}
	}

	// load router
	// load token maker
	return server, nil
}

func (server *Server) NotImplemented(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"status": "This resource is not yet implemented, but will be in the future"})
}
