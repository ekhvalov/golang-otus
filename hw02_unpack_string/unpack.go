package hw02unpackstring

import (
	"errors"
	"strings"
)

var ErrInvalidString = errors.New("invalid string")

const (
	asciiDecimalZero = 48
	asciiDecimalNine = 57
)

func Unpack(str string) (string, error) {
	result := strings.Builder{}
	var buffer byte
	for i := 0; i < len(str); i++ {
		symbol := str[i]
		if buffer == 0 {
			if isDigit(symbol) {
				return "", ErrInvalidString
			}
			buffer = symbol
			continue
		}
		if isDigit(symbol) {
			n := toDigit(symbol)
			for j := 0; j < n; j++ { // TODO: Optimize?
				result.WriteByte(buffer)
			}
			buffer = 0
		} else {
			result.WriteByte(buffer)
			buffer = symbol
		}
	}
	if buffer != 0 {
		if isDigit(buffer) {
			return "", ErrInvalidString
		}
		result.WriteByte(buffer)
	}
	return result.String(), nil
}

func isDigit(b byte) bool {
	return b >= asciiDecimalZero && b <= asciiDecimalNine
}

func toDigit(b byte) int {
	return int(b - asciiDecimalZero)
}
