//go:build integration

package integration_test

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/ekhvalov/golang-otus/hw12_13_14_15_calendar/internal/domain/event"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const (
	defaultTargetFile = "/tmp/writer.txt"
)

func TestNotification(t *testing.T) {
	suite.Run(t, new(notificationSuite))
}

type notificationSuite struct {
	suite.Suite
	ctx            context.Context
	cancel         context.CancelFunc
	tick           time.Duration
	waitFor        time.Duration
	db             *pgx.Conn
	senderFileName string
}

func (s *notificationSuite) SetupSuite() {
	s.tick = time.Second
	s.waitFor = s.tick * 30
	s.ctx, s.cancel = context.WithCancel(context.Background())
	s.senderFileName = getEnv("TESTS_WRITER_TARGET_FILE", defaultTargetFile)

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
}

func (s *notificationSuite) TearDownSuite() {
	s.cleanDatabase()
	if s.db != nil {
		_ = s.db.Close(s.ctx)
	}
	s.cancel()
}

func (s *notificationSuite) SetupTest() {
	s.cleanDatabase()
}

func (s *notificationSuite) cleanDatabase() {
	_, err := s.db.Exec(s.ctx, "TRUNCATE TABLE events;")
	if err != nil {
		panic(fmt.Errorf("database clean error: %s", err))
	}
}

func (s *notificationSuite) Test_SendNotification() {
	e := event.Event{
		Title:        fmt.Sprintf("Event %d", time.Now().UnixMicro()),
		DateTime:     time.Now().Add(time.Minute + time.Second*10).UTC(),
		Duration:     time.Minute * 20,
		NotifyBefore: time.Minute,
	}
	seedEventsIntoDB(s.T(), s.ctx, s.db, []event.Event{e})

	require.Eventually(s.T(), func() bool {
		data, err := os.ReadFile(s.senderFileName)
		if err != nil {
			return false
		}
		return strings.Contains(string(data), e.Title)
	}, s.waitFor, s.tick, "can not read file")
}
