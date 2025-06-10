package main

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	env := make(Environment)

	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		filePath := filepath.Join(dir, file.Name())
		f, err := os.Open(filePath)
		if err != nil {
			return nil, err
		}
		defer f.Close()

		scanner := bufio.NewScanner(f)
		value := ""
		if scanner.Scan() {
			value = scanner.Text()
		}

		if err := scanner.Err(); err != nil {
			return nil, err
		}

		value = strings.TrimRight(value, " \t")
		value = strings.ReplaceAll(value, "\x00", "\n")

		env[file.Name()] = EnvValue{
			Value:      value,
			NeedRemove: value == "",
		}
	}

	return env, nil
}
