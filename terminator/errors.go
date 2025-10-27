package terminator

import (
	"errors"
)

var (
	ErrInvalidTimeout = errors.New("timeout must be greater than 0")
)
