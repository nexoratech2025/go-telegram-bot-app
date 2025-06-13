package tgbotapp

import (
	"log/slog"

	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	botCommands []tgbotapi.BotCommand
)

// Control the application option.
type OptionFunc func(*Application)

type Application struct {
	middlewares *MiddlewareChain
	handler     HandlerFunc

	SessionManager SessionManager
	Logger         *slog.Logger
	Router         Router
	BotAPI         *tgbotapi.BotAPI
}

func New(botAPI tgbotapi.BotAPI, opts ...OptionFunc) *Application {
	app := &Application{
		middlewares: NewMiddlewareChain(),
	}
	return app.With(opts...)
}

func (a *Application) With(opts ...OptionFunc) *Application {
	for _, opt := range opts {
		opt(a)
	}

	return a
}

func Default(botAPI tgbotapi.BotAPI, opts ...OptionFunc) *Application {

	app := New(botAPI, opts...)
	app.SessionManager = NewInMemoryManager()
	app.Router = NewRouteTable()
	app.UseSession()
	app.UseRouting()
	return app
}

func (a *Application) RegisterCommand(name CommandName, description string, handler HandlerFunc) error {

	botCommands = append(botCommands, tgbotapi.BotCommand{
		Command:     string(name),
		Description: description,
	})

	return a.Router.AddCommandHandler(name, handler)
}

func (a *Application) RegisterCallback(name CallbackName, handler HandlerFunc) error {
	return a.Router.AddCallbackHandler(name, handler)
}

func (a *Application) RegisterMessage(state StateName, handler HandlerFunc) error {
	return a.Router.AddMessageHandler(state, handler)
}

func (a *Application) Use(middlewares ...Middleware) {
	a.middlewares.Append(middlewares...)
}

func (a *Application) UseRouting() {

	a.middlewares.Append(RouterMiddleware(a.Router))

}

func (a *Application) UseSession() {

	a.middlewares.Append(SessionMiddleware(a.SessionManager))

}

// TODO: Try to make update handling concurrent without problems

func (a *Application) Start(ctx context.Context) error {

	err := a.initBotCommands()
	if err != nil {
		a.Logger.ErrorContext(ctx, "Cannot set commands list.", "error_detail", err)
	}

	updateCfg := tgbotapi.NewUpdate(0)
	updateCfg.Timeout = 60

	updates := a.BotAPI.GetUpdatesChan(updateCfg)

	for {
		select {
		case <-ctx.Done():
			return a.shutdown()

		case update := <-updates:
			a.handleUpdate(ctx, &update)
		}
	}

}

func (a *Application) shutdown() error {
	a.BotAPI.StopReceivingUpdates()
	return nil
}

func (a *Application) handleUpdate(ctx context.Context, update *tgbotapi.Update) {

	botCtx := NewBotContext(ctx, a.BotAPI, update)

	f := a.middlewares.Wrap(func(ctx *BotContext) {
		if ctx.app.handler != nil {
			ctx.app.handler(botCtx)
		} else {
			ctx.Logger().WarnContext(ctx.Ctx, "No Handler Found. Returning nothing.")
		}
	})

	f(botCtx)

}

func (a *Application) initBotCommands() error {

	cmds := tgbotapi.NewSetMyCommands(botCommands...)

	_, err := a.BotAPI.Send(cmds)

	return err

}
