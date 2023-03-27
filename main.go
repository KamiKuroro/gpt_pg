package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/valyala/fasthttp"
)

const (
	APIBaseURL = "https://api.openai.com/v1/chat/completions"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Choice struct {
	Msg          Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
	Index        int     `json:"index"`
}

type OpenAIRequest struct {
	Messages []Message `json:"messages"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type OpenAIResponse struct {
	ChoiceList []Choice `json:"choices"`
	UsageInfo  Usage    `json:"usage"`
}

var APIKey string

func main() {
	apiKey := flag.String("key", "", "API key for the service")
	// Parse the command-line arguments
	flag.Parse()
	// Check if the API key was provided
	if *apiKey == "" {
		fmt.Println("Please provide an API key using the -key flag")
		os.Exit(1)
	}
	APIKey = *apiKey
	router := gin.Default()
	// Configure CORS settings
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowWildcard:    true,
		MaxAge:           12 * time.Hour,
	}))
	router.Static("/", "dist")
	router.POST("/gpt", handleGPTRequest)
	_ = router.Run(":8080")
}

func handleGPTRequest(c *gin.Context) {
	var openAIReq OpenAIRequest

	if err := c.BindJSON(&openAIReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := generateText(openAIReq.Messages)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"response": response})
}

func generateText(messages []Message) (string, error) {
	requestBody, err := json.Marshal(map[string]interface{}{
		"messages":    messages,
		"model":       "gpt-3.5-turbo",
		"max_tokens":  512,
		"temperature": 0.7,
	})
	if err != nil {
		return "", err
	}
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(APIBaseURL)
	req.Header.SetContentType("application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", APIKey))
	req.SetBody(requestBody)
	req.Header.SetMethod("POST")

	resp := fasthttp.AcquireResponse()
	client := &fasthttp.Client{}
	err = client.Do(req, resp)

	if err != nil {
		return "", err
	}

	if resp.StatusCode() != http.StatusOK {
		return "", fmt.Errorf("API request failed with status %d, %v", resp.StatusCode(), resp)
	}

	bodyBytes := resp.Body()
	var openAIRes OpenAIResponse
	err = json.Unmarshal(bodyBytes, &openAIRes)
	if err != nil {
		return "", err
	}
	b, err := json.Marshal(openAIRes.UsageInfo)
	if err != nil {
		return "", err
	}
	log.Println(string(b))
	if len(openAIRes.ChoiceList) == 0 {
		return "", fmt.Errorf("no choices returned from API")
	}
	return openAIRes.ChoiceList[0].Msg.Content, nil
}
