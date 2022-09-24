package memorystorage

import (
	"context"
	"strconv"
	"sync"
	"testing"

	"github.com/ekhvalov/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) {
	ctx := context.Background()
	s := New()

	e1 := storage.Event{ID: "1", Title: "1"}
	err := s.CreateEvent(ctx, e1)
	require.NoError(t, err)
	e2 := storage.Event{ID: "2", Title: "2"}
	err = s.CreateEvent(ctx, e2)
	require.NoError(t, err)

	e1.Title = "one"
	err = s.UpdateEvent(ctx, e1)
	require.NoError(t, err)
	e2.Title = "two"
	err = s.UpdateEvent(ctx, e2)
	require.NoError(t, err)

	events, err := s.GetEvents(ctx)
	require.NoError(t, err)
	require.Len(t, events, 2)
	require.Contains(t, events, e1)
	require.Contains(t, events, e2)
}

func TestConcurrent(t *testing.T) {
	s := New()
	ctx := context.Background()
	wg := &sync.WaitGroup{}
	wg.Add(3)

	go func() {
		defer wg.Done()
		var event storage.Event
		for i := 0; i < 10_000; i++ {
			event.ID = strconv.Itoa(i)
			event.Title = strconv.Itoa(i)
			_ = s.CreateEvent(ctx, event)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 10_000; i++ {
			_, _ = s.GetEvents(ctx)
		}
	}()

	go func() {
		defer wg.Done()
		var event storage.Event
		for i := 0; i < 10_000; i++ {
			event.ID = strconv.Itoa(i)
			event.Title = strconv.Itoa(i)
			_ = s.DeleteEvent(ctx, event)
		}
	}()

	wg.Wait()
}
