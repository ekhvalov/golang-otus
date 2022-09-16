package validators

import (
	"fmt"
	"strings"
	"sync"
)

const (
	argsValidatorSeparator     = "|"
	argsValidatorNameSeparator = ":"
	validatorMinName           = "min"
	validatorMaxName           = "max"
	validatorLenName           = "len"
	validatorInName            = "in"
	validatorRegexpName        = "regexp"
)

type TagArguments string

type nameArgs struct {
	name string
	args string
}

type ValidatorFactory struct {
	builders map[string]func(string, interface{}) (Validator, error)
	once     sync.Once
}

func (f *ValidatorFactory) CreateValidator(args TagArguments, value interface{}) (Validator, error) {
	f.init()
	validators := make([]Validator, 0)
	for _, na := range parseTagArguments(args) {
		builder, ok := f.builders[na.name]
		if !ok {
			return nil, fmt.Errorf("undefined validator: %s", na.name)
		}
		v, err := builder(na.args, value)
		if err != nil {
			return nil, fmt.Errorf("create validator '%s' error: %v", na.name, err)
		}
		validators = append(validators, v)
	}
	return combinedValidator{validators: validators}, nil
}

func (f *ValidatorFactory) init() {
	f.once.Do(func() {
		f.builders = map[string]func(string, interface{}) (Validator, error){
			validatorMinName:    newValidatorMin,
			validatorMaxName:    newValidatorMax,
			validatorLenName:    newValidatorLen,
			validatorInName:     newValidatorIn,
			validatorRegexpName: newValidatorRegexp,
		}
	})
}

func parseTagArguments(args TagArguments) []nameArgs {
	nArgs := make([]nameArgs, 0)
	validators := strings.Split(string(args), argsValidatorSeparator)
	for len(validators) > 0 {
		tagArgs := strings.SplitN(validators[0], argsValidatorNameSeparator, 2)
		if len(tagArgs) > 1 {
			nArgs = append(nArgs, nameArgs{name: tagArgs[0], args: tagArgs[1]})
		} else {
			nArgs = append(nArgs, nameArgs{name: tagArgs[0]})
		}
		validators = validators[1:]
	}
	return nArgs
}
