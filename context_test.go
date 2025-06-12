package tgbotapp_test

import (
	tgbotapp "go-telegram-bot-app/v1"
	"testing"
)

const (
	TestKey      = "TEST"
	TestData int = 123
)

// func TestContextShouldPanicWithoutLoggerWithErrLoggerNotFound(t *testing.T) {
// 	ctx := tgbotapp.NewBotContext(nil, nil, nil)

// 	panicValue := testutil.AssertPanic(t, func() {
// 		ctx.Logger()
// 	})

// }

func TestContextDataShouldBeSettableAndGettable(t *testing.T) {
	ctx := tgbotapp.NewBotContext(t.Context(), nil, nil)

	ctx.SetData(TestKey, TestData)

	data, ok := ctx.GetData(TestKey)
	if !ok {
		t.Errorf("Expected to get data %d with key %s. Found nothing.", TestData, TestKey)
	}

	switch d := data.(type) {
	case int:
		if d != TestData {
			t.Errorf("Expected %d, found %d", TestData, d)
		}
	default:
		t.Errorf("Unable to convert data to int")
	}

}
