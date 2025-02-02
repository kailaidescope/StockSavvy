package server

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

func (server *Server) InitializeModel() {
	ctx := context.Background()

	// Initialize Gemini client
	client, err := genai.NewClient(ctx, option.WithAPIKey(server.geminiKey))
	if err != nil {
		log.Fatal("Failed to get Gemini client", err)
	}
	server.GeminiClient = client

	// Get model
	model := server.GeminiClient.GenerativeModel("gemini-1.5-flash")
	server.GeminiModel = model
}

func (server *Server) GenerateContent(c *gin.Context) {
	// Get prompt from request
	prompt := c.PostForm("prompt")
	if prompt == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "prompt is required"})
		return
	}

	ctx := context.Background()
	resp, err := server.GeminiModel.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		log.Fatal("Failed to get response from Gemini", err)
	}

	//fmt.Print(resp)

	printResponse(resp)
	c.JSON(http.StatusOK, resp)
}

func printResponse(resp *genai.GenerateContentResponse) {
	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				fmt.Println(part)
			}
		}
	}
	fmt.Println("---")
}
