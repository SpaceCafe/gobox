package types

import (
	"errors"
)

const (
	SqliteConstraintNotNull    = 1299
	SqliteConstraintPrimaryKey = 1555
	PostgreSQLNotNullViolation = "23502"
	PostgreSQLUniqueViolation  = "23505"
	AWSEntityTooLarge          = "EntityTooLarge"
)

var (
	ErrNotAuthorized = errors.New("not authorized to perform this operation")
	ErrNotFound      = errors.New("resource not found")
	ErrDuplicatedKey = errors.New("duplicated key")
)
