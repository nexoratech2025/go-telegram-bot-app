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

// Return completely new application with no configuration.
func New(botAPI *tgbotapi.BotAPI, opts ...OptionFunc) *Application {
	app := &Application{
		middlewares: NewMiddlewareChain(),
		BotAPI:      botAPI,
	}
	return app.With(opts...)
}

func (a *Application) With(opts ...OptionFunc) *Application {
	for _, opt := range opts {
		opt(a)
	}

	return a
}

func defaultOptions(a *Application) {
	a.Router = NewRouteTable()
	a.SessionManager = NewInMemoryManager()
}

// Return new application with default configured Middlewares (Session and Router)
func Default(botAPI *tgbotapi.BotAPI, opts ...OptionFunc) *Application {

	options := []OptionFunc{defaultOptions}

	options = append(options, opts...)

	app := New(botAPI, options...)
	app.UseSession()
	app.UseRouting()
	return app
}

func (a *Application) RegisterCommand(name string, description string, handler HandlerFunc) error {

	botCommands = append(botCommands, tgbotapi.BotCommand{
		Command:     string(name),
		Description: description,
	})

	return a.Router.AddHandler(name, CommandHandler, handler)
}

func (a *Application) RegisterCallback(name string, handler HandlerFunc) error {
	return a.Router.AddHandler(name, CallbackHandler, handler)
}

func (a *Application) RegisterMessage(state string, handler HandlerFunc) error {
	return a.Router.AddHandler(state, MessageHandler, handler)
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

	go func() {
		for {
			select {
			case <-ctx.Done():
				err := a.shutdown()
				if err != nil {
					a.Logger.ErrorContext(ctx, "Error shutting down the application", "error_detail", err)
				}

			case update := <-updates:
				a.handleUpdate(ctx, &update)
			}
		}
	}()

	return nil

}

func (a *Application) shutdown() error {
	a.BotAPI.StopReceivingUpdates()
	return nil
}

func (a *Application) handleUpdate(ctx context.Context, update *tgbotapi.Update) {

	botCtx := NewBotContext(ctx, a, update)

	f := a.middlewares.Wrap(func(ctx *BotContext) {
		if a.handler != nil {
			a.handler(ctx)
		} else {
			a.Logger.ErrorContext(ctx.Ctx, "Error: Default handler should be set in routing middleware.")
		}
	})

	f(botCtx)

}

func (a *Application) initBotCommands() error {

	cmds := tgbotapi.NewSetMyCommands(botCommands...)

	_, err := a.BotAPI.Send(cmds)

	return err

}
