package validators

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type likeString string

func Test_newValidatorIn(t *testing.T) {
	tests := []struct {
		name            string
		args            args
		wantCreateErr   error
		wantValidateErr []error
	}{
		{
			name:            "unsupported type",
			args:            args{args: "10", value: func() {}},
			wantCreateErr:   fmt.Errorf("unsupported type: func()"),
			wantValidateErr: nil,
		},
		{
			name:            "invalid argument",
			args:            args{args: "", value: []int{10, 15}},
			wantCreateErr:   fmt.Errorf("argument parse error: strconv.ParseInt: parsing \"\": invalid syntax"),
			wantValidateErr: nil,
		},
		{
			name:          "validate errors (int)",
			args:          args{args: "20,25", value: []int{10, 15}},
			wantCreateErr: nil,
			wantValidateErr: []error{
				fmt.Errorf("unexpected value: 10"),
				fmt.Errorf("unexpected value: 15"),
			},
		},
		{
			name:          "validate errors (uint)",
			args:          args{args: "20,25", value: []uint{10, 15}},
			wantCreateErr: nil,
			wantValidateErr: []error{
				fmt.Errorf("unexpected value: 10"),
				fmt.Errorf("unexpected value: 15"),
			},
		},
		{
			name:          "validate errors (string)",
			args:          args{args: "20,25", value: []string{"10", "15"}},
			wantCreateErr: nil,
			wantValidateErr: []error{
				fmt.Errorf("unexpected value: 10"),
				fmt.Errorf("unexpected value: 15"),
			},
		},
		{
			name:            "no errors (string)",
			args:            args{args: "10,15", value: []string{"10", "15"}},
			wantCreateErr:   nil,
			wantValidateErr: nil,
		},
		{
			name:            "no errors (int)",
			args:            args{args: "10,15", value: []int{10, 15}},
			wantCreateErr:   nil,
			wantValidateErr: nil,
		},
		{
			name:            "no errors (uint)",
			args:            args{args: "10,15", value: []uint{10, 15}},
			wantCreateErr:   nil,
			wantValidateErr: nil,
		},
		{
			name:            "no errors (likeString)",
			args:            args{args: "10,15", value: []likeString{"10", "15"}},
			wantCreateErr:   nil,
			wantValidateErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v, createErr := newValidatorIn(tt.args.args, tt.args.value)
			if tt.wantCreateErr != nil {
				require.NotNil(t, createErr)
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

type args struct {
	args  string
	value interface{}
}
