package validators

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

// nolint:funlen
func TestValidatorMin(t *testing.T) {
	type args struct {
		args  string
		value interface{}
	}
	tests := []struct {
		name            string
		args            args
		wantCreateErr   bool
		wantValidateErr []error
	}{
		{
			name:            "unsupported type",
			args:            args{args: "10", value: func() {}},
			wantCreateErr:   true,
			wantValidateErr: nil,
		},
		{
			name:            "nil",
			args:            args{args: "10", value: nil},
			wantCreateErr:   true,
			wantValidateErr: nil,
		},
		{
			name:            "integer 10 < 10",
			args:            args{args: "10", value: 10},
			wantCreateErr:   false,
			wantValidateErr: nil,
		},
		{
			name:            "negative integer -10 < 10",
			args:            args{args: "-10", value: 10},
			wantCreateErr:   false,
			wantValidateErr: nil,
		},
		{
			name:            "negative integer -15 < 10",
			args:            args{args: "-15", value: -10},
			wantCreateErr:   false,
			wantValidateErr: nil,
		},
		{
			name:            "integer x < 10 (create error)",
			args:            args{args: "x", value: 10},
			wantCreateErr:   true,
			wantValidateErr: nil,
		},
		{
			name:            "integer 20 < 10 (validate error)",
			args:            args{args: "20", value: 10},
			wantCreateErr:   false,
			wantValidateErr: []error{fmt.Errorf("10 is less than 20")},
		},
		{
			name:            "negative integer 20 < -10 (validate error)",
			args:            args{args: "20", value: -10},
			wantCreateErr:   false,
			wantValidateErr: []error{fmt.Errorf("-10 is less than 20")},
		},
		{
			name:            "int8",
			args:            args{args: "10", value: int8(11)},
			wantCreateErr:   false,
			wantValidateErr: nil,
		},
		{
			name:            "int8 (validate error)",
			args:            args{args: "20", value: int8(-10)},
			wantCreateErr:   false,
			wantValidateErr: []error{fmt.Errorf("-10 is less than 20")},
		},
		{
			name:          "[]int16 (validate error)",
			args:          args{args: "20", value: []int16{int16(10), int16(19)}},
			wantCreateErr: false,
			wantValidateErr: []error{
				fmt.Errorf("10 is less than 20"),
				fmt.Errorf("19 is less than 20"),
			},
		},
		{
			name:            "[]int32",
			args:            args{args: "10", value: []int32{int32(10), int32(19)}},
			wantCreateErr:   false,
			wantValidateErr: nil,
		},
		{
			name:            "uint (create error)",
			args:            args{args: "-10", value: uint(10)},
			wantCreateErr:   true,
			wantValidateErr: nil,
		},
		{
			name:            "uint8 (validate error)",
			args:            args{args: "20", value: uint8(10)},
			wantCreateErr:   false,
			wantValidateErr: []error{fmt.Errorf("10 is less than 20")},
		},
		{
			name:            "uint16",
			args:            args{args: "10", value: uint16(10)},
			wantCreateErr:   false,
			wantValidateErr: nil,
		},
		{
			name:            "[]uint32",
			args:            args{args: "10", value: []uint32{uint32(10), uint32(15)}},
			wantCreateErr:   false,
			wantValidateErr: nil,
		},
		{
			name:            "float32 (create error)",
			args:            args{args: "d01.5", value: float32(10.1)},
			wantCreateErr:   true,
			wantValidateErr: nil,
		},
		{
			name:            "float32 (validate error)",
			args:            args{args: "10.5", value: float32(10.1)},
			wantCreateErr:   false,
			wantValidateErr: []error{fmt.Errorf("10.100000 is less than 10.500000")},
		},
		{
			name:            "float64",
			args:            args{args: "10.5", value: 10.6},
			wantCreateErr:   false,
			wantValidateErr: nil,
		},
		{
			name:            "string (create error)",
			args:            args{args: "x", value: "hello world"},
			wantCreateErr:   true,
			wantValidateErr: nil,
		},
		{
			name:            "string (validate error)",
			args:            args{args: "20", value: "hello world"},
			wantCreateErr:   false,
			wantValidateErr: []error{fmt.Errorf("11 is less than 20")},
		},
		{
			name:            "string utf8 (validate error)",
			args:            args{args: "20", value: "こんいちは"},
			wantCreateErr:   false,
			wantValidateErr: []error{fmt.Errorf("5 is less than 20")},
		},
		{
			name:            "string utf8",
			args:            args{args: "3", value: []string{"こんいちは", "こんばんは"}},
			wantCreateErr:   false,
			wantValidateErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v, err := newValidatorMin(tt.args.args, tt.args.value)
			if tt.wantCreateErr {
				require.NotNil(t, err)
				return
			}
			require.Nil(t, err)
			require.NotNil(t, v)
			errors := v.Validate()
			require.Equalf(t, tt.wantValidateErr, errors, "")
			if tt.wantValidateErr == nil {
				require.Nilf(t, errors, "unexpected errors: %v", errors)
			} else {
				require.NotNil(t, errors)
			}
		})
	}
}
