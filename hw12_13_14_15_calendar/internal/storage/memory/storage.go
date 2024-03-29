package memorystorage

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ekhvalov/golang-otus/hw12_13_14_15_calendar/internal/domain/event"
)

type Storage struct {
	events     map[string]event.Event
	mu         sync.RWMutex
	idProvider IDProvider
}

var _ event.Storage = (*Storage)(nil)

func New(idProvider IDProvider) *Storage {
	return &Storage{
		idProvider: idProvider,
		events:     make(map[string]event.Event),
	}
}

func (s *Storage) Create(_ context.Context, e event.Event) (event.Event, error) {
	id, err := s.idProvider.GenerateID()
	if err != nil {
		return event.Event{}, fmt.Errorf("generate ID error: %w", err)
	}
	e.ID = id
	s.mu.RLock()
	for _, e2 := range s.events {
		if isOverlapped(e, e2) {
			s.mu.RUnlock()
			return event.Event{}, event.ErrDateBusy
		}
	}
	if _, ok := s.events[e.ID]; ok {
		return event.Event{}, fmt.Errorf("event with id '%s' is already exist", e.ID)
	}
	s.mu.RUnlock()
	s.mu.Lock()
	defer s.mu.Unlock()
	s.events[e.ID] = e
	return e, nil
}

func (s *Storage) Update(_ context.Context, eventID string, e event.Event) error {
	s.mu.RLock()
	for _, e2 := range s.events {
		if eventID != e2.ID && isOverlapped(e, e2) {
			s.mu.RUnlock()
			return event.ErrDateBusy
		}
	}
	if _, ok := s.events[eventID]; !ok {
		return event.ErrEventNotFound
	}
	s.mu.RUnlock()
	e.ID = eventID
	s.mu.Lock()
	s.events[eventID] = e
	s.mu.Unlock()
	return nil
}

func (s *Storage) Delete(_ context.Context, eventID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.events[eventID]; !ok {
		return fmt.Errorf("event with id '%s' is not exist", eventID)
	}
	delete(s.events, eventID)
	return nil
}

func (s *Storage) GetDayEvents(_ context.Context, date time.Time) ([]event.Event, error) {
	return s.filterByTimeRange(date.Unix(), date.AddDate(0, 0, 1).Unix()), nil
}

func (s *Storage) GetWeekEvents(_ context.Context, date time.Time) ([]event.Event, error) {
	return s.filterByTimeRange(date.Unix(), date.AddDate(0, 0, 7).Unix()), nil
}

func (s *Storage) GetMonthEvents(_ context.Context, date time.Time) ([]event.Event, error) {
	return s.filterByTimeRange(date.Unix(), date.AddDate(0, 1, 0).Unix()), nil
}

func (s *Storage) DeleteEventsOlderThan(_ context.Context, date time.Time) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for key, e := range s.events {
		if e.DateTime.Before(date) {
			delete(s.events, key)
		}
	}
	return nil
}

func (s *Storage) GetEventsNotifyBetween(_ context.Context, from time.Time, to time.Time) ([]event.Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	events := make([]event.Event, 0)
	f := from.Unix()
	t := to.Unix()
	for _, e := range s.events {
		if e.Duration.Seconds() > 0 {
			nt := getNotifyTime(e)
			if nt >= f && nt <= t {
				events = append(events, e)
			}
		}
	}
	return events, nil
}

func (s *Storage) filterByTimeRange(start, end int64) []event.Event {
	events := make([]event.Event, 0)
	s.mu.RLock()
	for _, e := range s.events {
		if e.DateTime.Unix() >= start && e.DateTime.Unix() < end {
			events = append(events, e)
		}
	}
	s.mu.RUnlock()
	return events
}

func isOverlapped(e1, e2 event.Event) bool {
	startA := e1.DateTime.Unix()
	endA := e1.DateTime.Add(e1.Duration).Unix()
	startB := e2.DateTime.Unix()
	endB := e2.DateTime.Add(e2.Duration).Unix()
	if startA >= startB && startA < endB {
		return true
	}
	if endA > startB && endA <= endB {
		return true
	}
	if startB >= startA && startB < endA {
		return true
	}
	if endB > startA && endB <= endA {
		return true
	}
	return false
}

func getNotifyTime(e event.Event) int64 {
	return e.DateTime.Unix() - int64(e.Duration.Seconds())
}
