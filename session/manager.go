package session

import (
	"fmt"
	tgbotapp "go-telegram-bot-app/v1"
	"sync"
)

const (
	ErrSessionNotFound = "Session not found for chat: %d"
)

type SessionManager struct {
	registry map[int64]tgbotapp.Session
	mu       sync.RWMutex
}

func New() tgbotapp.SessionManager {
	return &SessionManager{
		registry: make(map[int64]tgbotapp.Session),
	}
}

func (s *SessionManager) GetOrCreateSession(chatID int64) (*tgbotapp.Session, error) {
	s.mu.RLock()
	sess, ok := s.registry[chatID]
	s.mu.RUnlock()
	if !ok {
		s.mu.Lock()
		defer s.mu.Unlock()
		sess = tgbotapp.Session{
			ChatID: chatID,
			State:  tgbotapp.StateDefault,
		}

		s.registry[chatID] = sess
	}

	return &sess, nil

}

func (s *SessionManager) SetSession(chatID int64, session *tgbotapp.Session) error {

	_, ok := s.registry[chatID]
	if !ok {
		return fmt.Errorf(ErrSessionNotFound, chatID)
	}

	s.registry[chatID] = *session

	return nil
}
