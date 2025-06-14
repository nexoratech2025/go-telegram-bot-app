package tgbotapp

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidBotToken = errors.New("Invalid bot token.")
)

type ErrHandlerAlreadyExists struct {
	name   string
	action HandlerAction
}

func (e *ErrHandlerAlreadyExists) Error() string {
	return fmt.Sprintf("Handler function already exists for name %s of type %s", e.name, e.action)
}

func NewErrHandlerAlreadyExists(name string, action HandlerAction) error {
	return &ErrHandlerAlreadyExists{
		name:   name,
		action: action,
	}
}

type ErrInvalidArgument struct {
	reason  string
	argName string
}

func (e *ErrInvalidArgument) Error() string {
	return fmt.Sprintf("Invalid Argument %q: %s", e.argName, e.reason)
}

func NewErrInvalidArgument(reason string, argName string) error {
	return &ErrInvalidArgument{
		reason,
		argName,
	}
}
