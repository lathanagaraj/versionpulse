package vp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/sashabaranov/go-openai"
)

// Huggingface
// const baseURL = "https://router.huggingface.co/together"
// const modelName = "deepseek-ai/DeepSeek-R1"

// Open router
// const baseURL = "https://openrouter.ai/api/v1"
// const modelName = "deepseek/deepseek-r1:free"

// Azure Open AI
var apiKey = os.Getenv("VP_API_KEY")

const baseURL = "https://version-pulse.openai.azure.com"
const modelName = "gpt-4o"

type Checker struct {
	ToolName string
	Content  string
}

type promptData struct {
	Tool       string
	WebContent string
}
type ToolVersion struct {
	Tool    string `json:"tool"`
	Version string `json:"version"`
	Date    string `json:"date"`
	Link    string
}

func NewChecker(toolName, content string) Checker {
	return Checker{
		ToolName: toolName,
		Content:  content,
	}
}

func (c Checker) CheckVersion() (*ToolVersion, error) {
	response, err := queryLLM(c.ToolName, c.Content)
	if err != nil {
		log.Printf("Error querying LLM: %v", err)
		return nil, err
	}
	toolVersion, err := extractJSONObject(response)
	if err != nil {
		log.Printf("Error extracting JSON object: %v", err)
		return nil, err
	}
	return toolVersion, nil

}

func createAzureClient() *openai.ClientConfig {
	config := openai.DefaultAzureConfig(apiKey, baseURL)
	config.APIType = openai.APITypeAzure
	config.APIVersion = "2024-08-01-preview"
	config.AzureModelMapperFunc = func(model string) string {
		azureModelMapping := map[string]string{
			"gpt-4o": model,
		}
		return azureModelMapping[model]
	}
	config.HTTPClient = createhttpClient()
	return &config
}

// func createClient() *openai.ClientConfig {
// 	config := openai.DefaultConfig(apiKey)
// 	config.BaseURL = baseURL
// 	return &config
// }

func createhttpClient() *http.Client {
	// Configure retryable HTTP client
	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 5
	retryClient.RetryWaitMin = 2 * time.Second
	retryClient.RetryWaitMax = 10 * time.Second
	retryClient.CheckRetry = retryablehttp.DefaultRetryPolicy // Automatic handling of 429 errors

	return retryClient.StandardClient()
}

// Queries an LLM using OpenAI client with a custom base URL
func queryLLM(toolName, extractedText string) (string, error) {

	// Combine extracted text with user query
	userprompt, err := createPrompt(toolName, extractedText)
	if err != nil {
		return "", err
	}

	client := openai.NewClientWithConfig(*createAzureClient())

	// Define chat messages
	messages := []openai.ChatCompletionMessage{
		{
			Role: openai.ChatMessageRoleSystem,
			Content: `
You are a helpful assistant that only responds in valid JSON format.
**Rules:**
1. Your response must **strictly** be in valid JSON format and nothing else.
2. Do not include any additional explanations or text in your response.
3. If the version or date cannot be found, return an empty JSON with the attributes set to null.
			`,
		},
		{
			Role:    openai.ChatMessageRoleUser,
			Content: userprompt,
		},
	}

	// Send request to OpenAI
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    modelName,
			Messages: messages,
			ResponseFormat: &openai.ChatCompletionResponseFormat{
				Type: openai.ChatCompletionResponseFormatTypeJSONObject,
			},
			MaxTokens:   50,
			N:           1,
			Temperature: 0,
		},
	)
	if err != nil {
		return "", fmt.Errorf("API error: %v", err)
	}

	// Return model response
	if len(resp.Choices) > 0 {
		return resp.Choices[0].Message.Content, nil
	}
	return "", fmt.Errorf("no response received from the model")
}

func createPrompt(toolName, extractedText string) (string, error) {
	promptTemplate, err := template.ParseFiles("prompt.txt")
	if err != nil {
		return "", err
	}

	// Define the data for the prompt
	data := promptData{
		Tool:       toolName,
		WebContent: extractedText,
	}

	// Execute the template with the data
	var buf bytes.Buffer
	err = promptTemplate.Execute(&buf, data)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	return buf.String(), nil
}

func extractJSONObject(text string) (*ToolVersion, error) {
	// Find the JSON object in the text
	startIndex := strings.Index(text, "{")
	if startIndex == -1 {
		return nil, fmt.Errorf("no JSON object found in text")
	}

	endIndex := strings.Index(text, "}")
	if endIndex == -1 {
		return nil, fmt.Errorf("no closing bracket found in JSON object")
	}

	// Extract the JSON object
	jsonString := text[startIndex : endIndex+1]
	var toolVersion ToolVersion

	err := json.Unmarshal([]byte(jsonString), &toolVersion)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON object: %v", err)
	}

	return &toolVersion, nil
}
