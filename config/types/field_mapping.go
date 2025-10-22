package types

import (
	"reflect"
)

// FieldMapping represents a mapping between YAML path and struct field.
type FieldMapping struct {
	EnvName   string
	EnvAlias  string
	FieldPath []int
	FieldType reflect.Type
}
