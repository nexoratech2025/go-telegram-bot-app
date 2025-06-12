package tgbotapp

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

const (
	StateDefault StateName = ""
)

type Session struct {
	ChatID int64
	State  StateName
}

type SessionManager interface {
	GetOrCreateSession(chatID int64) (*Session, error)
	SetSession(chatID int64, session *Session) error
}

func SessionMiddleware(manager SessionManager) Middleware {

	return func(ctx *BotContext, next HandlerFunc) {

		chatID, ok := tryGetChatID(ctx.Update)

		if ok {
			session, err := manager.GetOrCreateSession(chatID)
			if err != nil {
				ctx.Logger().WarnContext(ctx.Ctx, "Faild to retrieve session", "error", err)
			}

			ctx.Session = session

		} else {
			ctx.Logger().WarnContext(ctx.Ctx, "Cannot retrieve chatID from chat update", "updateId", ctx.Update.UpdateID)
		}

		next(ctx)
	}

}

func tryGetChatID(update *tgbotapi.Update) (chatID int64, ok bool) {

	if update.CallbackQuery != nil {
		chatID = update.CallbackQuery.Message.Chat.ID
		ok = true

		return
	}

	if update.Message != nil {
		chatID = update.Message.Chat.ID
		ok = true

		return
	}

	return

}
