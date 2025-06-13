package tgbotapp_test

import (
	"errors"
	"testing"

	tgbotapp "github.com/StridersTech2025/go-telegram-bot-app/v1"
)

const (
	expectNoError   = "Should not return error. Got error: %#v"
	expectError     = "Should return error. got no error"
	expectErrorType = "Should return error type %#v. Got error type %#v"
)

func dummyHandler(*tgbotapp.BotContext) {}

func TestAddCommandHandlerShouldReturnErrorForDuplicateName(t *testing.T) {
	router := tgbotapp.NewRouteTable()

	var commandName tgbotapp.CommandName = "test_command"
	err := router.AddCommandHandler(commandName, dummyHandler)

	if err != nil {
		t.Errorf(expectNoError, err)
	}

	err = router.AddCommandHandler(commandName, dummyHandler)

	if err == nil {
		t.Error(expectError)
	}

	var expected *tgbotapp.ErrCommandExists

	if !errors.As(err, &expected) {
		t.Errorf(expectErrorType, expected, err)
	}

}

func TestAddCallbackHandlerShouldReturnErrorForDuplicateName(t *testing.T) {
	router := tgbotapp.NewRouteTable()

	var callbackName tgbotapp.CallbackName = "test_callback"
	err := router.AddCallbackHandler(callbackName, dummyHandler)

	if err != nil {
		t.Errorf(expectNoError, err)
	}

	err = router.AddCallbackHandler(callbackName, dummyHandler)

	if err == nil {
		t.Error(expectError)
	}

	var expected *tgbotapp.ErrCallbackExists

	if !errors.As(err, &expected) {
		t.Errorf(expectErrorType, expected, err)
	}

}

func TestAddMessageHandlerShouldReturnErrorForDuplicateName(t *testing.T) {
	router := tgbotapp.NewRouteTable()

	var stateName tgbotapp.StateName = "test_state"
	err := router.AddMessageHandler(stateName, dummyHandler)

	if err != nil {
		t.Errorf(expectNoError, err)
	}

	err = router.AddMessageHandler(stateName, dummyHandler)

	if err == nil {
		t.Error(expectError)
	}

	var expected *tgbotapp.ErrMessageStateExists

	if !errors.As(err, &expected) {
		t.Errorf(expectErrorType, expected, err)
	}

}
