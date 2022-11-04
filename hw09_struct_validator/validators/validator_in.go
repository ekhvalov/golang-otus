package validators

import (
	"fmt"
	"reflect"
)

func newValidatorIn(args string, value interface{}) (Validator, error) {
	builder := validatorBuilder{}
	builder.setFunctionsForKind(
		func(args string) (parsedArgs, error) {
			return parseInts(args)
		},
		func(pArgs parsedArgs, value reflect.Value) error {
			for _, arg := range pArgs.([]int64) {
				if value.Int() == arg {
					return nil
				}
			}
			return fmt.Errorf("unexpected value: %d", value.Int())
		},
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
	)
	builder.setFunctionsForKind(
		func(args string) (parsedArgs, error) {
			return parseUints(args)
		},
		func(pArgs parsedArgs, value reflect.Value) error {
			for _, arg := range pArgs.([]uint64) {
				if value.Uint() == arg {
					return nil
				}
			}
			return fmt.Errorf("unexpected value: %d", value.Uint())
		},
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
	)
	builder.setFunctionsForKind(
		func(args string) (parsedArgs, error) {
			return parseStrings(args)
		},
		func(pArgs parsedArgs, value reflect.Value) error {
			for _, arg := range pArgs.([]string) {
				if value.String() == arg {
					return nil
				}
			}
			return fmt.Errorf("unexpected value: %s", value.String())
		},
		reflect.String,
	)
	return builder.build(args, value)
}
