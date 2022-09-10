package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	type testCase struct {
		cmd      []string
		env      Environment
		exitCode int
	}

	testCases := map[string]testCase{
		"direct exit code": {
			cmd:      []string{"sh", "-c", "exit 10"},
			env:      Environment{},
			exitCode: 10,
		},
		"exit code from environment variable": {
			cmd:      []string{"sh", "-c", "exit $EXIT_CODE"},
			env:      Environment{"EXIT_CODE": EnvValue{Value: "101"}},
			exitCode: 101,
		},
		"plain exit": {
			cmd:      []string{"sh"},
			env:      Environment{},
			exitCode: 0,
		},
	}
	for testName, tc := range testCases {
		t.Run(testName, func(t *testing.T) {
			actualExitCode := RunCmd(tc.cmd, tc.env)
			require.Equal(t, tc.exitCode, actualExitCode)
		})
	}
}
