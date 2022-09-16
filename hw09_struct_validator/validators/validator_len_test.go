package validators

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_newValidatorLen(t *testing.T) {
	type args struct {
		args  string
		value interface{}
	}
	tests := []struct {
		name            string
		args            args
		wantCreateErr   error
		wantValidateErr []error
	}{
		{
			name:            "unsupported type",
			args:            args{args: "10", value: 10},
			wantCreateErr:   fmt.Errorf("unsupported type: int"),
			wantValidateErr: nil,
		},
		{
			name:            "invalid argument",
			args:            args{args: "-10", value: "10"},
			wantCreateErr:   fmt.Errorf("argument parse error: strconv.ParseUint: parsing \"-10\": invalid syntax"),
			wantValidateErr: nil,
		},
		{
			name:            "validate error",
			args:            args{args: "10", value: "10"},
			wantCreateErr:   nil,
			wantValidateErr: []error{fmt.Errorf("length mismatched; expected 10, got 2")},
		},
		{
			name:            "no error",
			args:            args{args: "10", value: "0123456789"},
			wantCreateErr:   nil,
			wantValidateErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v, createErr := newValidatorLen(tt.args.args, tt.args.value)
			if tt.wantCreateErr != nil {
				require.NotNil(t, createErr)
				require.Equal(t, createErr, tt.wantCreateErr)
				return
			}
			require.Nil(t, createErr)
			require.NotNil(t, v)
			errors := v.Validate()
			if tt.wantValidateErr == nil {
				require.Nilf(t, errors, "unexpected validate error(s): %v", errors)
			} else {
				require.NotNilf(t, errors, "expected validate error(s) not found")
				require.Equal(t, tt.wantValidateErr, errors)
			}
		})
	}
}
