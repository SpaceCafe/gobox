package job_manager

var (
	_ HookContext = (*KVStorage)(nil)
)

// KVStorage is a map type that allows storing and retrieving any values by string keys.
type KVStorage map[string]any

// Get retrieves the value associated with the given key from the storage.
func (r KVStorage) Get(key string) any {
	return r[key]
}

// Set stores the provided value in the storage under the specified key.
// It overwrites any existing value associated with that key.
func (r KVStorage) Set(key string, value any) {
	r[key] = value
}
