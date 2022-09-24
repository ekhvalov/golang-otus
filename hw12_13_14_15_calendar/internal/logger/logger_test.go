package logger

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLogger(t *testing.T) {
	type testCase struct {
		loggerLevel  string
		messageLevel string
		shouldPrint  bool
	}
	testCases := []testCase{
		// Debug
		{
			loggerLevel:  levelDebug,
			messageLevel: levelDebug,
			shouldPrint:  true,
		},
		{
			loggerLevel:  levelDebug,
			messageLevel: levelInfo,
			shouldPrint:  true,
		},
		{
			loggerLevel:  levelDebug,
			messageLevel: levelWarning,
			shouldPrint:  true,
		},
		{
			loggerLevel:  levelDebug,
			messageLevel: levelError,
			shouldPrint:  true,
		},
		// Info
		{
			loggerLevel:  levelInfo,
			messageLevel: levelDebug,
			shouldPrint:  false,
		},
		{
			loggerLevel:  levelInfo,
			messageLevel: levelInfo,
			shouldPrint:  true,
		},
		{
			loggerLevel:  levelInfo,
			messageLevel: levelWarning,
			shouldPrint:  true,
		},
		{
			loggerLevel:  levelInfo,
			messageLevel: levelError,
			shouldPrint:  true,
		},
		// Warning
		{
			loggerLevel:  levelWarning,
			messageLevel: levelDebug,
			shouldPrint:  false,
		},
		{
			loggerLevel:  levelWarning,
			messageLevel: levelInfo,
			shouldPrint:  false,
		},
		{
			loggerLevel:  "warn",
			messageLevel: levelWarning,
			shouldPrint:  true,
		},
		{
			loggerLevel:  levelWarning,
			messageLevel: levelError,
			shouldPrint:  true,
		},
		// Error
		{
			loggerLevel:  levelError,
			messageLevel: levelDebug,
			shouldPrint:  false,
		},
		{
			loggerLevel:  levelError,
			messageLevel: levelInfo,
			shouldPrint:  false,
		},
		{
			loggerLevel:  levelError,
			messageLevel: levelWarning,
			shouldPrint:  false,
		},
		{
			loggerLevel:  levelWarning,
			messageLevel: levelError,
			shouldPrint:  true,
		},
	}
	for _, tc := range testCases {
		var msg string
		if tc.shouldPrint {
			msg = fmt.Sprintf("%s logger should print %s message", tc.loggerLevel, tc.messageLevel)
		} else {
			msg = fmt.Sprintf("%s logger should not print %s message", tc.loggerLevel, tc.messageLevel)
		}
		t.Run(msg, func(t *testing.T) {
			w := &bytes.Buffer{}
			l := New(tc.loggerLevel, w)

			switch tc.messageLevel {
			case levelDebug:
				l.Debug(msg)
			case levelInfo:
				l.Info(msg)
			case levelWarning:
				l.Warn(msg)
			case levelError:
				l.Error(msg)
			default:
				require.Failf(t, "undefined message level: %s", tc.messageLevel)
			}

			if tc.shouldPrint {
				require.Contains(t, w.String(), msg)
			} else {
				require.NotContains(t, w.String(), msg)
			}
		})
	}
}
