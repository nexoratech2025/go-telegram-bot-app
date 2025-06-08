package middleware_test

import (
	"go-telegram-bot-app/v1/internal/middleware"
	"go-telegram-bot-app/v1/internal/server"
	"log"
	"testing"

	"github.com/google/uuid"
)

const (
	dataKey       = "Test"
	expectedValue = "ACEDB"
)

func middlewareFunc1(ctx *server.Context, next middleware.HandlerFunc) {

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

	ctx.Logger.Println("Middleware 1 end")

}

func middlewareFunc2(ctx *server.Context, next middleware.HandlerFunc) {
	ctx.Logger.Println("Middleware 2 start")
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
	ctx.Logger.Println("Middleware 2 end")

}

func handlerFunc(ctx *server.Context) {
	v, ok := ctx.GetData(dataKey)
	if ok {
		e := v.(string)
		e += "E"
		ctx.SetData(dataKey, e)
	}
	ctx.Logger.Println("Handler Function called.")
}

func TestMiddlewareChainCorrectOrder(t *testing.T) {

	t.Logf("Running Test %s", uuid.NewString())

	chain := middleware.New()
	ctx := server.NewContext()

	ctx.Logger = log.Default()
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
