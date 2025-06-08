package server_test

import (
	tgbotapp "go-telegram-bot-app/v1"
	"go-telegram-bot-app/v1/internal/server"
	"go-telegram-bot-app/v1/testutil"
	"testing"
)

func TestContextShouldPanicWithoutLoggerWithErrLoggerNotFound(t *testing.T) {
	ctx := server.NewContext()

	panicValue := testutil.AssertPanic(t, func() {
		ctx.Logger()
	})

	if panicValue != tgbotapp.ErrLoggerNotFound {
		t.Errorf("Expected panic %v, got panic %v", tgbotapp.ErrLoggerNotFound, panicValue)
	}
}
