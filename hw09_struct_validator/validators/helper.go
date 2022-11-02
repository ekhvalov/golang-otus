package validators

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const argsSeparator = ","

func parseInts(args string) ([]int64, error) {
	ints := make([]int64, 0)
	for _, s := range strings.Split(args, argsSeparator) {
		i, err := strconv.ParseInt(s, 10, 0)
		if err != nil {
			return nil, err
		}
		ints = append(ints, i)
	}
	return ints, nil
}

func parseUints(args string) ([]uint64, error) {
	uints := make([]uint64, 0)
	for _, s := range strings.Split(args, argsSeparator) {
		i, err := strconv.ParseUint(s, 10, 0)
		if err != nil {
			return nil, err
		}
		uints = append(uints, i)
	}
	return uints, nil
}

func parseStrings(args string) ([]string, error) {
	if args == "" {
		return nil, fmt.Errorf("empty argument")
	}
	return strings.Split(args, argsSeparator), nil
}

func parseRegexp(args string) (regexp.Regexp, error) {
	r, err := regexp.Compile(args)
	return *r, err
}
