package tgbotapp

import (
	"errors"
	"log/slog"
	"sync"

	"context"

	"github.com/StridersTech2025/go-telegram-bot-app/session"
	"github.com/StridersTech2025/go-telegram-bot-app/util"
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
	wg          sync.WaitGroup

	SessionManager session.SessionManager[int64]
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
	a.Logger = slog.Default()
	a.Router = NewRouteTable()
	a.SessionManager = NewDefaultInMemoryManager()
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
		Command:     name,
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

func (a *Application) RegisterDocument(handler HandlerFunc) error {
	return a.Router.AddHandler("document", DocumentHandler, handler)
}

func (a *Application) RegisterDocumentByType(docType string, handler HandlerFunc) error {
	return a.Router.AddHandler(docType, DocumentHandler, handler)
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

	a.Logger.InfoContext(ctx, "Starting application...")

	err := a.initBotCommands()
	if err != nil {
		a.Logger.ErrorContext(ctx, "Cannot set commands list.", "error_detail", err)
	} else {
		a.Logger.InfoContext(ctx, "Command list set successfully.")
	}

	updateCfg := tgbotapi.NewUpdate(0)
	updateCfg.Timeout = 60

	updates := a.BotAPI.GetUpdatesChan(updateCfg)

	// Poll loop
	a.wg.Add(1)
	go func() {
		defer a.wg.Done()
		a.Logger.Info("Listening for updates from bot.", "bot_id", a.BotAPI.Self.ID, "bot_username", a.BotAPI.Self.UserName)
		for {
			select {
			case <-ctx.Done():
				a.shutdown()
				return

			case update := <-updates:
				a.handleUpdate(ctx, &update)
			}
		}
	}()

	<-ctx.Done()

	a.wg.Wait()
	return nil

}

func (a *Application) shutdown() {
	a.Logger.Info("Shutting Down the application...")
	a.BotAPI.StopReceivingUpdates()
	a.Logger.Info("Application stopped successfully.")

}

func (a *Application) handleUpdate(ctx context.Context, update *tgbotapi.Update) {

	botCtx := NewBotContext(ctx, a, update)

	f := a.middlewares.Wrap(func(ctx *BotContext) {
		a.Logger.InfoContext(ctx.Ctx, "Processing update.", "from", ctx.Update.SentFrom().ID)

		if a.handler != nil {
			a.handler(ctx)
		} else {
			a.Logger.ErrorContext(ctx.Ctx, "Error: Default handler should be set in routing middleware.")
		}
	})

	f(botCtx)

}

func (a *Application) initBotCommands() error {

	if len(botCommands) < 1 {
		a.Logger.Warn("No bot commands found.")
		return nil
	}

	cmds := tgbotapi.NewSetMyCommands(botCommands...)

	// tgbotapi Send method handles message response only.
	// setMyCommands method return boolean.
	// Thus custom setMyCommand function is used here.

	ok, err := util.SendSetMyCommands(*a.BotAPI, cmds)

	if err != nil {
		return err
	}

	if !ok {
		return errors.New("Cannot set command.")
	}
	return nil

}
