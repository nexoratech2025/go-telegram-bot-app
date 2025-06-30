package tgbotapp

import (
	"errors"
	"fmt"
	"sync"

	"github.com/StridersTech2025/go-telegram-bot-app/session"
)

const (
	ErrSessionNotFound = "Session not found for chat: %d"
)

var (
	ErrEmptySessionManager = errors.New("Session Manager is nil.")
)

func SessionMiddleware(manager session.SessionManager[int64]) Middleware {

	return func(ctx *BotContext, next HandlerFunc) {

		if manager == nil {
			ctx.Logger().ErrorContext(ctx.Ctx, "No session manager available.")
			next(ctx)
		} else {

			chat := ctx.Update.FromChat()

			if chat != nil {
				chatID := chat.ChatConfig().ChatID
				session, err := manager.GetOrCreate(chatID)
				if err != nil {
					ctx.Logger().WarnContext(ctx.Ctx, "Failed to retrieve session", "error", err)
				}

				ctx.Session = session
				next(ctx)
				manager.Set(chatID, ctx.Session)

			} else {
				ctx.Logger().WarnContext(ctx.Ctx, "Cannot retrieve chatID from chat update", "update_id", ctx.Update.UpdateID)
				next(ctx)
			}
		}

	}

}

type DefaultSession struct {
	data  map[string]any
	state string
}

func NewDefaultSession() session.Sessioner {
	return &DefaultSession{
		data: make(map[string]any),
	}
}

func (s *DefaultSession) CurrentState() string {
	return s.state
}

func (s *DefaultSession) SetState(state string) {
	s.state = state
}

func (s *DefaultSession) Get(key string) (value any, ok bool) {
	value, ok = s.data[key]
	return
}

func (s *DefaultSession) Set(key string, value any) {
	s.data[key] = value
}

func (s *DefaultSession) Delete(key string) {
	delete(s.data, key)
}

func (s *DefaultSession) GetAllKeys() (keys []string) {

	for k := range s.data {
		keys = append(keys, k)
	}

	return
}

// ClearSessionData clears all session data by deleting all keys
func ClearSessionData(session interface {
	GetAllKeys() []string
	Delete(string)
}) {
	for _, key := range session.GetAllKeys() {
		session.Delete(key)
	}
}

// Default Implementation for Session In Memory Manager.
type DefaultInMemoryManager struct {
	registry map[int64]session.Sessioner
	mu       sync.RWMutex
}

func NewDefaultInMemoryManager() session.SessionManager[int64] {
	return &DefaultInMemoryManager{
		registry: make(map[int64]session.Sessioner),
	}
}

func (s *DefaultInMemoryManager) GetOrCreate(chatID int64) (session.Sessioner, error) {
	s.mu.RLock()
	sess, ok := s.registry[chatID]
	s.mu.RUnlock()
	if !ok {
		s.mu.Lock()
		defer s.mu.Unlock()
		sess = NewDefaultSession()
		s.registry[chatID] = sess
	}

	return sess, nil

}

func (s *DefaultInMemoryManager) Set(chatID int64, session session.Sessioner) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	sess, ok := s.registry[chatID]
	if !ok {
		return fmt.Errorf(ErrSessionNotFound, chatID)
	}

	s.registry[chatID] = sess

	return nil
}

func (s *DefaultInMemoryManager) Delete(chatID int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.registry, chatID)
	return nil
}

func (s *DefaultSession) ClearData() {
	ClearSessionData(s)
}

func (s *DefaultSession) ClearState() {
	s.state = ""
}

func (s *DefaultSession) ClearAll() {
	s.ClearData()
	s.ClearState()
}
