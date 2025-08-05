package tgbotapp_test

import (
	"os"
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "github.com/joho/godotenv/autoload"
	tgbotapp "github.com/nexoratech2025/go-telegram-bot-app"
)

const (
	botTokenEnvKey = "BOT_TOKEN"

	expectsNoError   = "Should not return error. Got error: %#v"
	expectsError     = "Should return error. got no error"
	expectsErrorType = "Should return error type %#v. Got error type %#v"
	expectsNotNil    = "Expects %s to not nil."
)

func TestApplicationShouldStartCorrectly(t *testing.T) {
	// Arrange

	token := os.Getenv(botTokenEnvKey)
	botAPI, err := tgbotapi.NewBotAPI(token)

	if err != nil {
		t.Errorf("Expected botAPI to initialised. Received Error: %#v", err)
	}

	app := tgbotapp.Default(botAPI)

	// Act

	err = app.Start(t.Context())

	//Assert

	if err != nil {
		t.Errorf("Expected app to start smoothly. found error: %#v", err)
	}

}

func TestApplicationShouldRegisterCommandCorrectly(t *testing.T) {
	// Arrange

	token := os.Getenv(botTokenEnvKey)
	botAPI, err := tgbotapi.NewBotAPI(token)

	if err != nil {
		t.Errorf("Expected botAPI to initialised. Received Error: %#v", err)
	}

	app := tgbotapp.Default(botAPI)

	app.RegisterCommand("ping", "pong", func(bc *tgbotapp.BotContext) {})

	// Act

	err = app.Start(t.Context())

	//Assert

	if err != nil {
		t.Errorf("Expected app to start smoothly. found error: %#v", err)
	}

}
