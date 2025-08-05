package tgbotapp_test

import (
	"errors"
	"testing"

	tgbotapp "github.com/nexoratech2025/go-telegram-bot-app"
)

const (
	testKey  = "test"
	testData = "data123"
)

func dummyHandler(*tgbotapp.BotContext) {}
func dummyHandlerWithAction(ctx *tgbotapp.BotContext) {
	ctx.SetData(testKey, testData)
}

func TestAddHandlerShouldReturnErrorForDuplicateName(t *testing.T) {
	// Arrange
	router := tgbotapp.NewRouteTable()

	var commandName = "test_command"
	err := router.AddHandler(commandName, tgbotapp.CommandHandler, dummyHandler)

	if err != nil {
		t.Errorf(expectsNoError, err)
	}

	// Act

	err = router.AddHandler(commandName, tgbotapp.CommandHandler, dummyHandler)

	// Assert
	if err == nil {
		t.Error(expectsError)
	}

	var expected *tgbotapp.ErrHandlerAlreadyExists

	if !errors.As(err, &expected) {
		t.Errorf(expectsErrorType, expected, err)
	}

}

func TestAddHandlerShouldReturnOkForSameNameWithDifferentType(t *testing.T) {
	// Arrange
	router := tgbotapp.NewRouteTable()

	var commandName = "test_command"
	err := router.AddHandler(commandName, tgbotapp.CommandHandler, dummyHandler)

	if err != nil {
		t.Errorf(expectsNoError, err)
	}

	// Act

	err = router.AddHandler(commandName, tgbotapp.CallbackHandler, dummyHandler)

	// Assert
	if err != nil {
		t.Errorf(expectsNoError, err)
	}

}

func TestGetHandlerShouldReturnSameFunction(t *testing.T) {
	// Arrange
	router := tgbotapp.NewRouteTable()

	var commandName = "test_command"
	err := router.AddHandler(commandName, tgbotapp.CommandHandler, dummyHandlerWithAction)

	if err != nil {
		t.Errorf(expectsNoError, err)
	}
	ctx := tgbotapp.NewBotContext(t.Context(), nil, nil)

	// Act

	h, ok := router.GetHandler(commandName, tgbotapp.CommandHandler)

	// Assert
	if !ok {
		t.Errorf("Expects to return handler info struct. Found none.")
	}

	h.Func(ctx)

	d, ok := ctx.GetData(testKey)

	if !ok {
		t.Errorf("Expects context to have data with key %s", testKey)

	}

	switch v := d.(type) {
	case string:
		if v != testData {
			t.Errorf("Expects context to have value %s, found %s", testData, v)
		}
	default:
		t.Errorf("Expects context data to have string value.")
	}

}
