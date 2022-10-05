package memorystorage_test

import (
	"context"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/ekhvalov/hw12_13_14_15_calendar/internal/domain/event"
	memorystorage "github.com/ekhvalov/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/stretchr/testify/require"
)

type idProvider struct {
	id int
}

func (p *idProvider) GenerateID() (string, error) {
	p.id++
	return strconv.Itoa(p.id), nil
}

func TestConcurrent(t *testing.T) {
	s := memorystorage.New(&idProvider{})
	ctx := context.Background()
	wg := &sync.WaitGroup{}
	wg.Add(3)

	go func() {
		defer wg.Done()
		var e event.Event
		for i := 0; i < 2000; i++ {
			e.ID = strconv.Itoa(i)
			e.Title = strconv.Itoa(i)
			_, _ = s.Create(ctx, e)
		}
	}()

	go func() {
		defer wg.Done()
		now := time.Now()
		for i := 0; i < 2000; i++ {
			_, _ = s.GetDayEvents(ctx, now)
		}
	}()

	go func() {
		defer wg.Done()
		var e event.Event
		for i := 0; i < 2000; i++ {
			e.ID = strconv.Itoa(i)
			e.Title = strconv.Itoa(i)
			_ = s.Delete(ctx, e.ID)
		}
	}()

	wg.Wait()
}

func TestStorage_GetDayEvents(t *testing.T) {
	event1Day1 := event.Event{
		ID:       "1",
		DateTime: time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC),
	}
	event2Day1 := event.Event{
		ID:       "2",
		DateTime: time.Date(2022, time.January, 1, 0, 0, 1, 0, time.UTC),
	}
	event3Day1 := event.Event{
		ID:       "3",
		DateTime: time.Date(2022, time.January, 1, 23, 59, 59, 0, time.UTC),
	}
	event4Day2 := event.Event{
		ID:       "4",
		DateTime: time.Date(2022, time.January, 2, 12, 59, 59, 0, time.UTC),
	}
	events := []event.Event{event1Day1, event2Day1, event3Day1, event4Day2}
	tests := map[string]struct {
		date    time.Time
		want    []event.Event
		wantErr bool
	}{
		"2022-01-01": {
			date: time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC),
			want: []event.Event{event1Day1, event2Day1, event3Day1},
		},
		"2022-01-02": {
			date: time.Date(2022, time.January, 2, 0, 0, 0, 0, time.UTC),
			want: []event.Event{event4Day2},
		},
		"2022-01-03": {
			date: time.Date(2022, time.January, 3, 0, 0, 0, 0, time.UTC),
			want: []event.Event{},
		},
	}
	for testName, tt := range tests {
		t.Run(testName, func(t *testing.T) {
			s := memorystorage.New(&idProvider{})
			for _, e := range events {
				_, err := s.Create(context.Background(), e)
				require.NoError(t, err)
			}
			got, err := s.GetDayEvents(context.Background(), tt.date)
			require.NoError(t, err)
			require.ElementsMatch(t, tt.want, got)
		})
	}
}

func TestStorage_GetWeekEvents(t *testing.T) {
	event1Week1 := event.Event{
		ID:       "1",
		DateTime: time.Date(2022, time.January, 3, 0, 0, 0, 0, time.UTC),
	}
	event2Week1 := event.Event{
		ID:       "2",
		DateTime: time.Date(2022, time.January, 5, 0, 0, 1, 0, time.UTC),
	}
	event3Week1 := event.Event{
		ID:       "3",
		DateTime: time.Date(2022, time.January, 9, 23, 59, 59, 0, time.UTC),
	}
	event4Week2 := event.Event{
		ID:       "4",
		DateTime: time.Date(2022, time.January, 10, 12, 59, 59, 0, time.UTC),
	}
	events := []event.Event{event1Week1, event2Week1, event3Week1, event4Week2}
	tests := map[string]struct {
		date    time.Time
		want    []event.Event
		wantErr bool
	}{
		"2022-01-03": {
			date: time.Date(2022, time.January, 3, 0, 0, 0, 0, time.UTC),
			want: []event.Event{event1Week1, event2Week1, event3Week1},
		},
		"2022-01-10": {
			date: time.Date(2022, time.January, 10, 0, 0, 0, 0, time.UTC),
			want: []event.Event{event4Week2},
		},
		"2022-01-17": {
			date: time.Date(2022, time.January, 17, 0, 0, 0, 0, time.UTC),
			want: []event.Event{},
		},
	}
	for testName, tt := range tests {
		t.Run(testName, func(t *testing.T) {
			s := memorystorage.New(&idProvider{})
			for _, e := range events {
				_, err := s.Create(context.Background(), e)
				require.NoError(t, err)
			}
			got, err := s.GetWeekEvents(context.Background(), tt.date)
			require.NoError(t, err)
			require.ElementsMatch(t, tt.want, got)
		})
	}
}

func TestStorage_GetMonthEvents(t *testing.T) {
	event1Month1 := event.Event{
		ID:       "1",
		DateTime: time.Date(2022, time.January, 3, 0, 0, 0, 0, time.UTC),
	}
	event2Month1 := event.Event{
		ID:       "2",
		DateTime: time.Date(2022, time.January, 5, 0, 0, 1, 0, time.UTC),
	}
	event3Month1 := event.Event{
		ID:       "3",
		DateTime: time.Date(2022, time.January, 9, 23, 59, 59, 0, time.UTC),
	}
	event4Month2 := event.Event{
		ID:       "4",
		DateTime: time.Date(2022, time.February, 10, 12, 59, 59, 0, time.UTC),
	}
	events := []event.Event{event1Month1, event2Month1, event3Month1, event4Month2}
	tests := map[string]struct {
		date    time.Time
		want    []event.Event
		wantErr bool
	}{
		"2022-01-03": {
			date: time.Date(2022, time.January, 3, 0, 0, 0, 0, time.UTC),
			want: []event.Event{event1Month1, event2Month1, event3Month1},
		},
		"2022-02-10": {
			date: time.Date(2022, time.February, 10, 0, 0, 0, 0, time.UTC),
			want: []event.Event{event4Month2},
		},
		"2022-03-17": {
			date: time.Date(2022, time.March, 17, 0, 0, 0, 0, time.UTC),
			want: []event.Event{},
		},
	}
	for testName, tt := range tests {
		t.Run(testName, func(t *testing.T) {
			s := memorystorage.New(&idProvider{})
			for _, e := range events {
				_, err := s.Create(context.Background(), e)
				require.NoError(t, err)
			}
			got, err := s.GetMonthEvents(context.Background(), tt.date)
			require.NoError(t, err)
			require.ElementsMatch(t, tt.want, got)
		})
	}
}

func TestStorage_Create(t *testing.T) {
	event1000to1130 := event.Event{
		DateTime: time.Date(2022, time.January, 1, 10, 0, 0, 0, time.UTC),
		Duration: time.Minute * 90,
	}
	event1230to1530 := event.Event{
		DateTime: time.Date(2022, time.January, 1, 12, 30, 0, 0, time.UTC),
		Duration: time.Hour * 3,
	}
	events := []event.Event{event1000to1130, event1230to1530}
	tests := map[string]struct {
		date     time.Time
		duration time.Duration
		err      error
	}{
		"09:00 -> 10:00 OK": {
			date:     time.Date(2022, time.January, 1, 9, 0, 0, 0, time.UTC),
			duration: time.Hour,
		},
		"11:30 -> 12:00 OK": {
			date:     time.Date(2022, time.January, 1, 11, 30, 0, 0, time.UTC),
			duration: time.Minute * 30,
		},
		"11:30 -> 12:30 OK": {
			date:     time.Date(2022, time.January, 1, 11, 30, 0, 0, time.UTC),
			duration: time.Minute * 60,
		},
		"11:30 -> 12:31 Err": {
			date:     time.Date(2022, time.January, 1, 11, 30, 0, 0, time.UTC),
			duration: time.Minute * 61,
			err:      event.ErrDateBusy,
		},
		"09:00 -> 10:01 Err": {
			date:     time.Date(2022, time.January, 1, 9, 0, 0, 0, time.UTC),
			duration: time.Minute * 61,
			err:      event.ErrDateBusy,
		},
		"09:00 -> 12:00 Err": {
			date:     time.Date(2022, time.January, 1, 9, 0, 0, 0, time.UTC),
			duration: time.Hour * 3,
			err:      event.ErrDateBusy,
		},
		"09:00 -> 16:00 Err": {
			date:     time.Date(2022, time.January, 1, 9, 0, 0, 0, time.UTC),
			duration: time.Hour * 7,
			err:      event.ErrDateBusy,
		},
		"10:00 -> 11:30 Err": {
			date:     time.Date(2022, time.January, 1, 10, 0, 0, 0, time.UTC),
			duration: time.Minute * 90,
			err:      event.ErrDateBusy,
		},
		"10:30 -> 11:00 Err": {
			date:     time.Date(2022, time.January, 1, 10, 30, 0, 0, time.UTC),
			duration: time.Minute * 30,
			err:      event.ErrDateBusy,
		},
		"10:30 -> 15:00 Err": {
			date:     time.Date(2022, time.January, 1, 10, 30, 0, 0, time.UTC),
			duration: time.Hour*4 + time.Minute*30,
			err:      event.ErrDateBusy,
		},
		"15:00 -> 16:00 Err": {
			date:     time.Date(2022, time.January, 1, 15, 0, 0, 0, time.UTC),
			duration: time.Hour,
			err:      event.ErrDateBusy,
		},
	}
	for testName, tt := range tests {
		t.Run(testName, func(t *testing.T) {
			s := memorystorage.New(&idProvider{})
			for _, e := range events {
				_, err := s.Create(context.Background(), e)
				require.NoError(t, err)
			}
			e, err := s.Create(context.Background(), event.Event{
				DateTime: tt.date,
				Duration: tt.duration,
			})
			if tt.err != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tt.err)
			} else {
				require.NoError(t, err)
				require.Equal(t, event.Event{
					ID:       "3",
					DateTime: tt.date,
					Duration: tt.duration,
				}, e)
			}
		})
	}
}

func TestStorage_Update(t *testing.T) {
	event1000to1130 := event.Event{
		DateTime: time.Date(2022, time.January, 1, 10, 0, 0, 0, time.UTC),
		Duration: time.Minute * 90,
	}
	event1230to1530 := event.Event{
		DateTime: time.Date(2022, time.January, 1, 12, 30, 0, 0, time.UTC),
		Duration: time.Hour * 3,
	}
	updatableEventID := "3"
	updatableEvent := event.Event{
		ID:       updatableEventID,
		DateTime: time.Date(2022, time.January, 1, 11, 30, 0, 0, time.UTC),
		Duration: time.Hour,
	}
	events := []event.Event{event1000to1130, event1230to1530, updatableEvent}
	tests := map[string]struct {
		date     time.Time
		duration time.Duration
		err      error
	}{
		"09:00 -> 10:00 OK": {
			date:     time.Date(2022, time.January, 1, 9, 0, 0, 0, time.UTC),
			duration: time.Hour,
		},
		"11:30 -> 12:00 OK": {
			date:     time.Date(2022, time.January, 1, 11, 30, 0, 0, time.UTC),
			duration: time.Minute * 30,
		},
		"11:30 -> 12:30 OK": {
			date:     time.Date(2022, time.January, 1, 11, 30, 0, 0, time.UTC),
			duration: time.Minute * 60,
		},
		"11:30 -> 12:31 Err": {
			date:     time.Date(2022, time.January, 1, 11, 30, 0, 0, time.UTC),
			duration: time.Minute * 61,
			err:      event.ErrDateBusy,
		},
		"09:00 -> 10:01 Err": {
			date:     time.Date(2022, time.January, 1, 9, 0, 0, 0, time.UTC),
			duration: time.Minute * 61,
			err:      event.ErrDateBusy,
		},
		"09:00 -> 12:00 Err": {
			date:     time.Date(2022, time.January, 1, 9, 0, 0, 0, time.UTC),
			duration: time.Hour * 3,
			err:      event.ErrDateBusy,
		},
		"09:00 -> 16:00 Err": {
			date:     time.Date(2022, time.January, 1, 9, 0, 0, 0, time.UTC),
			duration: time.Hour * 7,
			err:      event.ErrDateBusy,
		},
		"10:00 -> 11:30 Err": {
			date:     time.Date(2022, time.January, 1, 10, 0, 0, 0, time.UTC),
			duration: time.Minute * 90,
			err:      event.ErrDateBusy,
		},
		"10:30 -> 11:00 Err": {
			date:     time.Date(2022, time.January, 1, 10, 30, 0, 0, time.UTC),
			duration: time.Minute * 30,
			err:      event.ErrDateBusy,
		},
		"10:30 -> 15:00 Err": {
			date:     time.Date(2022, time.January, 1, 10, 30, 0, 0, time.UTC),
			duration: time.Hour*4 + time.Minute*30,
			err:      event.ErrDateBusy,
		},
		"15:00 -> 16:00 Err": {
			date:     time.Date(2022, time.January, 1, 15, 0, 0, 0, time.UTC),
			duration: time.Hour,
			err:      event.ErrDateBusy,
		},
	}
	for testName, tt := range tests {
		t.Run(testName, func(t *testing.T) {
			s := memorystorage.New(&idProvider{})
			for _, e := range events {
				_, err := s.Create(context.Background(), e)
				require.NoError(t, err)
			}
			err := s.Update(context.Background(), updatableEventID, event.Event{
				DateTime: tt.date,
				Duration: tt.duration,
			})
			if tt.err != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tt.err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
