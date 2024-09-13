package types

import (
	"strings"

	"gorm.io/gorm/schema"
)

// NamingStrategy is a custom naming strategy that embeds gorm's NamingStrategy.
type NamingStrategy struct {
	schema.NamingStrategy
}

// TableName overrides the default TableName method to remove the "Model" suffix from the table name.
func (r NamingStrategy) TableName(name string) string {
	return r.NamingStrategy.TableName(strings.TrimSuffix(name, "Model"))
}
