package validators

import (
	"fmt"
	"reflect"
	"regexp"
)

func newValidatorRegexp(args string, value interface{}) (Validator, error) {
	builder := validatorBuilder{}
	builder.setFunctionsForKind(
		func(args string) (parsedArgs, error) {
			return parseRegexp(args)
		},
		func(pArgs parsedArgs, value reflect.Value) error {
			r := pArgs.(regexp.Regexp)
			if !r.Match([]byte(value.String())) {
				return fmt.Errorf("regexp mismatched")
			}
			return nil
		},
		reflect.String,
	)
	return builder.build(args, value)
}
