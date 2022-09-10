package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	type testCase struct {
		path        string
		expectedEnv Environment
	}
	testCases := map[string]testCase{
		"when files existed then should read variables": {
			path: "./testdata/env",
			expectedEnv: Environment{
				"BAR":   {Value: "bar"},
				"EMPTY": {Value: ""},
				"FOO":   {Value: "   foo\nwith new line"},
				"HELLO": {Value: "\"hello\""},
				"UNSET": {Value: "", NeedRemove: true},
			},
		},
		"when read empty dir then environment should be empty": {
			path:        "./testdata/env-empty",
			expectedEnv: Environment{},
		},
	}
	for testName, tc := range testCases {
		t.Run(testName, func(t *testing.T) {
			actualEnv, err := ReadDir(tc.path)
			require.NoError(t, err)
			require.NotNil(t, actualEnv)
			for key, value := range tc.expectedEnv {
				actualValue, ok := actualEnv[key]
				require.Truef(t, ok, "no value for key: %s", key)
				require.Equalf(t, value.Value, actualValue.Value, "wrong value for key: %s", key)
			}
			for key := range actualEnv {
				_, ok := tc.expectedEnv[key]
				require.Truef(t, ok, "key %s is present but should not")
			}
		})
	}
}
