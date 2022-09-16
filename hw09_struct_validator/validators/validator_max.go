package validators

import (
	"fmt"
	"reflect"
	"strconv"
	"unicode/utf8"
)

func newValidatorMax(args string, value interface{}) (Validator, error) {
	builder := validatorBuilder{}
	builder.setFunctionsForKind(
		func(args string) (parsedArgs, error) {
			return strconv.ParseInt(args, 10, 0)
		},
		func(pArgs parsedArgs, value reflect.Value) error {
			threshold := pArgs.(int64)
			if value.Int() > threshold {
				return fmt.Errorf("%d is greater than %d", value.Int(), threshold)
			}
			return nil
		},
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
	)
	builder.setFunctionsForKind(
		func(args string) (parsedArgs, error) {
			return strconv.ParseUint(args, 10, 0)
		},
		func(pArgs parsedArgs, value reflect.Value) error {
			threshold := pArgs.(uint64)
			if value.Uint() > threshold {
				return fmt.Errorf("%d is greater than %d", value.Uint(), threshold)
			}
			return nil
		},
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
	)
	builder.setFunctionsForKind(
		func(args string) (parsedArgs, error) {
			return strconv.ParseFloat(args, 10)
		},
		func(pArgs parsedArgs, value reflect.Value) error {
			threshold := pArgs.(float64)
			if value.Float() > threshold {
				return fmt.Errorf("%f is greater than %f", value.Float(), threshold)
			}
			return nil
		},
		reflect.Float64, reflect.Float32,
	)
	builder.setFunctionsForKind(
		func(args string) (parsedArgs, error) {
			return strconv.ParseInt(args, 10, 0)
		},
		func(pArgs parsedArgs, value reflect.Value) error {
			threshold := pArgs.(int64)
			l := utf8.RuneCountInString(value.String())
			if int64(l) > threshold {
				return fmt.Errorf("%d is greater than %d", l, threshold)
			}
			return nil
		},
		reflect.String,
	)
	return builder.build(args, value)
}
