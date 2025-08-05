package tgbotapp_test

import (
	"testing"

	tgbotapp "github.com/nexoratech2025/go-telegram-bot-app"
)

func TestGetOrCreateSessionShouldCreateNewSessionIfNotExists(t *testing.T) {

	sessionManager := tgbotapp.NewDefaultInMemoryManager()

	var chatID int64 = 123

	s, err := sessionManager.GetOrCreate(chatID)

	if err != nil {
		t.Errorf(expectsNoError, err)
	}

	if s == nil {
		t.Errorf(expectsNotNil, "s")
	}

}

func TestSetSessionShouldUpdateSession(t *testing.T) {
	// Arrange
	mgr := tgbotapp.NewDefaultInMemoryManager()
	const chatID int64 = 123
	testData := 420
	testKey := "test"
	var err error
	s, _ := mgr.GetOrCreate(chatID)

	s.Set(testKey, testData)

	// Act
	err = mgr.Set(chatID, s)

	// Assert
	if err != nil {
		t.Errorf(expectsNoError, err)
	}

	s, _ = mgr.GetOrCreate(chatID)

	if s == nil {
		t.Errorf(expectsNotNil, "s")
	}

	state, ok := s.Get(testKey)
	if !ok {
		t.Errorf("Expected session to contain %q key", "state")
	}

	switch d := state.(type) {
	case int:
		if d != testData {
			t.Errorf("Expected %d. Received: %d", testData, d)
		}
	default:
		t.Errorf("Cannot convert session data back to int")
	}

}

func TestSetStateShouldUpdateSessionState(t *testing.T) {
	// Arrange
	mgr := tgbotapp.NewDefaultInMemoryManager()
	const chatID int64 = 123
	var testState = "TEST_STATE"
	var err error
	s, _ := mgr.GetOrCreate(chatID)

	s.SetState(testState)

	// Act
	err = mgr.Set(chatID, s)

	// Assert
	if err != nil {
		t.Errorf(expectsNoError, err)
	}

	s, _ = mgr.GetOrCreate(chatID)

	if s == nil {
		t.Errorf(expectsNotNil, "s")
	}

	state := s.CurrentState()

	if state != testState {
		t.Errorf("Expected state to be %q, found state %q", testState, state)
	}

}
