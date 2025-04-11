package custom_errors

import (
	"fmt"
	"strings"
)

// Helper function to display errors with special handling for rate limits
func HandleError(err error) {
	errStr := err.Error()

	// Check if it's a rate limit error with no more models to try
	if strings.Contains(errStr, "All available models have been tried") {
		fmt.Println(
			"All models have reached their rate limits. Please try again later.",
		)
	} else if strings.Contains(errStr, "RATE LIMIT ERROR") {
		// This should now be handled by auto-routing
		fmt.Println("Rate limit reached, attempting to switch models...")
	} else {
		// Regular error display
		fmt.Printf("Error: %v\n", err)
	}
}
