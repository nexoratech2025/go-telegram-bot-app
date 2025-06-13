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
	Params  []string
	Bot     *tgbotapi.BotAPI
	Update  *tgbotapi.Update
	Session *Session
}

func NewBotContext(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) *BotContext {

	ctxID := uuid.NewString()
	ctx = context.WithValue(ctx, CtxKeyRequestID, ctxID)

	return &BotContext{
		Ctx:  ctx,
		data: make(map[string]any),

		Bot:    bot,
		Update: update,
	}
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
