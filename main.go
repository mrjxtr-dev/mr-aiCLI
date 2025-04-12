package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/mrjxtr-dev/mr-aiCLI/config"
	"github.com/mrjxtr-dev/mr-aiCLI/custom_errors"
)

func main() {
	// Initialize client with configuration
	client := config.LoadConfig()

	// Set available models for auto-routing
	client.SetAvailableModels(config.AvailableModels)

	fmt.Println("------------------------------")
	fmt.Println("          Mr-AI Chat          ")
	fmt.Println("------------------------------")

	// Set up input scanner
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Printf("\nYou: ")
		if !scanner.Scan() {
			break
		}
		userInput := scanner.Text()

		// Handle exit commands
		if userInput == "exit" || userInput == "quit" {
			fmt.Println("Exiting chat.")
			break
		}

		err := client.SendMessage(userInput)
		if err != nil {
			custom_errors.HandleError(err)
			continue
		}

		// Display the last assistant response
		for i := len(client.Messages) - 1; i >= 0; i-- {
			if client.Messages[i].Role == "assistant" {
				fmt.Printf("\nMr-AI: %s", client.Messages[i].Content)
				break
			}
		}
	}
}
