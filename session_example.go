package tgbotapp

import (
	"fmt"
	"log"
)

// ExampleBotHandler shows how to use session clearing in a bot handler
func ExampleBotHandler(ctx *BotContext) {
	// Clear all session data when user wants to start fresh
	if ctx.Session != nil {
		// Option 1: Use the ClearData method (recommended)
		ctx.Session.ClearData()

		// Option 2: Clear specific keys
		ctx.Session.Delete("user_preferences")
		ctx.Session.Delete("temp_data")

		// Option 3: Clear state
		ctx.Session.SetState("")

		log.Println("Session data cleared successfully")
	}
}

// ExampleUsage demonstrates how to use session clearing
func ExampleUsage() {
	// Create a session manager
	sessionManager := NewDefaultInMemoryManager()

	// Get or create a session for a user
	userSession, err := sessionManager.GetOrCreate(12345)
	if err != nil {
		log.Fatal(err)
	}

	// Add some data to the session
	userSession.Set("user_id", 12345)
	userSession.Set("username", "john_doe")
	userSession.Set("temp_data", "some temporary information")
	userSession.SetState("waiting_for_input")

	fmt.Println("Before clearing:")
	fmt.Printf("State: %s\n", userSession.CurrentState())
	fmt.Printf("Keys: %v\n", userSession.GetAllKeys())

	// Clear all session data
	userSession.ClearData()

	fmt.Println("\nAfter clearing data:")
	fmt.Printf("State: %s\n", userSession.CurrentState()) // State is preserved
	fmt.Printf("Keys: %v\n", userSession.GetAllKeys())

	// Clear state separately
	userSession.SetState("")

	fmt.Println("\nAfter clearing state:")
	fmt.Printf("State: %s\n", userSession.CurrentState())
}
