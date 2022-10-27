package validators

import "reflect"

type Validator interface {
	Validate() []error
}

type validator struct {
	pArgs    parsedArgs
	value    reflect.Value
	isSlice  bool
	validate validateFunc
}

func (v validator) Validate() []error {
	if v.isSlice {
		errors := make([]error, 0)
		for i := 0; i < v.value.Len(); i++ {
			err := v.validate(v.pArgs, v.value.Index(i))
			if err != nil {
				errors = append(errors, err)
			}
		}
		if len(errors) > 0 {
			return errors
		}
		return nil
	}
	err := v.validate(v.pArgs, v.value)
	if err != nil {
		return []error{err}
	}
	return nil
}

type combinedValidator struct {
	validators []Validator
}

func (v combinedValidator) Validate() []error {
	if v.validators == nil {
		return nil
	}
	errors := make([]error, 0)
	for _, val := range v.validators {
		err := val.Validate()
		if err != nil {
			errors = append(errors, err...)
		}
	}
	if len(errors) > 0 {
		return errors
	}
	return nil
}
