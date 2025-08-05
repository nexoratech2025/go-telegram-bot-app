package tgbotapp_test

import (
	"log/slog"
	"testing"

	tgbotapp "github.com/nexoratech2025/go-telegram-bot-app"

	"github.com/google/uuid"
)

const (
	dataKey       = "Test"
	expectedValue = "ACEDB"
)

func middlewareFunc1(ctx *tgbotapp.BotContext, next tgbotapp.HandlerFunc) {
	ctx.Logger().Debug("Middleware 1 start")
	v, _ := ctx.GetData(dataKey)

	a, ok := v.(string)
	if !ok {
		a = ""
	}
	a += "A"
	ctx.SetData(dataKey, a)

	next(ctx)

	v, ok = ctx.GetData(dataKey)

	if ok {
		b := v.(string)
		b += "B"
		ctx.SetData(dataKey, b)
	}

	ctx.Logger().Debug("Middleware 1 end")

}

func middlewareFunc2(ctx *tgbotapp.BotContext, next tgbotapp.HandlerFunc) {
	ctx.Logger().Debug("Middleware 2 start")
	v, ok := ctx.GetData(dataKey)
	if ok {
		c := v.(string)
		c += "C"
		ctx.SetData(dataKey, c)
	}

	next(ctx)
	v, ok = ctx.GetData(dataKey)
	if ok {
		d := v.(string)
		d += "D"
		ctx.SetData(dataKey, d)
	}
	ctx.Logger().Debug("Middleware 2 end")

}

func handlerFunc(ctx *tgbotapp.BotContext) {
	v, ok := ctx.GetData(dataKey)
	if ok {
		e := v.(string)
		e += "E"
		ctx.SetData(dataKey, e)
	}
	ctx.Logger().Debug("Handler Function called.")
}

func TestMiddlewareChainCorrectOrder(t *testing.T) {

	t.Logf("Running Test %s", uuid.NewString())

	chain := tgbotapp.NewMiddlewareChain()
	app := tgbotapp.Application{
		Logger: slog.Default(),
	}
	ctx := tgbotapp.NewBotContext(t.Context(), &app, nil)

	chain.Append(middlewareFunc1, middlewareFunc2)
	m := chain.Wrap(handlerFunc)

	m(ctx)

	v, ok := ctx.GetData(dataKey)
	if !ok {
		t.Errorf("Expects data, found nothing")
	}

	data := v.(string)

	t.Logf("Got String: %s", data)

	if data != expectedValue {
		t.Errorf("Expects %s. Found %s", expectedValue, data)
	}

}
