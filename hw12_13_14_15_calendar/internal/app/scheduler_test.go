package app

import (
	"context"
	"testing"
	"time"

	"github.com/ekhvalov/golang-otus/hw12_13_14_15_calendar/internal/app/notification/queue"
	queuemock "github.com/ekhvalov/golang-otus/hw12_13_14_15_calendar/internal/app/notification/queue/mock"
	"github.com/ekhvalov/golang-otus/hw12_13_14_15_calendar/internal/domain/event"
	"github.com/ekhvalov/golang-otus/hw12_13_14_15_calendar/internal/domain/event/mock"
	"github.com/ekhvalov/golang-otus/hw12_13_14_15_calendar/internal/domain/notification"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func Test_scheduler_FindNotificationReadyEvents(t *testing.T) {
	errStorage := event.NewErrStorage("storage error")
	errProducer := queue.NewErrProducer("producer error")
	tests := map[string]struct {
		getStorage       func(controller *gomock.Controller) event.Storage
		getProducer      func(controller *gomock.Controller) queue.Producer
		contextTimeout   time.Duration
		getEventsTimeout time.Duration
		err              error
	}{
		"when storage error occurred then should return error": {
			getStorage: func(controller *gomock.Controller) event.Storage {
				s := mock.NewMockStorage(controller)
				s.EXPECT().
					GetEventsNotifyBetween(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, errStorage)
				return s
			},
			getProducer: func(controller *gomock.Controller) queue.Producer {
				p := queuemock.NewMockProducer(controller)
				p.EXPECT().Close().Return(nil)
				return p
			},
			contextTimeout:   time.Millisecond * 100,
			getEventsTimeout: time.Millisecond * 10,
			err:              errStorage,
		},
		"when queue error occurred then should return error": {
			getStorage: func(controller *gomock.Controller) event.Storage {
				s := mock.NewMockStorage(controller)
				s.EXPECT().
					GetEventsNotifyBetween(gomock.Any(), gomock.Any(), gomock.Any()).
					Return([]event.Event{{Title: "event 1"}}, nil)
				return s
			},
			getProducer: func(controller *gomock.Controller) queue.Producer {
				p := queuemock.NewMockProducer(controller)
				p.EXPECT().
					Put(gomock.Any()).
					Return(errProducer)
				p.EXPECT().Close().Return(nil)
				return p
			},
			contextTimeout:   time.Millisecond * 100,
			getEventsTimeout: time.Millisecond * 10,
			err:              errProducer,
		},
		"when context timed out then no calls should be done": {
			getStorage: func(controller *gomock.Controller) event.Storage {
				return mock.NewMockStorage(controller)
			},
			getProducer: func(controller *gomock.Controller) queue.Producer {
				p := queuemock.NewMockProducer(controller)
				p.EXPECT().Close().Return(nil)
				return p
			},
			contextTimeout:   time.Millisecond * 10,
			getEventsTimeout: time.Millisecond * 100,
			err:              nil,
		},
	}
	for testName, tt := range tests {
		t.Run(testName, func(t *testing.T) {
			c := gomock.NewController(t)
			s, err := NewScheduler(tt.getStorage(c), tt.getProducer(c))
			require.NoError(t, err)
			s.(*scheduler).scanInterval = tt.getEventsTimeout
			ctx, cancel := context.WithTimeout(context.Background(), tt.contextTimeout)
			defer cancel()

			err = s.FindNotificationReadyEvents(ctx)

			if tt.err != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tt.err)
			} else {
				require.NoError(t, err)
			}
		})
	}

	t.Run("when events provided then they should be put to the queue provider", func(t *testing.T) {
		event1 := event.Event{Title: "event 1", ID: "1"}
		event2 := event.Event{Title: "event 2", ID: "2"}
		event3 := event.Event{Title: "event 3", ID: "3"}
		event4 := event.Event{Title: "event 4", ID: "4"}
		eventsToReturn := [][]event.Event{{event1, event2}, {}, {event3}, {event4}}
		callNumber := 0
		expectedEventsTitles := []string{event1.Title, event2.Title, event3.Title, event4.Title}
		actualEventsTitles := make([]string, 0)

		controller := gomock.NewController(t)
		producer := queuemock.NewMockProducer(controller)
		producer.EXPECT().
			Put(gomock.Any()).
			Do(func(n notification.Notification) {
				actualEventsTitles = append(actualEventsTitles, n.EventTitle)
			}).
			Return(nil).
			Times(len(expectedEventsTitles))
		producer.EXPECT().Close().Return(nil)
		storage := mock.NewMockStorage(controller)
		storage.EXPECT().
			GetEventsNotifyBetween(gomock.Any(), gomock.Any(), gomock.Any()).
			DoAndReturn(func(_ context.Context, _ time.Time, _ time.Time) ([]event.Event, error) {
				defer func() {
					callNumber++
				}()
				if callNumber < len(eventsToReturn) {
					return eventsToReturn[callNumber], nil
				}
				return make([]event.Event, 0), nil
			}).
			AnyTimes()
		s, err := NewScheduler(storage, producer)
		require.NoError(t, err)
		s.(*scheduler).scanInterval = time.Millisecond * 10

		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*45)
		defer cancel()
		err = s.FindNotificationReadyEvents(ctx)

		require.NoError(t, err)
		require.Equal(t, expectedEventsTitles, actualEventsTitles)
	})
}

func Test_scheduler_CleanOldEvents(t *testing.T) {
	errStorage := event.NewErrStorage("storage error")
	tests := map[string]struct {
		getStorage     func(controller *gomock.Controller) event.Storage
		contextTimeout time.Duration
		cleanTimeout   time.Duration
		err            error
	}{
		"when storage error occurred then error should be returned": {
			getStorage: func(controller *gomock.Controller) event.Storage {
				s := mock.NewMockStorage(controller)
				s.EXPECT().
					DeleteEventsOlderThan(gomock.Any(), gomock.Any()).
					Return(errStorage)
				return s
			},
			contextTimeout: time.Millisecond * 100,
			cleanTimeout:   time.Millisecond * 20,
			err:            errStorage,
		},
		"when context timed out then no calls should be done": {
			getStorage: func(controller *gomock.Controller) event.Storage {
				return mock.NewMockStorage(controller)
			},
			contextTimeout: time.Millisecond * 20,
			cleanTimeout:   time.Millisecond * 100,
			err:            nil,
		},
	}
	for testName, tt := range tests {
		t.Run(testName, func(t *testing.T) {
			controller := gomock.NewController(t)
			s, err := NewScheduler(tt.getStorage(controller), queuemock.NewMockProducer(controller))
			require.NoError(t, err)
			s.(*scheduler).cleanInterval = tt.cleanTimeout
			ctx, cancel := context.WithTimeout(context.Background(), tt.contextTimeout)
			defer cancel()

			err = s.CleanOldEvents(ctx, time.Hour)

			if tt.err != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tt.err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
