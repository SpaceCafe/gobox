package utils

import (
	"reflect"
)

// IsStructPointer checks if the provided input is a non-nil pointer to a struct and returns a boolean value.
func IsStructPointer(in any) bool {
	value := reflect.ValueOf(in)
	return in != nil && value.Kind() == reflect.Ptr && value.Elem().Kind() == reflect.Struct
}
