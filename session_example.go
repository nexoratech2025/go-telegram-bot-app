package tgbotapp

import (
	"fmt"
	"log"
)

// ExampleSession demonstrates how to use session clearing in other applications
type ExampleSession struct {
	data  map[string]any
	state string
}

func NewExampleSession() *ExampleSession {
	return &ExampleSession{
		data: make(map[string]any),
	}
}

// Implement the required interface for ClearSessionData
func (s *ExampleSession) GetAllKeys() []string {
	keys := make([]string, 0, len(s.data))
	for k := range s.data {
		keys = append(keys, k)
	}
	return keys
}

func (s *ExampleSession) Delete(key string) {
	delete(s.data, key)
}

// Additional methods for the example
func (s *ExampleSession) Set(key string, value any) {
	s.data[key] = value
}

func (s *ExampleSession) Get(key string) (any, bool) {
	value, ok := s.data[key]
	return value, ok
}

func (s *ExampleSession) SetState(state string) {
	s.state = state
}

func (s *ExampleSession) GetState() string {
	return s.state
}

// ExampleUsage demonstrates how to use ClearSessionData in other applications
func ExampleUsage() {
	// Create a session
	session := NewExampleSession()

	// Add some data
	session.Set("user_id", 12345)
	session.Set("username", "john_doe")
	session.Set("temp_data", "some temporary information")
	session.SetState("waiting_for_input")

	fmt.Println("Before clearing:")
	fmt.Printf("State: %s\n", session.GetState())
	fmt.Printf("Keys: %v\n", session.GetAllKeys())

	// Clear all session data using the exported function
	ClearSessionData(session)

	fmt.Println("\nAfter clearing data:")
	fmt.Printf("State: %s\n", session.GetState()) // State is preserved
	fmt.Printf("Keys: %v\n", session.GetAllKeys())

	// Clear state separately
	session.SetState("")

	fmt.Println("\nAfter clearing state:")
	fmt.Printf("State: %s\n", session.GetState())
}

// ExampleBotHandler shows how to use it in a bot handler
func ExampleBotHandler(ctx *BotContext) {
	// Clear all session data when user wants to start fresh
	if ctx.Session != nil {
		// Option 1: Use the generic function
		ClearSessionData(ctx.Session)

		// Option 2: Use the method if available
		if clearable, ok := ctx.Session.(interface{ ClearData() }); ok {
			clearable.ClearData()
		}

		// Option 3: Clear specific keys
		ctx.Session.Delete("user_preferences")
		ctx.Session.Delete("temp_data")

		log.Println("Session data cleared successfully")
	}
}
