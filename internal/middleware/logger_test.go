package middleware_test

import (
	"go-telegram-bot-app/v1/internal/middleware"
	"go-telegram-bot-app/v1/internal/server"
	"testing"
)

func printLogHandler(ctx *server.Context) {
	ctx.Logger().Println("Called from print log handler")

}

func TestLoggerMiddlewareShouldPassLoggerInContext(t *testing.T) {

	ctx := server.NewContext()

	m := middleware.New()
	m.Append(middleware.Logger())

	f := m.Wrap(printLogHandler)

	f(ctx)

}
