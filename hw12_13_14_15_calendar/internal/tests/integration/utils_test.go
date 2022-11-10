//go:build integration

package integration_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/ekhvalov/golang-otus/hw12_13_14_15_calendar/internal/domain/event"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
)

const (
	defaultHttpServerHost   = "localhost"
	defaultHttpServerPort   = "8080"
	defaultDatabaseHost     = "localhost"
	defaultDatabasePort     = "5432"
	defaultDatabaseUsername = "postgres"
	defaultDatabasePassword = "password"
	defaultDatabaseName     = "postgres"
)

func seedEventsIntoDB(t *testing.T, ctx context.Context, db *pgx.Conn, events []event.Event) {
	sql := `
INSERT INTO events (user_id, title, description, start_time, end_time, notify_time) 
VALUES ($1, $2, $3, $4, $5, $6)`
	stmt, err := db.Prepare(ctx, "insert event", sql)
	require.NoError(t, err)
	for _, e := range events {
		notifyTime := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
		if e.NotifyBefore > 0 {
			notifyTime = time.Unix(e.DateTime.Unix()-int64(e.NotifyBefore.Seconds()), 0).UTC()
		}
		if e.UserID == "" {
			e.UserID = "10"
		}
		_, err = db.Exec(ctx, stmt.Name, e.UserID, e.Title, e.Description, e.DateTime, e.DateTime.Add(e.Duration), notifyTime)
		require.NoError(t, err)
	}
	require.NoError(t, db.Deallocate(ctx, stmt.Name))
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

func getDatabaseAddress() string {
	host := getEnv("TESTS_DATABASE_HOST", defaultDatabaseHost)
	port := getEnv("TESTS_DATABASE_PORT", defaultDatabasePort)
	username := getEnv("TESTS_DATABASE_USERNAME", defaultDatabaseUsername)
	password := getEnv("TESTS_DATABASE_PASSWORD", defaultDatabasePassword)
	name := getEnv("TESTS_DATABASE_NAME", defaultDatabaseName)
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s", username, password, host, port, name)
}

func getEnv(name, defaultValue string) string {
	value := os.Getenv(name)
	if value == "" {
		return defaultValue
	}
	return value
}
