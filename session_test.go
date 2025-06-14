package tgbotapp_test

import (
	"testing"

	tgbotapp "github.com/StridersTech2025/go-telegram-bot-app"
)

func TestGetOrCreateSessionShouldCreateNewSessionIfNotExists(t *testing.T) {

	sessionManager := tgbotapp.NewInMemoryManager()

	var chatID int64 = 123

	s, err := sessionManager.GetOrCreateSession(chatID)

	if err != nil {
		t.Errorf(expectsNoError, err)
	}

	if s == nil {
		t.Errorf(expectsNotNil, "s")
	}

	if s.ChatID != chatID {
		t.Errorf("Expect chatID to be %d. Found chatID %d", chatID, s.ChatID)
	}

}

func TestSetSessionShouldUpdateSession(t *testing.T) {
	// Arrange
	mgr := tgbotapp.NewInMemoryManager()
	var chatID int64 = 123
	var testState tgbotapp.StateName = "TEST_STATE"
	var err error
	s, _ := mgr.GetOrCreateSession(chatID)

	s.State = testState

	// Act
	err = mgr.SetSession(chatID, s)

	// Assert
	if err != nil {
		t.Errorf(expectsNoError, err)
	}

	s, _ = mgr.GetOrCreateSession(chatID)

	if s == nil {
		t.Errorf(expectsNotNil, "s")
	}

	if s.State != testState {
		t.Errorf("Expected %s, found %s", testState, s.State)
	}

}
