package validators

import (
	"fmt"
	"reflect"
	"sync"
)

type (
	parsedArgs    interface{}
	parseArgsFunc func(args string) (parsedArgs, error)
	validateFunc  func(pArgs parsedArgs, value reflect.Value) error
)

type validatorBuilder struct {
	parseArgsFunctions map[reflect.Kind]parseArgsFunc
	validateFunctions  map[reflect.Kind]validateFunc
	once               sync.Once
}

func (b *validatorBuilder) init() {
	b.once.Do(func() {
		b.parseArgsFunctions = make(map[reflect.Kind]parseArgsFunc)
		b.validateFunctions = make(map[reflect.Kind]validateFunc)
	})
}

func (b *validatorBuilder) setFunctionsForKind(
	parseArgsF parseArgsFunc,
	validateF validateFunc,
	kind reflect.Kind,
	kinds ...reflect.Kind,
) {
	b.init()
	b.parseArgsFunctions[kind] = parseArgsF
	b.validateFunctions[kind] = validateF
	for _, k := range kinds {
		b.parseArgsFunctions[k] = parseArgsF
		b.validateFunctions[k] = validateF
	}
}

func (b *validatorBuilder) build(args string, value interface{}) (Validator, error) {
	v := reflect.ValueOf(value)
	isSlice := false
	kind := v.Kind()
	if kind == reflect.Slice {
		if v.IsNil() {
			return nil, fmt.Errorf("can not create validator for <nil> value")
		}
		isSlice = true
		kind = v.Type().Elem().Kind()
	}
	parseArgs, ok := b.parseArgsFunctions[kind]
	if !ok {
		return nil, fmt.Errorf("unsupported type: %T", value)
	}
	pArgs, err := parseArgs(args)
	if err != nil {
		return nil, fmt.Errorf("argument parse error: %w", err)
	}
	return validator{
		pArgs:    pArgs,
		value:    v,
		isSlice:  isSlice,
		validate: b.validateFunctions[kind],
	}, nil
}
