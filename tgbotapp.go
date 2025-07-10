package tgbotapp

import (
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"sync"

	"context"

	"github.com/StridersTech2025/go-telegram-bot-app/botwrapper"
	"github.com/StridersTech2025/go-telegram-bot-app/helper"
	"github.com/StridersTech2025/go-telegram-bot-app/session"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	botCommands []tgbotapi.BotCommand
)

// Control the application option.
type OptionFunc func(*Application)
type LifeCycleHook func(*Application) error

type Application struct {
	middlewares *MiddlewareChain
	handler     HandlerFunc
	wg          sync.WaitGroup

	before []LifeCycleHook
	after  []LifeCycleHook

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

	err := a.botSetup()
	if err != nil {
		a.Logger.ErrorContext(ctx, "Cannot set commands list.", "error_detail", err)
	} else {
		a.Logger.InfoContext(ctx, "Command list set successfully.")
	}

	for _, f := range a.before {
		err = f(a)
		if err != nil {
			a.Logger.ErrorContext(ctx, "Initialisation failed.", "error_detail", err)
			return err
		}
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
	for _, f := range a.after {
		f(a)
	}
	a.Logger.Info("Application stopped successfully.")

}

func (a *Application) BeforeStart(initFuncs ...LifeCycleHook) {
	a.before = append(a.before, initFuncs...)
}

func (a *Application) BeforeShutdown(cleanupFuncs ...LifeCycleHook) {
	a.after = append(a.after, cleanupFuncs...)
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

func (a *Application) botSetup() error {

	// Command initialisation
	if len(botCommands) < 1 {
		a.Logger.Warn("No bot commands found.")
		return nil
	}

	bot := botwrapper.NewBotAPIWrapper(a.BotAPI)

	cmds := tgbotapi.NewSetMyCommands(botCommands...)

	// tgbotapi Send method handles message response only.
	// setMyCommands method return boolean.
	// Thus custom setMyCommand function is used here.

	ok, err := bot.SendSetMyCommands(cmds)

	if err != nil {
		return err
	}

	if !ok {
		return errors.New("Cannot set command.")
	}

	// Config initialisation
	// check if yaml or json exists

	config := NewAppConfig()

	if helper.IsPathExists(DefaultConfigFilePathJson) {
		config.FromJson(DefaultConfigFilePathJson)
	} else if helper.IsPathExists(DefaultConfigFilePathYaml) {
		config.FromYaml(DefaultConfigFilePathYaml)
	} else {
		return nil
	}

	for _, c := range config.BotConfigs {
		langCode := strings.TrimSpace(c.LanguageCode)
		if len(langCode) > 0 {
			name := strings.TrimSpace(c.Bot.Name)
			desc := strings.TrimSpace(c.Bot.Description)
			sdesc := strings.TrimSpace(c.Bot.ShortDescription)

			if len(name) > 0 {
				cfg := botwrapper.NewSetMyName(name, langCode)
				ok, err = bot.SetConfig(cfg)
				if err != nil {
					return err
				}
				if !ok {
					return fmt.Errorf("Cannot set name for language code: %s", langCode)
				}
			}

			if len(desc) > 0 {
				cfg := botwrapper.NewSetMyDescription(desc, langCode)
				ok, err = bot.SetConfig(cfg)
				if err != nil {
					return err
				}
				if !ok {
					return fmt.Errorf("Cannot set description for language code: %s", langCode)
				}

			}

			if len(sdesc) > 0 {
				cfg := botwrapper.NewSetMyShortDescription(sdesc, langCode)
				ok, err = bot.SetConfig(cfg)
				if err != nil {
					return err
				}
				if !ok {
					return fmt.Errorf("Cannot set short description for language code: %s", langCode)
				}

			}

		}

	}

	return nil

}
