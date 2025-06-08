package middleware

import (
	"go-telegram-bot-app/v1/internal/server"
	"log"
	"os"
)

func Logger() Middleware {

	return func(context *server.Context, next HandlerFunc) {

		flags := log.Ldate | log.Ltime | log.Lmicroseconds | log.Llongfile

		logger := log.New(os.Stdout, "Telegram Bot App: ", flags)

		context.SetLogger(logger)

		next(context)
	}
}
