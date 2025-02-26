package types

import (
	"reflect"
)

// Resources is a map that associates resource names with their corresponding IResource instances.
type Resources map[string]IResource

// Add inserts a new resource into the Resources map. The provided resource must be a non-nil pointer to a struct.
func (r *Resources) Add(res IResource) {
	if reflect.TypeOf(res).Kind() != reflect.Pointer || reflect.TypeOf(res).Elem().Kind() != reflect.Struct {
		panic("resource type must be a non-nil pointer to a struct")
	}
	(*r)[res.Name()] = res
}

// Get retrieves the resource with the specified name from the Resources map. If the resource is not found, it panics.
func (r *Resources) Get(name string) IResource {
	if res, ok := (*r)[name]; ok {
		return res
	}
	panic("resource not found: " + name)
}
