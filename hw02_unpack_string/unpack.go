package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(s string) (string, error) {
	waitRepeatNumber, activeSymbol := false, ""
	var str strings.Builder

	for _, v := range s {
		if repeater, err := strconv.Atoi(string(v)); err == nil {
			if !waitRepeatNumber {
				return "", ErrInvalidString
			}

			str.WriteString(strings.Repeat(activeSymbol, repeater))
			waitRepeatNumber, activeSymbol = false, ""
		} else {
			if activeSymbol != "" {
				str.WriteString(activeSymbol)
			}
			waitRepeatNumber, activeSymbol = true, string(v)
		}
	}

	if activeSymbol != "" {
		str.WriteString(activeSymbol)
	}

	return str.String(), nil
}
