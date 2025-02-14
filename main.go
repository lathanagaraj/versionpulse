package main

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"text/template"

	"github.com/PuerkitoBio/goquery"
	openai "github.com/sashabaranov/go-openai"
	"gopkg.in/yaml.v2"
)

const apiKey = ""
const baseURL = "https://router.huggingface.co/together"

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

	// Configure OpenAI client with custom base URL
	defaultConfig := openai.DefaultConfig(apiKey)
	defaultConfig.BaseURL = baseURL
	client := openai.NewClientWithConfig(defaultConfig)

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
			Model:    "deepseek-ai/DeepSeek-R1", // Specify the Hugging Face model
			Messages: messages,
			ResponseFormat: &openai.ChatCompletionResponseFormat{
				Type: openai.ChatCompletionResponseFormatTypeJSONObject,
			},
			MaxTokens:   500,
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
	for _, tool := range toolsMap.Tools {
		siteContent, err := extractTextFromHTML(tool.URL)
		if err != nil {
			log.Printf("Error fetching %s: %v", tool.URL, err)
			continue
		}

		response, err := queryLLM(apiKey, baseURL, tool.Name, siteContent)
		result, err := extractJSONObject(response)

		println("result " + result)
	}
}

func extractJSONObject(text string) (string, error) {
	// Find the JSON object in the text
	startIndex := strings.Index(text, "{")
	if startIndex == -1 {
		return "", fmt.Errorf("no JSON object found in text")
	}

	endIndex := strings.Index(text, "}")
	if endIndex == -1 {
		return "", fmt.Errorf("no closing bracket found in JSON object")
	}

	// Extract the JSON object
	jsonObject := text[startIndex : endIndex+1]

	return jsonObject, nil
}

func main() {

	checkToolVersions()
}
