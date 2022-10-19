//go:build integration
// +build integration

package integrationt_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"testing"
	"time"

	"github.com/ekhvalov/hw12_13_14_15_calendar/pkg/api/openapi"
	"github.com/stretchr/testify/require"
)

const (
	defaultHttpServerHost = "localhost"
	defaultHttpServerPort = "8080"
)

func Test_CreateEvent_Http(t *testing.T) {
	client, err := openapi.NewClient(getHttpServerAddress())
	require.NoError(t, err)

	t.Run("create event 200 OK", func(t *testing.T) {
		event := openapi.NewEvent{
			Date:     time.Now().Add(time.Hour).Unix(),
			Duration: 20,
			Title:    "Event 1",
		}
		response, err := client.PostEvents(context.Background(), event)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, response.StatusCode)
	})

	t.Run("create event with empty title", func(t *testing.T) {
		event := openapi.NewEvent{
			Date:     time.Now().Add(time.Hour * 2).Unix(),
			Duration: 20,
			Title:    "",
		}
		response, err := client.PostEvents(context.Background(), event)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, response.StatusCode)
	})

	t.Run("create event busy date", func(t *testing.T) {
		event1 := openapi.NewEvent{
			Date:     time.Now().Add(time.Hour * 3).Unix(),
			Duration: 20,
			Title:    "Event 3",
		}
		response, err := client.PostEvents(context.Background(), event1)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, response.StatusCode)
		event2 := openapi.NewEvent{
			Date:     time.Now().Add(time.Hour * 3).Unix(),
			Duration: 20,
			Title:    "Event 4",
		}
		response, err = client.PostEvents(context.Background(), event2)
		require.NoError(t, err)
		require.Equal(t, http.StatusConflict, response.StatusCode)
	})
}

func Test_GetEvents(t *testing.T) {
	client, err := openapi.NewClient(getHttpServerAddress())
	require.NoError(t, err)
	now := time.Now()
	year, month, _ := now.AddDate(0, 1, 0).Date()
	date := time.Date(year, month, 1, 0, 0, 0, 0, now.Location()) // start of the next month
	event1 := openapi.NewEvent{
		Date:     date.Unix(),
		Duration: 20,
		Title:    "Event 1",
	}
	event2 := openapi.NewEvent{
		Date:     date.Add(time.Hour).Unix(), // one hour after the date
		Duration: 20,
		Title:    "Event 2",
	}
	event3 := openapi.NewEvent{
		Date:     date.AddDate(0, 0, 1).Unix(), // next day after the date
		Duration: 20,
		Title:    "Event 3",
	}
	event4 := openapi.NewEvent{
		Date:     date.AddDate(0, 0, 14).Unix(), // two weeks after the date
		Duration: 20,
		Title:    "Event 4",
	}
	for _, event := range []openapi.NewEvent{event1, event2, event3, event4} {
		response, err := client.PostEvents(context.Background(), event)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, response.StatusCode)
	}

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
		t.Run(testName, func(t *testing.T) {
			response, err := client.GetEvents(context.Background(), &openapi.GetEventsParams{
				Period: tt.period,
				Date:   date.Unix(),
			})
			require.NoError(t, err)
			events := make([]openapi.Event, 0)
			body, err := io.ReadAll(response.Body)
			require.NoError(t, err)
			err = json.Unmarshal(body, &events)
			require.NoError(t, err)
			eventTitles := make([]string, len(events))
			for i, event := range events {
				eventTitles[i] = event.Title
			}
			sort.Strings(eventTitles)
			require.Equal(t, tt.wantTitles, eventTitles)
		})
	}
}

func getHttpServerAddress() string {
	host := os.Getenv("TESTS_HTTP_SERVER_HOST")
	if host == "" {
		host = defaultHttpServerHost
	}
	port := os.Getenv("TESTS_HTTP_SERVER_PORT")
	if port == "" {
		port = defaultHttpServerPort
	}
	return fmt.Sprintf("http://%s:%s", host, port)
}
