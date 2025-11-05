package terminator

import (
	"errors"
)

var ErrInvalidTimeout = errors.New("terminator timeout must be greater than 0")
