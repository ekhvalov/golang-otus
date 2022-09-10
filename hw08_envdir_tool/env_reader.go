package main

import (
	"os"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

const (
	newLine  byte = '\n'
	zeroByte byte = 0x0000
)

var rightTrimBytes = []byte{' ', '\t'}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	env := Environment{}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			return nil, err
		}
		if info.Size() == 0 {
			env[info.Name()] = EnvValue{NeedRemove: true}
			continue
		}
		value, err := getValueFromFile(dir + "/" + info.Name())
		if err != nil {
			return nil, err
		}
		env[entry.Name()] = EnvValue{Value: value}
	}
	return env, nil
}

func getValueFromFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	line := getFirstLine(data)
	return string(sanitize(line)), nil
}

func getFirstLine(data []byte) []byte {
	for i, d := range data {
		if d == newLine {
			return data[:i]
		}
	}
	return data
}

func sanitize(line []byte) []byte {
	for i, b := range line {
		if b == zeroByte {
			line[i] = newLine
		}
	}
	return trimRight(line)
}

func trimRight(line []byte) []byte {
	for len(line) > 0 {
		if !shouldBeTrimmed(line[len(line)-1]) {
			return line
		}
		line = line[:len(line)-1]
	}
	return line
}

func shouldBeTrimmed(b byte) bool {
	for _, trimByte := range rightTrimBytes {
		if b == trimByte {
			return true
		}
	}
	return false
}
