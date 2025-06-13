package tgbotapp

import (
	"errors"
	"fmt"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	StateDefault       StateName = ""
	ErrSessionNotFound           = "Session not found for chat: %d"
)

var (
	ErrEmptySessionManager = errors.New("Session Manager is nil.")
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

	if manager == nil {
		panic(ErrEmptySessionManager)
	}

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

// Default Implementation for Session In Memory Manager.

type InMemoryManager struct {
	registry map[int64]Session
	mu       sync.RWMutex
}

func NewInMemoryManager() SessionManager {
	return &InMemoryManager{
		registry: make(map[int64]Session),
	}
}

func (s *InMemoryManager) GetOrCreateSession(chatID int64) (*Session, error) {
	s.mu.RLock()
	sess, ok := s.registry[chatID]
	s.mu.RUnlock()
	if !ok {
		s.mu.Lock()
		defer s.mu.Unlock()
		sess = Session{
			ChatID: chatID,
			State:  StateDefault,
		}

		s.registry[chatID] = sess
	}

	return &sess, nil

}

func (s *InMemoryManager) SetSession(chatID int64, session *Session) error {

	_, ok := s.registry[chatID]
	if !ok {
		return fmt.Errorf(ErrSessionNotFound, chatID)
	}

	s.registry[chatID] = *session

	return nil
}
