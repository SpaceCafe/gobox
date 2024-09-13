package rest

import (
	"strings"
)

// joinMapKeys takes a map with string keys and empty struct values, and a separator string.
// It returns a single string with all the keys joined by the separator.
func joinMapKeys(m map[string]any, sep string) string {
	keys := make([]string, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	return strings.Join(keys, sep)
}
