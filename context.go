package tgbotapp

import (
	"context"
	"log/slog"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
)

const (
	CtxKeyRequestID = "request_id"
)

type HandlerFunc func(*BotContext)

type BotContext struct {
	data map[string]any
	app  *Application

	Ctx     context.Context
	BotAPI  *tgbotapi.BotAPI
	Update  *tgbotapi.Update
	Session *Session
	Params  []string
}

func NewBotContext(ctx context.Context, app *Application, update *tgbotapi.Update) *BotContext {

	ctxID := uuid.NewString()
	ctx = context.WithValue(ctx, CtxKeyRequestID, ctxID)

	c := &BotContext{
		data: make(map[string]any),
		app:  app,

		Ctx:    ctx,
		Update: update,
	}

	if app != nil {
		c.BotAPI = app.BotAPI
	}

	return c
}

func (c *BotContext) SetData(key string, value any) {
	c.data[key] = value
}

func (c *BotContext) GetData(key string) (value any, ok bool) {
	value, ok = c.data[key]
	return
}

func (c *BotContext) Logger() *slog.Logger {
	return c.app.Logger
}

func (c *BotContext) SetHandler(f HandlerFunc) {
	c.app.handler = f
}
