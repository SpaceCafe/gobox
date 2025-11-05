package httpserver

import (
	"errors"
)

var (
	ErrInvalidBasePath = errors.New(
		"http-server base path must be absolute and not end with a slash",
	)
	ErrInvalidPort = errors.New(
		"http-server port must be a number between 1 and 65535",
	)
	ErrInvalidReadHeaderTimeout = errors.New(
		"http-server read header timeout must be greater than 0",
	)
	ErrInvalidReadTimeout = errors.New("http-server read timeout must be greater than 0")
	ErrNoCertFile         = errors.New("http-server key file is set but cert_file is empty")
	ErrNoContext          = errors.New("http-server context can not be empty")
	ErrNoHost             = errors.New("http-server host cannot be empty")
	ErrNoKeyFile          = errors.New("http-server cert file is set but key-file is empty")
)
