package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/signal"
	"slices"
	"syscall"

	"github.com/mrjxtr-dev/mr-aiCLI/config"
	"github.com/mrjxtr-dev/mr-aiCLI/custom_errors"
)

func main() {
	// Set up channel to listen for interrupt signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		fmt.Printf("\n\nCTRL+C signal received, exiting!")
		os.Exit(0)
	}()

	// Initialize client with configuration
	client, err := config.LoadClient()
	if err != nil {
		log.Fatal(err)
	}

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

		// Check for exit commands
		exitCmds := []string{"exit", "quit", "bye", "goodbye", "q"}
		if slices.Contains(exitCmds, userInput) {
			os.Exit(0)
		}

		// Send user input to OpenRouter
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
