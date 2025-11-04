package authentication

import (
	"strings"
)

// hasCaseInsensitivePrefix checks if a string has a case-insensitive prefix.
func hasCaseInsensitivePrefix(s, prefix string) bool {
	return len(s) >= len(prefix) && strings.EqualFold(s[:len(prefix)], prefix)
}
