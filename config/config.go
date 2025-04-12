package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/mrjxtr-dev/mr-aiCLI/client"
)

// Available models with their priority order
var AvailableModels = []string{
	"google/gemini-2.0-flash-exp:free",
	"openrouter/optimus-alpha",
	"meta-llama/llama-4-scout:free",
	"nvidia/llama-3.1-nemotron-ultra-253b-v1:free",
}

// LoadConfig loads environment variables and initializes an OpenRouter client
func LoadConfig() *client.OpenRouterClient {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Get configuration from environment variables
	apiKey := os.Getenv("OPENROUTER_API_KEY")
	baseURL := os.Getenv("OPENROUTER_BASE_URL")

	// Use default URL if not set
	if baseURL == "" {
		baseURL = "https://openrouter.ai/api/v1/chat/completions"
		log.Println("No OPENROUTER_BASE_URL found in .env, using default:", baseURL)
	}

	// Set up client prompts
	systemPrompt := "You are a helpful personal AI assistant named Mr-AI. Respond in a concise manner and use casual language."
	customContext := ""

	// Create and initialize the client with the highest priority model
	client := client.NewOpenRouterClient(
		apiKey,
		baseURL,
		AvailableModels[0],
		systemPrompt,
		customContext,
	)

	// Set up initial system messages
	client.InitContext()

	return client
}
