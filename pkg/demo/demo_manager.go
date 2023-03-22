package demo

import (
	"context"
	"crypto/md5"
	"fmt"
	"math/rand"
	"time"

	"mackey/pkg/db"
	"mackey/pkg/logger"

	"go.uber.org/zap"
)

type Config struct {
	PushTimeout time.Duration `env:"DEMO_PUSH_TIMEOUT,default=1s"`
	PullTimeout time.Duration `env:"DEMO_PULL_TIMEOUT,default=3s"`
}

type Manager struct {
	conf Config
	repo db.DemoRepo
	log  *logger.Logger
}

func New(conf Config, repo db.DemoRepo, log *logger.Logger) *Manager {
	return &Manager{
		conf: conf,
		repo: repo,
		log:  log,
	}
}

func (m *Manager) Run() error {
	ctx := context.Background()
	m.log.Info(ctx, "start manager")

	go func() {
		pusher := time.NewTicker(m.conf.PushTimeout)
		for range pusher.C {
			if err := m.Push(); err == nil {
				m.log.Info(ctx, "pushed new message")
			}
		}
	}()

	puller := time.NewTicker(m.conf.PullTimeout)
	for range puller.C {
		cnt, err := m.Pull()
		if err != nil {
			return err
		}

		m.log.Info(ctx, fmt.Sprintf("pulled %d message(s)", cnt))
	}

	return nil
}

func (m *Manager) Push() error {
	userID := rand.Int()

	ctx := context.WithValue(context.Background(), "userId", userID)
	_, err := m.repo.AddMessage(ctx, &db.Message{
		Text: fmt.Sprintf("hello from %d", userID),
	})

	return err
}

func (m *Manager) Pull() (int, error) {
	ctx := context.Background()
	messages, err := m.repo.MessagesByFilters(ctx, nil, db.PagerDefault)
	if err != nil {
		return 0, err
	}

	for _, message := range messages {
		m.log.Info(ctx, fmt.Sprintf("got message: %s", message.Text), zap.Int("messageId", message.ID))
		if _, err := m.repo.DeleteMessage(ctx, message.ID); err != nil {
			return 0, err
		}
	}

	return len(messages), nil
}

func (m *Manager) Hash(id string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(id)))
}
