package server

import (
	"context"
	"fmt"
	"log"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

func (server *Server) InitializeModel() {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(server.geminiKey))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()
	server.GeminiClient = client

	// [START text_gen_text_only_prompt]
	model := server.GeminiClient.GenerativeModel("gemini-1.5-flash")
	server.GeminiModel = model
	resp, err := server.GeminiModel.GenerateContent(ctx, genai.Text("Write a story about a magic backpack."))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(resp)

	//printResponse(resp)
}
