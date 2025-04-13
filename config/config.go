package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/mrjxtr-dev/mr-aiCLI/client"
)

// Available models with their priority order
var availableModels = []string{
	"google/gemini-2.0-flash-exp:free",
	"openrouter/optimus-alpha",
	"meta-llama/llama-4-scout:free",
	"nvidia/llama-3.1-nemotron-ultra-253b-v1:free",
}

// LoadConfig loads environment variables and initializes an OpenRouter client
func LoadClient() (*client.OpenRouterClient, error) {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("ERROR: loading .env file")
	}

	// Get configuration from environment variables
	apiKey := os.Getenv("OPENROUTER_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("OPENROUTER_API_KEY not found in .env")
	}
	baseURL := os.Getenv("OPENROUTER_BASE_URL")
	if baseURL == "" {
		return nil, fmt.Errorf("OPENROUTER_BASE_URL not found in .env")
	}

	// Set up client prompts
	systemPrompt := "You are a helpful personal AI assistant named Mr-AI."
	customContext := `Always use casual language in your response.
	Your responses should be short, concise and to the point.
	Do your best to sound as human as possible.`

	// Create and initialize the client with the highest priority model
	client := client.New(
		apiKey,
		baseURL,
		availableModels[0],
		systemPrompt,
		customContext,
	)

	// Set available models for auto-routing
	client.SetAvailableModels(availableModels)

	// Set up initial system messages
	client.InitContext()

	return client, nil
}
