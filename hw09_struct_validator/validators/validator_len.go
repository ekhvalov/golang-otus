package validators

import (
	"fmt"
	"reflect"
	"strconv"
	"unicode/utf8"
)

func newValidatorLen(args string, value interface{}) (Validator, error) {
	builder := validatorBuilder{}
	builder.setFunctionsForKind(
		func(args string) (parsedArgs, error) {
			return strconv.ParseUint(args, 10, 0)
		},
		func(pArgs parsedArgs, value reflect.Value) error {
			length := pArgs.(uint64)
			l := utf8.RuneCountInString(value.String())
			if length != uint64(l) {
				return fmt.Errorf("length mismatched; expected %d, got %d", length, l)
			}
			return nil
		},
		reflect.String,
	)
	return builder.build(args, value)
}
