package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"text/template"

	"github.com/PuerkitoBio/goquery"
	"github.com/gorilla/feeds"
	"github.com/hashicorp/go-retryablehttp"
	openai "github.com/sashabaranov/go-openai"
	"gopkg.in/yaml.v2"
)

// Hugging facse
// const baseURL = "https://router.huggingface.co/together"
// const modelName = "deepseek-ai/DeepSeek-R1"

// Open router
// const baseURL = "https://openrouter.ai/api/v1"
// const modelName = "deepseek/deepseek-r1:free"

// Azure Open AI
const apiKey = ""
const baseURL = "https://version-pulse.openai.azure.com"
const modelName = "gpt-4o"

type Tool struct {
	Name string `yaml:"name"`
	URL  string `yaml:"url"`
}
type Tools struct {
	Tools []Tool `yaml:"tools"`
}

type PromptData struct {
	Tool       string
	WebContent string
}

type ToolVersion struct {
	Tool    string `json:"tool"`
	Version string `json:"version"`
	Date    string `json:"date"`
	Link    string
}

// Extracts text from an HTML string
func extractTextFromHTML(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Error fetching %s: %v", url, err)
		return "", err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Printf("Error parsing HTML for %s: %v", url, err)
		return "", err
	}

	text := doc.Text()
	text = strings.ReplaceAll(text, "\n", " ")
	text = strings.ReplaceAll(text, "\t", " ")
	text = strings.TrimSpace(text)

	return text, nil
}

// Queries an LLM using OpenAI client with a custom base URL
func queryLLM(apiKey, baseURL, toolName, extractedText string) (string, error) {

	// Combine extracted text with user query
	fullPrompt, err := createPrompt(toolName, extractedText)
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
			Content: fullPrompt,
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
	data := PromptData{
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

	// Print the resolved prompt
	//fmt.Println(buf.String())

	return buf.String(), nil
}

func checkToolVersions() {

	// Load tools.yaml file
	data, err := ioutil.ReadFile("tools.yaml")
	if err != nil {
		log.Fatal(err)
	}

	// Unmarshal YAML data into a slice of Tool structs
	var toolsMap Tools
	err = yaml.Unmarshal(data, &toolsMap)
	if err != nil {
		log.Fatal(err)
	}

	// Fetch HTML pages for each URL and extract text content

	var toolVersions []ToolVersion
	for _, tool := range toolsMap.Tools {
		siteContent, err := extractTextFromHTML(tool.URL)
		if err != nil {
			log.Printf("Error fetching %s: %v", tool.URL, err)
			continue
		}

		response, err := queryLLM(apiKey, baseURL, tool.Name, siteContent)
		if err != nil {
			log.Printf("Error querying LLM: %v", err)
		}
		toolVersion, err := extractJSONObject(response)

		if err != nil {
			log.Printf("Error extracting JSON object: %v", err)
		}
		toolVersion.Link = tool.URL
		toolVersions = append(toolVersions, *toolVersion)

		println("result " + toolVersion.Tool + " " + toolVersion.Version + " " + toolVersion.Date + " " + toolVersion.Link)
	}

	err = generateRSSFeed(toolVersions)
	if err != nil {
		log.Printf("Error generating RSS feed: %v", err)
	}
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

func main() {

	checkToolVersions()
	//test()

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

func createClient() *openai.ClientConfig {
	config := openai.DefaultConfig(apiKey)
	config.BaseURL = baseURL
	return &config
}

func createhttpClient() *http.Client {
	// Configure retryable HTTP client
	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 5
	retryClient.RetryWaitMin = 2 * time.Second
	retryClient.RetryWaitMax = 10 * time.Second
	retryClient.CheckRetry = retryablehttp.DefaultRetryPolicy // Automatic handling of 429 errors

	return retryClient.StandardClient()
}

func generateRSSFeed(toolVersions []ToolVersion) error {
	// Create a new RSS feed
	feed := &feeds.Feed{
		Title:       "Latest Tool Versions",
		Link:        &feeds.Link{Href: "https://example.com"},
		Description: "Latest Tool Versions RSS feed",
		Created:     time.Now(),
	}

	// Convert JSON items to RSS items
	for _, toolVersion := range toolVersions {

		rssItem := &feeds.Item{
			Title:       toolVersion.Tool,
			Link:        &feeds.Link{Href: toolVersion.Link},
			Description: toolVersion.Tool + " " + toolVersion.Version + " " + toolVersion.Date,
			Created:     time.Now(),
		}
		feed.Items = append(feed.Items, rssItem)
	}

	// Convert feed to RSS format
	rss, err := feed.ToRss()
	if err != nil {
		return fmt.Errorf("error generating RSS: %v", err)
	}

	os.MkdirAll("docs", os.ModePerm)

	// Write RSS to file
	err = ioutil.WriteFile("docs/feed.rss", []byte(rss), 0644)
	if err != nil {
		return fmt.Errorf("error writing RSS file: %v", err)
	}

	fmt.Println("RSS feed successfully generated as feed.rss")
	return nil
}
