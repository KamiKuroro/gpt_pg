package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/valyala/fasthttp"
)

const (
	APIKey     = ""
	APIBaseURL = "https://api.openai.com/v1/engines/text-davinci-002/completions"
)

type OpenAIRequest struct {
	Prompt string `json:"prompt"`
}

type OpenAIResponse struct {
	Choices []struct {
		Text string `json:"text"`
	} `json:"choices"`
}

func main() {
	router := gin.Default()
	router.POST("/gpt", handleGPTRequest)
	router.Run(":8080")
}

func handleGPTRequest(c *gin.Context) {
	var openAIReq OpenAIRequest

	if err := c.BindJSON(&openAIReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := generateText(openAIReq.Prompt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"response": response})
}

func generateText(prompt string) (string, error) {
	requestBody, err := json.Marshal(map[string]interface{}{
		"prompt":   prompt,
		"max_tokens": 50,
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

	if len(openAIRes.Choices) == 0 {
		return "", fmt.Errorf("No choices returned from API")
	}

	return openAIRes.Choices[0].Text, nil
}

