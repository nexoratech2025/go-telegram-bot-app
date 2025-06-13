package tgbotapp

import (
	"fmt"
	"strings"
)

type CommandName string
type CallbackName string
type StateName string

const (
	CommandDelimiter = "@"

	ErrEmptyArgument = "Argument %s must not be empty."
)

type Router interface {
	GetCommandHandler(name CommandName) (HandlerFunc, bool)
	AddCommandHandler(name CommandName, handler HandlerFunc) error

	GetCallbackHandler(name CallbackName) (HandlerFunc, bool)
	AddCallbackHandler(name CallbackName, handler HandlerFunc) error

	GetMessageStateHandler(name StateName) (HandlerFunc, bool)
	AddMessageHandler(name StateName, handler HandlerFunc) error
}

func RouterMiddleware(router Router) Middleware {

	return func(context *BotContext, next HandlerFunc) {
		logger := context.Logger()

		var f HandlerFunc
		var ok bool
		switch {
		case context.Update.CallbackQuery != nil:
			var action CallbackName
			callbackData := context.Update.CallbackQuery.Data
			action, context.Params = extractCallback(callbackData)

			f, ok = router.GetCallbackHandler(action)

			if !ok {
				logger.WarnContext(context.Ctx, "No handler found for callback.", "callbackName", callbackData)
			}

			context.SetHandler(f)

		case context.Update.Message != nil && context.Update.Message.IsCommand():
			command := context.Update.Message.Command()
			f, ok = router.GetCommandHandler(CommandName(command))

			if !ok {
				logger.WarnContext(context.Ctx, "No handler found for command.", "commandName", command)
			}

			context.SetHandler(f)
			context.Params = strings.Split(context.Update.Message.CommandArguments(), CommandDelimiter)

		case context.Update.Message != nil:

			if context.Session != nil {
				f, ok = router.GetMessageStateHandler(context.Session.State)
				if !ok {
					logger.WarnContext(context.Ctx, "No handler found for message state.", "messageState", context.Session.State)
				}
			}

			context.SetHandler(f)
		}

		next(context)

	}

}

func extractCallback(callbackData string) (action CallbackName, args []string) {
	if len(callbackData) < 1 {
		return
	}

	s := strings.Split(callbackData, CommandDelimiter)

	action = CallbackName(s[0])

	if len(s) > 1 {
		args = s[1:]
	}

	return

}

// Default Implementation for Route Table
type RouteTable struct {
	commandRegistry      map[CommandName]HandlerFunc
	callbackRegistry     map[CallbackName]HandlerFunc
	messageStateRegistry map[StateName]HandlerFunc
}

func NewRouteTable() Router {
	return &RouteTable{
		commandRegistry:      make(map[CommandName]HandlerFunc),
		callbackRegistry:     make(map[CallbackName]HandlerFunc),
		messageStateRegistry: make(map[StateName]HandlerFunc),
	}
}

func (r *RouteTable) GetCommandHandler(name CommandName) (f HandlerFunc, ok bool) {

	f, ok = r.commandRegistry[name]
	return

}

func (r *RouteTable) GetCallbackHandler(name CallbackName) (f HandlerFunc, ok bool) {

	f, ok = r.callbackRegistry[name]
	return

}

func (r *RouteTable) GetMessageStateHandler(name StateName) (f HandlerFunc, ok bool) {

	f, ok = r.messageStateRegistry[name]
	return

}

func (r *RouteTable) AddCallbackHandler(name CallbackName, handler HandlerFunc) error {

	if len(name) < 1 {
		return fmt.Errorf(ErrEmptyArgument, "name")
	}

	if _, ok := r.callbackRegistry[name]; ok {
		return NewErrCallbackExists(name)
	}

	r.callbackRegistry[name] = handler

	return nil

}

func (r *RouteTable) AddCommandHandler(name CommandName, handler HandlerFunc) error {
	if len(name) < 1 {
		return fmt.Errorf(ErrEmptyArgument, "name")
	}

	if _, ok := r.commandRegistry[name]; ok {
		return NewErrCommandExists(name)
	}

	r.commandRegistry[name] = handler

	return nil

}

func (r *RouteTable) AddMessageHandler(name StateName, handler HandlerFunc) error {
	if len(name) < 1 {
		return fmt.Errorf(ErrEmptyArgument, "name")
	}
	if _, ok := r.messageStateRegistry[name]; ok {
		return NewErrMessageStateExists(name)
	}

	r.messageStateRegistry[name] = handler

	return nil

}
