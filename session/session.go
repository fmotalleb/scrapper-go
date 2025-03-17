package session

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/fmotalleb/scrapper-go/config"
	"github.com/fmotalleb/scrapper-go/engine"
	"github.com/fmotalleb/scrapper-go/utils"
	"github.com/google/uuid"
)

type Session struct {
	ID             string
	lock           sync.Locker
	ctx            context.Context
	cancel         func()
	sendChannel    chan<- []config.Step
	receiveChannel <-chan map[string]any
	killTimer      *time.Timer
	timeout        time.Duration
}

type SessionStore struct {
	mutex    sync.Mutex
	sessions map[string]*Session
}

var (
	store = &SessionStore{sessions: make(map[string]*Session)}
)

func NewSession(cfg config.ExecutionConfig, timeout time.Duration) (*Session, error) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	sendChannel := make(chan []config.Step)
	receiveChannel, err := engine.ExecuteStream(ctx, cfg, sendChannel)
	if err != nil {
		cancel()
		slog.Error("Failed to execute stream", slog.Any("error", err))
		return nil, err
	}
	timer := time.NewTimer(timeout)
	session := &Session{
		ID:             uuid.New().String(),
		lock:           new(sync.Mutex),
		ctx:            ctx,
		sendChannel:    sendChannel,
		receiveChannel: receiveChannel,
		killTimer:      timer,
		timeout:        timeout,
		cancel:         cancel,
	}
	go session.killerDaemon()

	store.mutex.Lock()
	store.sessions[session.ID] = session
	store.mutex.Unlock()

	slog.Info("New session created", slog.String("session_id", session.ID))
	return session, nil
}

func GetSessions() []string {
	store.mutex.Lock()
	defer store.mutex.Unlock()
	result := make([]string, len(store.sessions))
	index := 0
	for id, _ := range store.sessions {
		result[index] = id
		index++
	}
	return result
}
func GetSession(id string) (*Session, bool) {
	store.mutex.Lock()
	defer store.mutex.Unlock()
	session, exists := store.sessions[id]
	if !exists {
		slog.Warn("Session not found", slog.String("session_id", id))
	} else {
		slog.Debug("Session retrieved", slog.String("session_id", id))
	}
	return session, exists
}

func (s *Session) Handle(steps ...config.Step) (*map[string]any, error) {
	s.resetTimer()
	defer s.resetTimer()
	s.sendChannel <- steps

	slog.Debug("Steps sent to session", slog.String("session_id", s.ID), slog.Any("steps", steps))
	result, err := utils.WithDeadline(s.receiveChannel, s.timeout)
	if err != nil {
		slog.Debug("Operation failed",
			slog.String("session_id", s.ID),
			slog.Any("steps", steps),
			slog.Any("error", err),
		)
		return nil, err
	}
	slog.Debug("Operation completed",
		slog.String("session_id", s.ID),
		slog.Any("steps", steps),
		slog.Any("result", result),
	)

	return result, nil
}

func (s *Session) resetTimer() {
	if !s.killTimer.Stop() {
		<-s.killTimer.C
	}
	s.killTimer.Reset(s.timeout)
	slog.Debug("Kill timer reset",
		slog.String("session_id", s.ID),
		slog.Duration("timeout", s.timeout),
	)
}

func (s *Session) killerDaemon() {
	defer func() {
		store.mutex.Lock()
		delete(store.sessions, s.ID)
		store.mutex.Unlock()

		slog.Debug("Session removed from store",
			slog.String("session_id", s.ID),
		)
	}()
	for {
		select {
		case <-s.ctx.Done():
			slog.Debug("Session context canceled",
				slog.String("session_id", s.ID),
			)
			return
		case <-s.killTimer.C:

			slog.Debug("Session timeout reached, canceling",
				slog.String("session_id", s.ID),
				slog.Duration("timeout", s.timeout),
			)
			s.cancel()
			return
		}
	}
}

func (s *Session) Kill() {
	s.cancel()
	slog.Info("Session manually killed", slog.String("session_id", s.ID))
}
