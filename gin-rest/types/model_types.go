package types

import (
	"github.com/guregu/null/v6"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

// DeletedAt represents a nullable time type that tracks soft deletion timestamps.
type DeletedAt null.Time

// QueryClauses returns the query clauses for handling soft delete operations in queries
func (DeletedAt) QueryClauses(f *schema.Field) []clause.Interface {
	return []clause.Interface{gorm.SoftDeleteQueryClause{Field: f}}
}

// DeleteClauses returns the delete clauses for handling soft delete operations in delete statements
func (DeletedAt) DeleteClauses(f *schema.Field) []clause.Interface {
	return []clause.Interface{gorm.SoftDeleteDeleteClause{Field: f}}
}

// UpdateClauses returns the update clauses for handling soft delete operations in update statements
func (DeletedAt) UpdateClauses(f *schema.Field) []clause.Interface {
	return []clause.Interface{gorm.SoftDeleteUpdateClause{Field: f}}
}
