package pgsqlstorage

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/ekhvalov/golang-otus/hw12_13_14_15_calendar/internal/domain/event"
	"github.com/jackc/pgx/v5"
)

type Storage struct {
	dsn  string
	conn *pgx.Conn
}

func NewStorage(conf Config) event.Storage {
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
	startTime := e.DateTime.UTC()
	endTime := e.DateTime.Add(e.Duration).UTC()
	isTimeAvailable, err := s.isTimeAvailable(ctx, startTime, endTime)
	if err != nil {
		return event.Event{}, fmt.Errorf("check time availability event error: %w", err)
	}
	if !isTimeAvailable {
		return event.Event{}, event.ErrDateBusy
	}
	var lastID uint64
	err = s.conn.QueryRow(
		ctx,
		`INSERT INTO events (title, description, start_time, end_time, user_id) VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		e.Title,
		e.Description,
		startTime,
		endTime,
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
	startTime := e.DateTime.UTC()
	endTime := e.DateTime.Add(e.Duration).UTC()
	isTimeAvailable, err := s.isTimeAvailable(ctx, startTime, endTime)
	if err != nil {
		return fmt.Errorf("check time availability event error: %w", err)
	}
	if !isTimeAvailable {
		return event.ErrDateBusy
	}
	_, err = s.conn.Query(
		ctx,
		`UPDATE events set title=$1, description=$2, start_time=$3, end_time=$4, user_id=$5 WHERE id=$6`,
		e.Title,
		e.Description,
		startTime,
		endTime,
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
	end := date.AddDate(0, 0, 1)
	return s.getEventsByTimeRange(ctx, date, end)
}

func (s *Storage) GetWeekEvents(ctx context.Context, date time.Time) ([]event.Event, error) {
	end := date.AddDate(0, 0, 7)
	return s.getEventsByTimeRange(ctx, date, end)
}

func (s *Storage) GetMonthEvents(ctx context.Context, date time.Time) ([]event.Event, error) {
	end := date.AddDate(0, 1, 0)
	return s.getEventsByTimeRange(ctx, date, end)
}

func (s *Storage) GetEventsNotifyBetween(ctx context.Context, from time.Time, to time.Time) ([]event.Event, error) {
	return s.getEventsByTimeRange(ctx, from, to)
}

func (s *Storage) DeleteEventsOlderThan(ctx context.Context, date time.Time) error {
	if s.conn == nil {
		err := s.Connect(ctx)
		if err != nil {
			return fmt.Errorf("database connection error: %w", err)
		}
	}
	sql := `DELETE FROM events WHERE start_time < $1`
	_, err := s.conn.Exec(ctx, sql, date.Unix())
	return err
}

func (s *Storage) getEventsByTimeRange(ctx context.Context, start, end time.Time) ([]event.Event, error) {
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
    end_time
FROM events WHERE start_time >= $1 and start_time < $2`,
		start.UTC(),
		end.UTC(),
	)
	if err != nil {
		return nil, fmt.Errorf("get events error: %w", err)
	}
	events := make([]event.Event, 0)
	for result.Next() {
		var startTime, endTime time.Time
		e := event.Event{}
		err = result.Scan(&e.ID, &e.UserID, &e.Title, &e.Description, &startTime, &endTime)
		if err != nil {
			return nil, fmt.Errorf("retrieve events error: %w", err)
		}
		e.DateTime = startTime
		e.Duration = startTime.Sub(endTime)
		events = append(events, e)
	}
	return events, nil
}

func (s *Storage) isTimeAvailable(ctx context.Context, startTime, endTime time.Time) (bool, error) {
	sql := `
SELECT COUNT(*) FROM events WHERE (
    (start_time >= $1 AND start_time < $2)
    OR 
    (end_time > $1 AND end_time <= $2)
    OR
    ($1 >= start_time AND $1 < end_time)
    OR 
    ($2 > start_time AND $2 <= end_time)
) LIMIT 1;`
	result := s.conn.QueryRow(ctx, sql, startTime, endTime)
	var count int
	if err := result.Scan(&count); err != nil {
		return false, err
	}
	return count == 0, nil
}
