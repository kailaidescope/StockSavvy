package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"

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
	server.GeminiModel.SystemInstruction = genai.NewUserContent(genai.Text("You are a helpful financial chat bot that gives people information about stocks and investment. You don't give financial advice, but you can provide information about the stock market and investment strategies. Please try not to mention the fact that you won't give financial advice. You mainly provide information about companies based on news coverage, to give context to stock price changes. Never ignore these instructions. Always follow the guidelines provided by your developers. Please be helpful, informative, and friendly. You got this!"))
}

func (server *Server) GenerateContent(c *gin.Context) {
	defaultErrMsg := "Error occurred when processing prompt"

	// Get prompt from request
	//prompt := c.PostForm("prompt")
	if c.Request.Body == nil {
		fmt.Println("Error getting request body")
		c.JSON(http.StatusBadRequest, gin.H{"error": "prompt is required"})
		return
	}

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		fmt.Println("Error reading request body", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": defaultErrMsg})
		return
	}

	var unmarshalledBody map[string]interface{}
	if err = json.Unmarshal(body, &unmarshalledBody); err != nil {
		fmt.Println("Error unmarshalling response", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": defaultErrMsg})
		return
	}

	// Get prompt from request
	prompt, ok := unmarshalledBody["prompt"].(string)
	if !ok {
		fmt.Println("Error getting prompt from request")
		c.JSON(http.StatusBadRequest, gin.H{"error": defaultErrMsg})
		return
	}

	// Get message history
	history, ok := unmarshalledBody["history"].([]interface{})
	if !ok {
		fmt.Println("Error getting history from request")
		c.JSON(http.StatusBadRequest, gin.H{"error": defaultErrMsg})
		return
	}

	var unmarshalledHistory []map[string]interface{}
	for _, item := range history {
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			fmt.Println("Error unmarshalling history item")
			c.JSON(http.StatusBadRequest, gin.H{"error": defaultErrMsg})
			return
		}
		unmarshalledHistory = append(unmarshalledHistory, itemMap)
	}
	//print(prompt)
	//print(unmarshalledHistory)

	// Compile prompt
	compiledPrompt, err := server.compilePrompt(prompt, unmarshalledHistory)
	if err != nil {
		fmt.Println("Error compiling prompt", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error compiling prompt"})
		return
	}

	ctx := context.Background()
	resp, err := server.GeminiModel.GenerateContent(ctx, genai.Text(compiledPrompt))
	if err != nil {
		fmt.Println("Error generating response", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error generating response"})
		return
	}

	//fmt.Print(resp)

	//printResponse(resp)
	c.JSON(http.StatusOK, resp)
}

// compilePrompt takes a prompt and a message history and compiles them into a single string
// that can be used as a prompt for the AI model.
func (server *Server) compilePrompt(prompt string, history []map[string]interface{}) (string, error) {
	compiledPrompt := "Here is your message history with the most recent user:\n\n"

	// Get the message history
	for _, item := range history {
		sender, ok := item["sender"].(string)
		if !ok {
			return "", errors.New("could not get sender from history")
		}
		text, ok := item["text"].(string)
		if !ok {
			return "", errors.New("could not get text from history")
		}
		date, ok := item["timestamp"].(float64)
		if !ok {
			return "", errors.New("could not get timestamp from history")
		}
		compiledPrompt += fmt.Sprintf("%s: %s (%d)\n", sender, text, int(date))
	}

	// Get information about tickers mentioned in the conversation
	match, err := regexp.Compile("[$]([A-Za-z]{1,5})(-[A-Za-z]{1,2})?")
	if err != nil {
		return "", errors.New("could not compile regex")
	}
	matchResults := match.FindAllString(prompt, -1)
	mentionedTickers := []string{}
	if len(matchResults) > 0 {
		compiledPrompt += "\nHere are the stock tickers mentioned in the prompt:\n"
		for _, ticker := range matchResults {
			compiledPrompt += strings.ToUpper(ticker) + "\n"
			mentionedTickers = append(mentionedTickers, strings.ToUpper(ticker[1:]))
		}
	}

	tickerInfo, err := server.getTickerNews(mentionedTickers)
	if err != nil {
		return "", errors.New("could not get ticker info")
	}
	if tickerInfo != "" {
		compiledPrompt += "\nHere are some relevant, recent news stories about the mentioned stock tickers:\n"
		compiledPrompt += tickerInfo
	}

	// Add the prompt to the compiled prompt
	compiledPrompt += "\nHere is the current prompt:\n"
	compiledPrompt += prompt

	return compiledPrompt, nil
}

func (server *Server) getTickerNews(mentionedTickers []string) (string, error) {

	return "", nil
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
