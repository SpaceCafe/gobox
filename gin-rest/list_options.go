package rest

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spacecafe/gobox/gin-rest/types"
	"gorm.io/gorm/clause"
)

const (
	DefaultPage     = 1
	DefaultPageSize = 10
)

// ListOptionsRequestParams represents the query parameters for pagination, sorting, and filtering.
type ListOptionsRequestParams struct {
	Page     int    `form:"page[number]" binding:"omitempty,min=1"`
	PageSize int    `form:"page[size]" binding:"omitempty,min=1,max=100"`
	Sort     string `form:"sort" binding:"omitempty,sort_by_attributes"`
}

// SortOptions parses the sort parameters from the request and returns a slice of SortOption.
// It validates the fields using the provided entity's sortable attributes.
func (r *ListOptionsRequestParams) SortOptions(ctx *gin.Context, entity any) *[]types.SortOption {
	sortableEntity, ok := entity.(types.IModelSortable)
	if !ok || r.Sort == "" {
		return &[]types.SortOption{}
	}

	sortable := sortableEntity.Sortable(ctx)
	params := strings.Split(r.Sort, ",")
	options := make([]types.SortOption, 0)
	for _, param := range params {
		option := types.SortOption{Field: strings.ToLower(param), Descending: false}
		if strings.HasPrefix(param, "-") {
			option.Field = option.Field[1:]
			option.Descending = true
		} else if strings.HasPrefix(param, "+") {
			option.Field = option.Field[1:]
		}

		// Check if field is sortable
		if _, ok = sortable[option.Field]; ok {
			options = append(options, option)
		}
	}
	return &options
}

// FilterOptions parses the filter parameters from the request and returns a slice of FilterOption.
// It validates the fields using the provided entity's filterable attributes.
func (*ListOptionsRequestParams) FilterOptions(ctx *gin.Context, entity any) *[]types.FilterOption {
	filterableEntity, ok := entity.(types.IModelFilterable)
	if !ok || ctx.Request == nil || ctx.Request.URL == nil {
		return &[]types.FilterOption{}
	}

	filterable := filterableEntity.Filterable(ctx)
	options := make([]types.FilterOption, 0)

	for key, value := range ctx.Request.URL.Query() {
		if !strings.HasPrefix(key, "filter[") {
			continue
		}

		field, operator, ok := parseFilterKey(key)
		if !ok || len(value) == 0 {
			continue
		}

		if _, ok = filterable[field]; !ok {
			continue
		}

		option := types.FilterOption{
			Field:    field,
			Value:    value[0],
			Operator: operator,
		}
		options = append(options, option)
	}

	return &options
}

// parseFilterKey parses a filter key from the query parameters to extract the field and operator.
// It returns the field name, the corresponding FilterOperator, and a boolean indicating success.
func parseFilterKey(key string) (field string, operator types.FilterOperator, ok bool) {
	endIndex := strings.IndexByte(key, ']')
	if endIndex <= 0 || !IsLetter(key[7:endIndex]) {
		return "", "", false
	}

	field = strings.ToLower(key[7:endIndex])
	operator = types.Equals

	if len(key) > endIndex+4 && key[endIndex+1] == '[' && key[len(key)-1] == ']' {
		operator = types.FilterOperator(key[endIndex+2 : len(key)-1])
		if _, ok = types.FilterOperators[operator]; !ok {
			return "", "", false
		}
	}

	return field, operator, true
}

// ListOptions encapsulates the pagination, sorting, and filtering options for database queries.
type ListOptions struct {
	Page     int
	PageSize int
	Sorts    *[]types.SortOption
	Filters  *[]types.FilterOption
}

// ParseListOptions parses the list options from the query parameters and sets them in the context.
func ParseListOptions(ctx *gin.Context, entity any) (err error) {
	params := &ListOptionsRequestParams{
		Page:     DefaultPage,
		PageSize: DefaultPageSize,
	}
	if err = ctx.ShouldBindQuery(params); err != nil {
		return
	}
	ctx.Set(types.ContextDataListOptions, &ListOptions{
		Page:     params.Page,
		PageSize: params.PageSize,
		Sorts:    params.SortOptions(ctx, entity),
		Filters:  params.FilterOptions(ctx, entity),
	})
	return nil
}

// GetListOptions retrieves pagination, sorting, and filtering options from the context.
func GetListOptions(ctx *gin.Context) *ListOptions {
	if val, ok := ctx.Get(types.ContextDataListOptions); ok {
		return val.(*ListOptions)
	}
	return &ListOptions{
		Page:     DefaultPage,
		PageSize: DefaultPageSize,
		Sorts:    &[]types.SortOption{},
		Filters:  &[]types.FilterOption{},
	}
}

// Paginate returns a clause to paginate the database result.
func (r *ListOptions) Paginate() clause.Interface {
	return clause.Limit{Limit: &r.PageSize, Offset: (r.Page - 1) * r.PageSize}
}

// Sort returns a clause to sort the database result based on the provided sort options.
func (r *ListOptions) Sort() clause.Interface {
	columns := make([]clause.OrderByColumn, 0, len(*r.Sorts))
	for i := range *r.Sorts {
		columns = append(columns, clause.OrderByColumn{
			Column: clause.Column{Name: (*r.Sorts)[i].Field},
			Desc:   (*r.Sorts)[i].Descending,
		})
	}
	return clause.OrderBy{Columns: columns}
}

// Filter returns a clause to filter the database result based on the provided filter options.
func (r *ListOptions) Filter() clause.Interface {
	expressions := make([]clause.Expression, 0, len(*r.Filters))
	for i := range *r.Filters {
		var expression clause.Expression

		switch (*r.Filters)[i].Operator {
		case types.Equals:
			expression = clause.Eq{Column: (*r.Filters)[i].Field, Value: (*r.Filters)[i].Value}
		case types.NotEquals:
			expression = clause.Neq{Column: (*r.Filters)[i].Field, Value: (*r.Filters)[i].Value}
		case types.GreaterThan:
			expression = clause.Gt{Column: (*r.Filters)[i].Field, Value: (*r.Filters)[i].Value}
		case types.GreaterThanOrEqual:
			expression = clause.Gte{Column: (*r.Filters)[i].Field, Value: (*r.Filters)[i].Value}
		case types.LessThan:
			expression = clause.Lt{Column: (*r.Filters)[i].Field, Value: (*r.Filters)[i].Value}
		case types.LessThanOrEqual:
			expression = clause.Lte{Column: (*r.Filters)[i].Field, Value: (*r.Filters)[i].Value}
		case types.Like:
			expression = clause.Like{Column: (*r.Filters)[i].Field, Value: (*r.Filters)[i].Value}
		default:
		}
		expressions = append(expressions, expression)
	}
	return clause.Where{Exprs: expressions}
}
