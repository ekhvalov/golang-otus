package scheduler

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/ekhvalov/hw12_13_14_15_calendar/internal/app/notification/queue"
	"github.com/ekhvalov/hw12_13_14_15_calendar/internal/domain/event"
	"github.com/ekhvalov/hw12_13_14_15_calendar/internal/domain/notification"
)

// ErrSchedule contains a list of schedule errors.
type ErrSchedule struct {
	Errors []error
}

func (e *ErrSchedule) Error() string {
	return strings.Join(errorsToStrings(e.Errors), "; ")
}

func (e *ErrSchedule) Add(err error) {
	e.Errors = append(e.Errors, err)
}

func errorsToStrings(errors []error) []string {
	s := make([]string, len(errors))
	for i, err := range errors {
		s[i] = err.Error()
	}
	return s
}

type Scheduler interface {
	FindNotificationReadyEvents(ctx context.Context, interval time.Duration) error
	CleanOldEvents(ctx context.Context, outDatePeriod, cleanInterval time.Duration) error
}

func NewScheduler(storage event.Storage, producer queue.Producer) (Scheduler, error) {
	if storage == nil {
		return nil, fmt.Errorf("required Storage, but <nil> provided")
	}
	if producer == nil {
		return nil, fmt.Errorf("required Producer, but <nil> provided")
	}
	return &scheduler{
		storage:       storage,
		producer:      producer,
		errors:        ErrSchedule{Errors: make([]error, 0)},
		cleanInterval: time.Hour,
		scanInterval:  time.Minute,
	}, nil
}

type scheduler struct {
	storage       event.Storage
	producer      queue.Producer
	errors        ErrSchedule
	cleanInterval time.Duration // Interval of removing old events from storage
	scanInterval  time.Duration // Interval of finding events that are ready to send notification about
}

func (s *scheduler) FindNotificationReadyEvents(ctx context.Context, interval time.Duration) error {
	t := time.NewTicker(interval)
	defer func() {
		t.Stop()
	}()
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			select {
			case <-ctx.Done():
				return nil
			case <-t.C:
				events, err := s.storage.GetEventsNotifyBetween(ctx, time.Now(), time.Now().Add(interval))
				if err != nil {
					return fmt.Errorf("get events error: %w", err)
				}
				for _, e := range events {
					err = s.producer.Put(notification.Notification{
						EventID:    e.ID,
						EventTitle: e.Title,
						EventDate:  e.DateTime,
						UserID:     e.UserID,
					})
					if err != nil {
						return fmt.Errorf("error while put notification into a queue: %w", err)
					}
				}
			}
		}
	}
}

func (s *scheduler) CleanOldEvents(ctx context.Context, outDatePeriod, cleanInterval time.Duration) error {
	t := time.NewTicker(cleanInterval)
	defer func() {
		t.Stop()
	}()
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			select {
			case <-ctx.Done():
				return nil
			case <-t.C:
				outDateTime := time.Now().Unix() - int64(outDatePeriod.Seconds())
				if err := s.storage.DeleteEventsOlderThan(ctx, time.Unix(outDateTime, 0)); err != nil {
					return fmt.Errorf("delete old events error: %w", err)
				}
			}
		}
	}
}
