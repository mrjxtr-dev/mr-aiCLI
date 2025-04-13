package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Message represents a single message in a conversation
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatRequest defines the structure for API requests to OpenRouter
type ChatRequest struct {
	Model     string    `json:"model"`
	Messages  []Message `json:"messages"`
	MaxTokens int       `json:"max_tokens"`
}

// OpenRouterClient handles communication with the OpenRouter API
type OpenRouterClient struct {
	APIKey          string
	BaseURL         string
	Referer         string
	Title           string
	Model           string
	Client          *http.Client
	SystemPrompt    string
	CustomContext   string
	Messages        []Message
	ModelIndex      int      // Tracks current model index for failover
	AvailableModels []string // List of models to try in order
}

// NewOpenRouterClient creates a new client with the provided configuration
func New(
	apiKey, baseURL, model, systemPrompt, customContext string,
) *OpenRouterClient {
	return &OpenRouterClient{
		APIKey:        apiKey,
		BaseURL:       baseURL,
		Model:         model,
		Client:        &http.Client{},
		SystemPrompt:  systemPrompt,
		CustomContext: customContext,
		Messages:      []Message{},
		ModelIndex:    0,
	}
}

// SendMessage sends a user message to the API and waits for full response
func (orc *OpenRouterClient) SendMessage(content string) error {
	// Add the user's new message to the conversation history
	orc.Messages = append(orc.Messages, Message{
		Role:    "user",
		Content: content,
	})

	// Keep trying models until one succeeds or we run out of options
	var lastError error
	for {
		err := orc.sendMessageWithCurrentModel()
		if err == nil {
			// Success!
			return nil
		}

		lastError = err

		// Check if this is a rate limit error
		if strings.Contains(err.Error(), "RATE LIMIT ERROR") {
			// Try to switch to next model
			if canSwitch, msg := orc.TryNextModel(); canSwitch {
				continue // Try again with the new model
			} else {
				// No more models to try
				return fmt.Errorf("%s", msg)
			}
		}

		// Not a rate limit error
		return lastError
	}
}

// sendMessageWithCurrentModel sends a message using the current model
func (orc *OpenRouterClient) sendMessageWithCurrentModel() error {
	// Prepare request body
	body := ChatRequest{
		Model:     orc.Model,
		Messages:  orc.Messages,
		MaxTokens: 4000,
	}

	jsonData, err := json.Marshal(body)
	if err != nil {
		return err
	}

	// Create and send request
	req, err := http.NewRequest("POST", orc.BaseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	orc.setHeaders(req)
	resp, err := orc.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close() // Ensure response is always closed

	// Read the response body once so we can inspect it
	bodyBytes, _ := io.ReadAll(resp.Body)
	bodyStr := string(bodyBytes)

	// Check for rate limit error indicators
	if resp.StatusCode == 429 || strings.Contains(bodyStr, "Rate limit exceeded") ||
		strings.Contains(
			bodyStr,
			"rate limit",
		) || strings.Contains(bodyStr, "ratelimit") {
		return fmt.Errorf("RATE LIMIT ERROR: %s", bodyStr)
	}

	// Handle other non-200 responses
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf(
			"API error (status %d): %s",
			resp.StatusCode,
			bodyStr,
		)
	}

	// Create a new reader from the captured body
	bodyReader := io.NopCloser(bytes.NewReader(bodyBytes))

	// Create a new response with the same headers but with the new reader
	newResp := &http.Response{
		Status:        resp.Status,
		StatusCode:    resp.StatusCode,
		Header:        resp.Header,
		Body:          bodyReader,
		ContentLength: int64(len(bodyBytes)),
	}

	// Parse response and update conversation history
	return orc.parseResponse(newResp)
}

// parseResponse processes the API response and adds to message history
func (orc *OpenRouterClient) parseResponse(resp *http.Response) error {
	// Parse JSON response - handle different possible response formats
	var standardResult struct {
		Choices []struct {
			Message Message `json:"message"`
		} `json:"choices"`
		Error struct {
			Message string `json:"message"`
			Code    int    `json:"code"`
		} `json:"error"`
	}

	var alternateResult struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
			Index   int    `json:"index"`
			Role    string `json:"role,omitempty"`
			Content string `json:"content,omitempty"`
		} `json:"choices"`
		Error struct {
			Message string `json:"message"`
			Code    int    `json:"code"`
		} `json:"error"`
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %v", err)
	}

	// Check for rate limit error in the raw response
	if strings.Contains(string(bodyBytes), "Rate limit exceeded") {
		return fmt.Errorf("RATE LIMIT ERROR: %s", string(bodyBytes))
	}

	// Try to parse as standard OpenAI-compatible format
	err1 := json.Unmarshal(bodyBytes, &standardResult)

	// Check for error field first
	if err1 == nil && standardResult.Error.Message != "" {
		if standardResult.Error.Code == 429 ||
			strings.Contains(standardResult.Error.Message, "Rate limit") {
			return fmt.Errorf("RATE LIMIT ERROR: %s", standardResult.Error.Message)
		}
		return fmt.Errorf("API Error: %s", standardResult.Error.Message)
	}

	// Check for valid choices
	if err1 == nil && len(standardResult.Choices) > 0 &&
		standardResult.Choices[0].Message.Role != "" {
		// Add assistant's response to conversation history
		orc.Messages = append(orc.Messages, standardResult.Choices[0].Message)
		return nil
	}

	// Try alternate format
	err2 := json.Unmarshal(bodyBytes, &alternateResult)

	// Check for error field
	if err2 == nil && alternateResult.Error.Message != "" {
		if alternateResult.Error.Code == 429 ||
			strings.Contains(alternateResult.Error.Message, "Rate limit") {
			return fmt.Errorf("RATE LIMIT ERROR: %s", alternateResult.Error.Message)
		}
		return fmt.Errorf("API Error: %s", alternateResult.Error.Message)
	}

	// Check for valid choices
	if err2 == nil && len(alternateResult.Choices) > 0 {
		// Construct a message from the alternate format
		var role string
		var content string

		if alternateResult.Choices[0].Role != "" {
			// Some models use role directly
			role = alternateResult.Choices[0].Role
			content = alternateResult.Choices[0].Content
		} else {
			// Others use the message.content field
			role = "assistant" // Default role
			content = alternateResult.Choices[0].Message.Content
		}

		message := Message{
			Role:    role,
			Content: content,
		}

		// Add to history
		orc.Messages = append(orc.Messages, message)
		return nil
	}

	// If both parsing methods failed
	if err1 != nil && err2 != nil {
		return fmt.Errorf("failed to parse response: %v, %v", err1, err2)
	}

	// Check one last time for rate limit error by keyword before failing
	if strings.Contains(string(bodyBytes), "Rate limit") {
		return fmt.Errorf("RATE LIMIT ERROR: %s", string(bodyBytes))
	}

	// If response was parsed but had no choices
	return fmt.Errorf("API returned empty choices array")
}

// setHeaders adds required and optional headers to the request
func (orc *OpenRouterClient) setHeaders(req *http.Request) {
	req.Header.Set("Authorization", "Bearer "+orc.APIKey)
	req.Header.Set("Content-Type", "application/json")

	// Add optional headers if provided
	if orc.Referer != "" {
		req.Header.Set("HTTP-Referer", orc.Referer)
	}
	if orc.Title != "" {
		req.Header.Set("X-Title", orc.Title)
	}
}

// InitContext sets up initial system messages based on provided prompts
func (orc *OpenRouterClient) InitContext() {
	// Add system prompt if provided
	if orc.SystemPrompt != "" {
		orc.Messages = append(orc.Messages, Message{
			Role:    "system",
			Content: orc.SystemPrompt,
		})
	}

	// Add custom context if provided
	if orc.CustomContext != "" {
		orc.Messages = append(orc.Messages, Message{
			Role:    "system",
			Content: "Context: " + orc.CustomContext,
		})
	}
}

// TryNextModel attempts to use the next available model in the list
func (orc *OpenRouterClient) TryNextModel() (bool, string) {
	// If models list not provided or at the end, can't switch

	if len(orc.AvailableModels) == 0 {
		// No models list provided, using only the default model
		return false, ""
	}

	orc.ModelIndex++
	if orc.ModelIndex >= len(orc.AvailableModels) {
		// We've tried all models
		return false, "All available models have been tried and reached rate limits"
	}

	// Switch to next model
	orc.Model = orc.AvailableModels[orc.ModelIndex]
	return true, orc.Model
}

// SetAvailableModels sets the list of models to try in order
func (orc *OpenRouterClient) SetAvailableModels(models []string) {
	orc.AvailableModels = models
	// Initialize current model to the first one if needed
	if orc.ModelIndex == 0 && len(models) > 0 {
		orc.Model = orc.AvailableModels[0]
	}
}
