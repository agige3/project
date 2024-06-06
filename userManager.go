package userManager

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/redis/go-redis/v9"

	"project/internal/handler"
	"project/internal/worker"
)

const (
	defaultExpirationTime = time.Minute * 5
)

var ErrWorkingWithNilManager = errors.New("You are trying to work with nil manager. Use userManager.New().")

type Manager struct {
	handler *handler.Handler
}

func NewManager(db *sql.DB, client *redis.Client) (*Manager, error) {
	return NewManagerWithExpirationTime(db, client, defaultExpirationTime)
}

func NewManagerWithExpirationTime(db *sql.DB, client *redis.Client, expirationTime time.Duration) (*Manager, error) {
	err := db.Ping()
	if err != nil {
		return nil, err
	}
	res := client.Ping(context.Background())
	if res.Err() != nil {
		return nil, err
	}
	worker := worker.NewWorker(client, db, expirationTime)
	handler := handler.NewHandler(worker)

	return &Manager{handler: handler}, err
}

func (m *Manager) SetExpirationTime(expirationTime time.Duration) error {
	if m == nil {
		return ErrWorkingWithNilManager
	}
	m.handler.Worker.ActualTime = expirationTime
}

func (m *Manager) StartServer(addr string) error {
	if m == nil {
		return ErrWorkingWithNilManager
	}
	r := chi.NewRouter()
	r.Get("/", m.handler.GetUsers)
	r.Post("/", m.handler.AddUser)
	err := http.ListenAndServe(addr, r)
	return err
}
