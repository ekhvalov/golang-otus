package memorystorage

import (
	"context"
	"sync"

	"github.com/ekhvalov/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	events map[string]storage.Event
	mu     sync.RWMutex
}

func New() *Storage {
	return &Storage{events: make(map[string]storage.Event)}
}

func (s *Storage) CreateEvent(_ context.Context, e storage.Event) error {
	return s.set(e)
}

func (s *Storage) UpdateEvent(_ context.Context, e storage.Event) error {
	return s.set(e)
}

func (s *Storage) DeleteEvent(_ context.Context, e storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.events, e.ID)
	return nil
}

func (s *Storage) GetEvents(_ context.Context) ([]storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	events := make([]storage.Event, len(s.events))
	i := 0
	for _, event := range s.events {
		events[i] = event
		i++
	}
	return events, nil
}

func (s *Storage) set(e storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.events[e.ID] = e
	return nil
}
