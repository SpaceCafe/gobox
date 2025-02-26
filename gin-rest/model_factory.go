package rest

import (
	"reflect"
)

// ModelFactory creates a new instance of the type represented by the provided data.
// It handles pointers, slices, and structs to return an appropriate new instance.
// This is useful for creating instances of models where you need to generate
// a new model object without knowing its specific type at compile time.
func ModelFactory(data any) any {
	rt := reflect.TypeOf(data)

	// Check if it is a pointer and dereference it if necessary.
	if rt.Kind() == reflect.Pointer {
		rt = rt.Elem()
	}

	// Check if it is a slice and dereference it if necessary.
	if rt.Kind() == reflect.Slice {
		rt = rt.Elem()
	}

	// Check if it is a pointer and dereference it if necessary.
	if rt.Kind() == reflect.Pointer {
		rt = rt.Elem()
	}

	// Check if it is a struct and return a new instance of the struct.
	if rt.Kind() == reflect.Struct {
		return reflect.New(rt).Interface()
	}

	return &struct{}{}
}
