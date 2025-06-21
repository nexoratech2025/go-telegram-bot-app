package tgbotapp

import (
	"errors"
	"fmt"
	"sync"
)

type StateName string

const (
	StateDefault StateName = ""
)
const (
	ErrSessionNotFound = "Session not found for chat: %d"
)

var (
	ErrEmptySessionManager = errors.New("Session Manager is nil.")
)

type Session struct {
	ChatID int64
	State  StateName
	Data   map[string]any
}

type SessionManager interface {
	GetOrCreateSession(chatID int64) (*Session, error)
	SetSession(chatID int64, session *Session) error
}

func SessionMiddleware(manager SessionManager) Middleware {

	return func(ctx *BotContext, next HandlerFunc) {

		if manager == nil {
			ctx.Logger().ErrorContext(ctx.Ctx, "No session manager available.")
			next(ctx)
		} else {

			chat := ctx.Update.FromChat()

			if chat != nil {
				chatID := chat.ID
				session, err := manager.GetOrCreateSession(chatID)
				if err != nil {
					ctx.Logger().WarnContext(ctx.Ctx, "Faild to retrieve session", "error", err)
				}

				ctx.Session = session
				next(ctx)
				manager.SetSession(chatID, ctx.Session)

			} else {
				ctx.Logger().WarnContext(ctx.Ctx, "Cannot retrieve chatID from chat update", "update_id", ctx.Update.UpdateID)
				next(ctx)
			}
		}
	}
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

	sess, ok := s.registry[chatID]
	if !ok {
		return fmt.Errorf(ErrSessionNotFound, chatID)
	}

	sess.State = session.State
	sess.Data = session.Data

	s.registry[chatID] = sess

	return nil
}
