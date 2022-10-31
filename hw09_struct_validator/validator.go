package hw09structvalidator

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/ekhvalov/hw09_struct_validator/validators"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

const (
	tagValidate = "validate"
)

func (v ValidationErrors) Error() string {
	sb := strings.Builder{}
	for _, validationError := range v {
		sb.WriteString(fmt.Sprintf("%s: %v; ", validationError.Field, validationError.Err))
	}
	return strings.TrimRight(sb.String(), "; ")
}

func Validate(v interface{}) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Struct {
		return fmt.Errorf("wrong validated type; expected struct, got: %T", v)
	}
	rt := reflect.TypeOf(v)
	factory := validators.ValidatorFactory{}
	errors := make([]ValidationError, 0)
	for i := 0; i < rt.NumField(); i++ {
		rf := rt.Field(i)
		if tagValue, ok := rf.Tag.Lookup(tagValidate); ok {
			validator, err := factory.CreateValidator(validators.TagArguments(tagValue), rv.Field(i).Interface())
			if err != nil {
				return fmt.Errorf("error while create validator for field '%s': %w", rf.Name, err)
			}
			for _, e := range validator.Validate() {
				errors = append(errors, ValidationError{
					Field: rf.Name,
					Err:   e,
				})
			}
		}
	}
	if len(errors) > 0 {
		return ValidationErrors(errors)
	}
	return nil
}
