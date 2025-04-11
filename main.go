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

	fmt.Println("Mr-AI Chat")
	fmt.Println("---------------------")
	fmt.Println("- Type 'exit' or 'quit' to exit")

	// NOTE: commented out print statements bellow to simplify outputs
	// fmt.Println("- Type 'stream:' before your message to use streaming mode")
	// Print out models used
	// fmt.Println("\nAvailable models (in priority order):")
	// for i, model := range config.AvailableModels {
	// 	fmt.Printf("%d. %s\n", i+1, model)
	// }

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

		// Check if streaming mode is requested
		if len(userInput) > 7 && userInput[:7] == "stream:" {
			// Handle streaming mode
			message := userInput[7:] // Extract the actual message content after "stream:"
			fmt.Printf("\nMr-AI: ")
			err := client.SendMessageStream(message, func(chunk string) {
				fmt.Print(chunk) // Print each token as it arrives
			})
			if err != nil {
				custom_errors.HandleError(err)
				continue
			}
		} else {
			// Handle regular (non-streaming) mode
			err := client.SendMessage(userInput)
			if err != nil {
				custom_errors.HandleError(err)
				continue
			}

			// Display the response
			if len(client.Messages) > 0 {
				lastMsgIndex := len(client.Messages) - 1
				lastMsg := client.Messages[lastMsgIndex]

				if lastMsg.Role == "assistant" {
					fmt.Printf("\nMr-AI: %s\n", lastMsg.Content)
				} else {
					// Try to find the last assistant message
					for i := lastMsgIndex - 1; i >= 0; i-- {
						if client.Messages[i].Role == "assistant" {
							fmt.Printf("\nMr-AI: %s\n", client.Messages[i].Content)
							break
						}
					}
				}
			}
		}
	}
}
