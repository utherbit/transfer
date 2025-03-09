package internal

import (
	"errors"
	"go/token"
	"path/filepath"
	"strconv"
	"strings"
)

// Регулярное выражение для проверки формата "file.go:line:column" или "file.go:line"
func isValidGoFilePosition(s string) bool {
	pos, err := ParsePosition(s)
	if err != nil {
		return false
	}

	if filepath.Ext(pos.Filename) != ".go" {
		return false
	}

	return pos.IsValid()
}

func ParsePosition(ref string) (pos token.Position, err error) {
	items := strings.Split(ref, ":")

	pos.Filename = items[0]

	if len(items) > 1 {
		pos.Line, err = strconv.Atoi(items[1])
		if err != nil {
			err = errors.New("ParsePosition: 'format - dir:line', line argument must be number")
			return
		}
	} else {
		return pos, errors.New("ParsePosition: ref must be contains line number")
	}

	if len(items) > 2 {
		pos.Column, err = strconv.Atoi(items[2])
		if err != nil {
			err = errors.New("ParsePosition: 'format - dir:line:column', column argument must be number")
			return
		}
	}

	return pos, nil
}
