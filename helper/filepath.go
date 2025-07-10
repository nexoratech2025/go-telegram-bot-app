package helper

import (
	"errors"
	"os"
)

func IsPathExists(pathname string) bool {
	_, err := os.Stat(pathname)

	return !errors.Is(err, os.ErrNotExist)
}
