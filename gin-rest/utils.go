package rest

import (
	"strings"
	"unicode"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// CamelToSnake converts a camel case string to a snake case string.
func CamelToSnake(s string) string {
	var result strings.Builder
	for i, r := range s {
		if unicode.IsUpper(r) && i != 0 {
			result.WriteRune('_')
		}
		result.WriteRune(unicode.ToLower(r))
	}
	return result.String()
}

// IsLetter checks if the given string consists entirely of letter characters.
func IsLetter(s string) bool {
	return !strings.ContainsFunc(s, func(r rune) bool {
		return !unicode.IsLetter(r)
	})
}

// AddClauses adds multiple clauses to a GORM database statement.
func AddClauses(stmt *gorm.Statement, clauses []clause.Interface) {
	for i := range clauses {
		stmt.AddClause(clauses[i])
	}
}
