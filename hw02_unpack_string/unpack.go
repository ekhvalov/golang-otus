package hw02unpackstring

import (
	"errors"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

const (
	asciiDecimalZero = 48
)

func Unpack(str string) (string, error) {
	result := strings.Builder{}
	var buffer rune
	for _, symbol := range str {
		if buffer == 0 {
			if isDigit(symbol) {
				return "", ErrInvalidString
			}
			buffer = symbol
			continue
		}
		if isDigit(symbol) {
			n := toDigit(symbol)
			for j := 0; j < n; j++ {
				result.WriteRune(buffer)
			}
			buffer = 0
		} else {
			result.WriteRune(buffer)
			buffer = symbol
		}
	}
	if buffer == 0 {
		return result.String(), nil
	}
	if isDigit(buffer) {
		return "", ErrInvalidString
	}
	result.WriteRune(buffer)
	return result.String(), nil
}

func isDigit(symbol rune) bool {
	return unicode.IsDigit(symbol)
}

func toDigit(digit rune) int {
	return int(byte(digit) - asciiDecimalZero)
}
