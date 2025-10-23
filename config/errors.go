package config

import (
	"errors"
)

var (
	ErrInvalidConfig      = errors.New("config must be a pointer to a struct")
	ErrConfigFileNotFound = errors.New("config file not found in search paths")
	ErrFieldNotSettable   = errors.New("field is not settable")
)
