package types

import (
	"errors"
)

const (
	SqliteConstraintNotNull    = 1299
	SqliteConstraintPrimaryKey = 1555
)

var (
	ErrNotAuthorized = errors.New("not authorized to perform this operation")
	ErrNotFound      = errors.New("resource not found")
	ErrDuplicatedKey = errors.New("duplicated key")
	ErrNotChanged    = errors.New("not changed")
)
