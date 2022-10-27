package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		meta   json.RawMessage
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in:          nil,
			expectedErr: fmt.Errorf("wrong validated type; expected struct, got: <nil>"),
		},
		{
			in:          10,
			expectedErr: fmt.Errorf("wrong validated type; expected struct, got: int"),
		},
		{
			in:          make([]string, 0),
			expectedErr: fmt.Errorf("wrong validated type; expected struct, got: []string"),
		},
		{
			in:          Token{},
			expectedErr: nil,
		},
		{
			in: User{
				ID:     "012345678901234567890123456789012345",
				Age:    28,
				Role:   "admin",
				Email:  "mail@example.com",
				Phones: []string{"+0123456789", "+0123456789"},
				meta:   nil,
			},
			expectedErr: nil,
		},
		{
			in: User{Age: 100},
			expectedErr: ValidationErrors{
				{
					Field: "ID",
					Err:   fmt.Errorf("length mismatched; expected 36, got 0"),
				},
			},
		},
		{
			in:          App{Version: "0.1.1"},
			expectedErr: nil,
		},
		{
			in:          Response{Code: 200},
			expectedErr: nil,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)
			if tt.expectedErr == nil {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				errors.Is(tt.expectedErr, err)
			}
			_ = tt
		})
	}
}
