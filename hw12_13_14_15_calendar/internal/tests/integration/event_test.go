//go:build integration

package integration_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"testing"
	"time"

	"github.com/ekhvalov/golang-otus/hw12_13_14_15_calendar/internal/domain/event"
	"github.com/ekhvalov/golang-otus/hw12_13_14_15_calendar/pkg/api/openapi"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func TestEvent(t *testing.T) {
	suite.Run(t, new(eventSuite))
}

type eventSuite struct {
	suite.Suite
	ctx        context.Context
	cancel     context.CancelFunc
	tick       time.Duration
	waitFor    time.Duration
	db         *pgx.Conn
	httpClient *openapi.Client
}

func (s *eventSuite) SetupSuite() {
	s.tick = time.Second
	s.waitFor = s.tick * 30
	s.ctx, s.cancel = context.WithCancel(context.Background())

	require.Eventually(s.T(), func() bool {
		connect, err := pgx.Connect(s.ctx, getDatabaseAddress())
		if err != nil {
			return false
		}
		s.db = connect
		return true
	}, s.waitFor, s.tick, "can not connect to database")

	require.Eventually(s.T(), func() bool {
		return s.db.Ping(s.ctx) == nil
	}, s.waitFor, s.tick, "can not ping to database")

	require.Eventually(s.T(), func() bool {
		sql := "SELECT EXISTS (SELECT FROM pg_tables WHERE schemaname = 'public' AND tablename = 'events');"
		result := s.db.QueryRow(s.ctx, sql)
		var isTableExists bool
		err := result.Scan(&isTableExists)
		require.NoError(s.T(), err)
		return isTableExists
	}, s.waitFor, s.tick, "migrations did not applied")

	require.Eventually(s.T(), func() bool {
		client, err := openapi.NewClient(getHttpServerAddress())
		if err != nil {
			return false
		}
		s.httpClient = client
		return true
	}, s.waitFor, s.tick, "can not connect to HTTP client")
}

func (s *eventSuite) TearDownSuite() {
	s.cleanDatabase()
	if s.db != nil {
		_ = s.db.Close(s.ctx)
	}
	s.cancel()
}

func (s *eventSuite) SetupTest() {
	s.cleanDatabase()
}

func (s *eventSuite) cleanDatabase() {
	_, err := s.db.Exec(s.ctx, "TRUNCATE TABLE events;")
	if err != nil {
		panic(fmt.Errorf("database clean error: %s", err))
	}
}

func (s *eventSuite) Test_CreateEvent_Http() {
	eventTime := time.Now().Add(time.Hour)
	e := openapi.NewEvent{
		Date:     eventTime.Unix(),
		Duration: 20,
		Title:    fmt.Sprintf("Event %d", eventTime.UnixMicro()),
	}
	response, err := s.httpClient.PostEvents(s.ctx, e)
	require.NoError(s.T(), err)
	require.Equal(s.T(), http.StatusOK, response.StatusCode)

	sql := `SELECT COUNT(*) FROM events WHERE title = $1`
	result := s.db.QueryRow(s.ctx, sql, e.Title)
	var count int
	err = result.Scan(&count)
	require.NoError(s.T(), err)
	require.Equalf(s.T(), 1, count, "event with title '%s' is not found", e.Title)
}

func (s *eventSuite) Test_CreateEvent_Http_Error_InvalidEvent() {
	tests := map[string]openapi.NewEvent{
		"empty title": {
			Title:    "",
			Date:     time.Now().Add(time.Hour).Unix(),
			Duration: 20,
		},
		"date is in the past": {
			Title:    "Event",
			Date:     time.Now().Unix() - 3600,
			Duration: 0,
		},
		"empty duration": {
			Title:    "Event",
			Date:     time.Now().Add(time.Hour).Unix(),
			Duration: 0,
		},
	}
	for testName, e := range tests {
		s.T().Run(testName, func(t *testing.T) {
			response, err := s.httpClient.PostEvents(s.ctx, e)
			require.NoError(s.T(), err)
			require.Equal(s.T(), http.StatusBadRequest, response.StatusCode)
		})
	}
}

func (s *eventSuite) Test_CreateEvent_Http_Error_DateBusy() {
	year, month, day := time.Now().AddDate(0, 0, 2).Date()
	startDate := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	events := []event.Event{
		{
			Title:    "Event 10:00 to 11:30",
			DateTime: startDate.Add(time.Hour * 10),
			Duration: time.Minute * 90,
		},
		{
			Title:    "Event 12:30 to 15:30",
			DateTime: startDate.Add(time.Hour*12 + time.Minute*30),
			Duration: time.Hour * 3,
		},
	}
	seedEventsIntoDB(s.T(), s.ctx, s.db, events)

	tests := []struct {
		testName     string
		date         time.Time
		duration     time.Duration
		responseCode int
	}{
		{
			testName:     "11:30 to 12:31 Err",
			date:         time.Date(year, month, day, 11, 30, 0, 0, time.UTC),
			duration:     time.Minute * 61,
			responseCode: http.StatusConflict,
		},
		{
			testName:     "09:00 to 10:01 Err",
			date:         time.Date(year, month, day, 9, 0, 0, 0, time.UTC),
			duration:     time.Minute * 61,
			responseCode: http.StatusConflict,
		},
		{
			testName:     "09:00 to 12:00 Err",
			date:         time.Date(year, month, day, 9, 0, 0, 0, time.UTC),
			duration:     time.Hour * 3,
			responseCode: http.StatusConflict,
		},
		{
			testName:     "09:00 to 16:00 Err",
			date:         time.Date(year, month, day, 9, 0, 0, 0, time.UTC),
			duration:     time.Hour * 7,
			responseCode: http.StatusConflict,
		},
		{
			testName:     "10:00 to 11:30 Err",
			date:         time.Date(year, month, day, 10, 0, 0, 0, time.UTC),
			duration:     time.Minute * 90,
			responseCode: http.StatusConflict,
		},
		{
			testName:     "10:30 to 11:00 Err",
			date:         time.Date(year, month, day, 10, 30, 0, 0, time.UTC),
			duration:     time.Minute * 30,
			responseCode: http.StatusConflict,
		},
		{
			testName:     "10:30 to 15:00 Err",
			date:         time.Date(year, month, day, 10, 30, 0, 0, time.UTC),
			duration:     time.Hour*4 + time.Minute*30,
			responseCode: http.StatusConflict,
		},
		{
			testName:     "15:00 to 16:00 Err",
			date:         time.Date(year, month, day, 15, 0, 0, 0, time.UTC),
			duration:     time.Hour,
			responseCode: http.StatusConflict,
		},
		{
			testName:     "09:00 to 10:00 OK",
			date:         time.Date(year, month, day, 9, 0, 0, 0, time.UTC),
			duration:     time.Hour,
			responseCode: http.StatusOK,
		},
		{
			testName:     "11:30 to 12:30 OK",
			date:         time.Date(year, month, day, 11, 30, 0, 0, time.UTC),
			duration:     time.Hour,
			responseCode: http.StatusOK,
		},
		{
			testName:     "15:30 to 16:30 OK",
			date:         time.Date(year, month, day, 15, 30, 0, 0, time.UTC),
			duration:     time.Hour,
			responseCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.testName, func(t *testing.T) {
			response, err := s.httpClient.PostEvents(s.ctx, openapi.NewEvent{
				Title:    tt.testName,
				Date:     tt.date.Unix(),
				Duration: int(tt.duration / time.Minute),
			})
			require.NoError(t, err)
			require.Equal(t, tt.responseCode, response.StatusCode)
		})
	}
}

func (s *eventSuite) Test_GetEvents_Http() {
	now := time.Now().UTC()
	year, month, _ := now.AddDate(0, 1, 0).Date()
	date := time.Date(year, month, 1, 0, 0, 0, 0, now.Location()) // start of the next month
	event1 := event.Event{
		Title:    "Event 1",
		DateTime: date,
		Duration: time.Minute * 20,
	}
	event2 := event.Event{
		Title:    "Event 2",
		DateTime: date.Add(time.Hour),
		Duration: time.Minute * 20,
	}
	event3 := event.Event{
		Title:    "Event 3",
		DateTime: date.AddDate(0, 0, 1),
		Duration: time.Minute * 20,
	}
	event4 := event.Event{
		Title:    "Event 4",
		DateTime: date.AddDate(0, 0, 14),
		Duration: time.Minute * 20,
	}
	events := []event.Event{event1, event2, event3, event4}
	seedEventsIntoDB(s.T(), s.ctx, s.db, events)

	tests := map[string]struct {
		period     openapi.EventsPeriod
		wantTitles []string
	}{
		"day events": {
			period:     openapi.Day,
			wantTitles: []string{event1.Title, event2.Title},
		},
		"week events": {
			period:     openapi.Week,
			wantTitles: []string{event1.Title, event2.Title, event3.Title},
		},
		"month events": {
			period:     openapi.Month,
			wantTitles: []string{event1.Title, event2.Title, event3.Title, event4.Title},
		},
	}
	for testName, tt := range tests {
		s.T().Run(testName, func(t *testing.T) {
			response, err := s.httpClient.GetEvents(s.ctx, &openapi.GetEventsParams{
				Date:   date.Unix(),
				Period: tt.period,
			})
			require.NoError(t, err)
			require.Equal(t, http.StatusOK, response.StatusCode)
			actualEvents := make([]openapi.Event, 0)
			body, err := io.ReadAll(response.Body)
			require.NoError(t, err)
			err = json.Unmarshal(body, &actualEvents)
			require.NoError(t, err)
			actualEventTitles := make([]string, len(actualEvents))
			for i, e := range actualEvents {
				actualEventTitles[i] = e.Title
			}
			sort.Strings(actualEventTitles)
			require.ElementsMatch(t, tt.wantTitles, actualEventTitles)
		})
	}
}
