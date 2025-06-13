package tgbotapp

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidBotToken = errors.New("Invalid bot token.")
)

type ErrCommandExists struct {
	commandName CommandName
}

func NewErrCommandExists(commandName CommandName) error {
	return &ErrCommandExists{
		commandName: commandName,
	}
}

func (e *ErrCommandExists) Error() string {
	return fmt.Sprintf("Command handler already exists for name: %s", e.commandName)
}

type ErrCallbackExists struct {
	callbackName CallbackName
}

func NewErrCallbackExists(callbackName CallbackName) error {
	return &ErrCallbackExists{
		callbackName: callbackName,
	}
}

func (e *ErrCallbackExists) Error() string {
	return fmt.Sprintf("Callback handler already exists for name: %s", e.callbackName)
}

type ErrMessageStateExists struct {
	messageState StateName
}

func NewErrMessageStateExists(messageState StateName) error {
	return &ErrMessageStateExists{
		messageState: messageState,
	}
}

func (e *ErrMessageStateExists) Error() string {
	return fmt.Sprintf("Message state handler already exists for state: %s", e.messageState)
}
