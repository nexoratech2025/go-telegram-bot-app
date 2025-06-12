package tgbotapp

import (
	"log"
	"log/slog"

	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	loggerFlags = log.Ldate | log.Ltime | log.Lmicroseconds | log.Llongfile
)

var (
	botCommands []tgbotapi.BotCommand

	defaultLoggers = slog.Default()
)

type Application struct {
	slog        *slog.Logger
	router      *RouteTable
	middlewares MiddlewareChain
	handler     HandlerFunc
	bot         *tgbotapi.BotAPI
}

func new(logger *slog.Logger) *Application {
	return &Application{
		slog:        logger,
		router:      NewRouteTable(),
		middlewares: *NewMiddlewareChain(),
	}

}

func Default() *Application {

	app := new(defaultLoggers)
	app.UseSession()
	app.UseRouting()
	return app
}

func (a *Application) RegisterCommand(name CommandName, description string, handler HandlerFunc) {

	botCommands = append(botCommands, tgbotapi.BotCommand{
		Command:     string(name),
		Description: description,
	})

	a.router.AddCommandHandler(name, handler)
}

func (a *Application) RegisterCallback(name CallbackName, handler HandlerFunc) {
	a.router.AddCallbackHandler(name, handler)
}

func (a *Application) RegisterMessage(state StateName, handler HandlerFunc) {
	a.router.AddMessageHandler(state, handler)
}

func (a *Application) Use(middlewares ...Middleware) {
	a.middlewares.Append(middlewares...)
}

func (a *Application) UseRouting() {

	a.middlewares.Append(Router(a.router))

}

func (a *Application) UseSession() {
	a.middlewares.Append()

}

// TODO: Try to make update handling concurrent without problems

func (a *Application) Start(ctx context.Context) error {

	updateCfg := tgbotapi.NewUpdate(0)
	updateCfg.Timeout = 60

	updates := a.bot.GetUpdatesChan(updateCfg)

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
	a.bot.StopReceivingUpdates()
	return nil
}

func (a *Application) handleUpdate(ctx context.Context, update *tgbotapi.Update) {

	botCtx := NewBotContext(ctx, a.bot, update)

	f := a.middlewares.Wrap(func(ctx *BotContext) {
		if ctx.app.handler != nil {
			ctx.app.handler(botCtx)
		}
		ctx.Logger().WarnContext(ctx.Ctx, "No Handler Found. Returning nothing.")
	})

	f(botCtx)

}
