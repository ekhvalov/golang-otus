package pgsqlstorage

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/ekhvalov/hw12_13_14_15_calendar/internal/domain/event"
	"github.com/jackc/pgx/v5"
)

type Storage struct {
	dsn  string
	conn *pgx.Conn
}

func NewStorage(conf Config) *Storage {
	return &Storage{dsn: conf.GetDSN()}
}

func (s *Storage) Connect(ctx context.Context) error {
	conn, err := pgx.Connect(ctx, s.dsn)
	if err != nil {
		return fmt.Errorf("connection error: %w", err)
	}
	s.conn = conn
	return nil
}

func (s *Storage) Close(ctx context.Context) error {
	if s.conn != nil {
		return s.conn.Close(ctx)
	}
	return nil
}

func (s *Storage) Create(ctx context.Context, e event.Event) (event.Event, error) {
	if s.conn == nil {
		err := s.Connect(ctx)
		if err != nil {
			return event.Event{}, fmt.Errorf("database connection error: %w", err)
		}
	}
	var lastID uint64
	err := s.conn.QueryRow(
		ctx,
		`INSERT INTO events (title, description, start_time, duration, user_id) VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		e.Title,
		e.Description,
		e.DateTime,
		e.Duration,
		e.UserID,
	).Scan(&lastID)
	if err != nil {
		return event.Event{}, fmt.Errorf("save event error: %w", err)
	}
	e.ID = strconv.Itoa(int(lastID))
	return e, nil
}

func (s *Storage) Update(ctx context.Context, eventID string, e event.Event) error {
	if s.conn == nil {
		err := s.Connect(ctx)
		if err != nil {
			return fmt.Errorf("database connection error: %w", err)
		}
	}
	_, err := s.conn.Query(
		ctx,
		`UPDATE events set title=$1, description=$2, start_time=$3, duration=$4, user_id=$5 WHERE id=$6`,
		e.Title,
		e.Description,
		e.DateTime,
		e.Duration,
		e.UserID,
		eventID,
	)
	if err != nil {
		return fmt.Errorf("update event error: %w", err)
	}
	return nil
}

func (s *Storage) Delete(ctx context.Context, eventID string) error {
	if s.conn == nil {
		err := s.Connect(ctx)
		if err != nil {
			return fmt.Errorf("database connection error: %w", err)
		}
	}
	_, err := s.conn.Query(ctx, `DELETE FROM events WHERE id=$1`, eventID)
	if err != nil {
		return fmt.Errorf("delete event error: %w", err)
	}
	return nil
}

func (s *Storage) GetDayEvents(ctx context.Context, date time.Time) ([]event.Event, error) {
	start := date.Unix()
	end := date.AddDate(0, 0, 1).Unix()
	return s.getEventsByTimeRange(ctx, start, end)
}

func (s *Storage) GetWeekEvents(ctx context.Context, date time.Time) ([]event.Event, error) {
	start := date.Unix()
	end := date.AddDate(0, 0, 7).Unix()
	return s.getEventsByTimeRange(ctx, start, end)
}

func (s *Storage) GetMonthEvents(ctx context.Context, date time.Time) ([]event.Event, error) {
	start := date.Unix()
	end := date.AddDate(0, 1, 0).Unix()
	return s.getEventsByTimeRange(ctx, start, end)
}

func (s *Storage) getEventsByTimeRange(ctx context.Context, start, end int64) ([]event.Event, error) {
	if s.conn == nil {
		err := s.Connect(ctx)
		if err != nil {
			return nil, fmt.Errorf("database connection error: %w", err)
		}
	}
	result, err := s.conn.Query(
		ctx,
		`
SELECT
    id,
    user_id,
    title,
    description,
    start_time,
    duration
FROM events WHERE start_time >= $1 and start_time <= $2`,
		start,
		end,
	)
	if err != nil {
		return nil, fmt.Errorf("get events error: %w", err)
	}
	events := make([]event.Event, 0)
	for result.Next() {
		var date int64
		e := event.Event{}
		err = result.Scan(&e.ID, &e.UserID, &e.Title, &e.Description, &date, &e.Duration)
		if err != nil {
			return nil, fmt.Errorf("retrieve events error: %w", err)
		}
		e.DateTime = time.Unix(date, 0)
		events = append(events, e)
	}
	return events, nil
}
