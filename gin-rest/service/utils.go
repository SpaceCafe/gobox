package service

import (
	"github.com/spacecafe/gobox/gin-rest/types"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// applyOptions applies pagination, sorting, and filtering options to a GORM database query.
// It returns a modified gorm.DB instance with the applied options.
func applyOptions(db *gorm.DB, options *types.ServiceOptions) (tx *gorm.DB) {
	tx = db.Limit(options.PageSize).Offset((options.Page - 1) * options.PageSize)

	for _, option := range *options.Sorts {
		tx = tx.Order(clause.OrderByColumn{Column: clause.Column{Name: option.Field}, Desc: option.Descending})
	}

	clauseExprs := make([]clause.Expression, 0, len(*options.Filters))
	for _, option := range *options.Filters {
		var clauseExpr clause.Expression

		switch option.Operator {
		case types.Equals:
			clauseExpr = clause.Eq{Column: option.Field, Value: option.Value}
		case types.NotEquals:
			clauseExpr = clause.Neq{Column: option.Field, Value: option.Value}
		case types.GreaterThan:
			clauseExpr = clause.Gt{Column: option.Field, Value: option.Value}
		case types.GreaterThanOrEqual:
			clauseExpr = clause.Gte{Column: option.Field, Value: option.Value}
		case types.LessThan:
			clauseExpr = clause.Lt{Column: option.Field, Value: option.Value}
		case types.LessThanOrEqual:
			clauseExpr = clause.Lte{Column: option.Field, Value: option.Value}
		case types.Like:
			clauseExpr = clause.Like{Column: option.Field, Value: option.Value}
		case types.Fuzzy:
		// 	TODO: implement fuzzy search
		default:
		}
		clauseExprs = append(clauseExprs, clauseExpr)
	}
	return tx.Clauses(clauseExprs...)
}
