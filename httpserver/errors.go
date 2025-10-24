package httpserver

import (
	"errors"
)

var (
	ErrNoHost                   = errors.New("host cannot be empty")
	ErrInvalidBasePath          = errors.New("base path must be absolute and not end with a slash")
	ErrNoCertFile               = errors.New("key file is set but cert_file is empty")
	ErrNoKeyFile                = errors.New("cert file is set but key_file is empty")
	ErrInvalidReadTimeout       = errors.New("read timeout must be greater than 0")
	ErrInvalidReadHeaderTimeout = errors.New("read header timeout must be greater than 0")
	ErrInvalidPort              = errors.New("port must be a number between 1 and 65535")
)
