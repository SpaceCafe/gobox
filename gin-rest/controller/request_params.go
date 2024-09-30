package controller

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spacecafe/gobox/gin-rest/types"
)

const (

	// DefaultPage is the default page number for pagination.
	DefaultPage = 1

	// DefaultPageSize is the default number of items per page.
	DefaultPageSize = 10

	ErrOptionsNil = "options must not be nil"
)

// RequestParams holds pagination, sorting, and filtering parameters from the request.
type RequestParams struct {
	Page     int    `form:"page[number]" binding:"omitempty,min=1"`
	PageSize int    `form:"page[size]" binding:"omitempty,min=1,max=100"`
	Sort     string `form:"sort" binding:"omitempty,sort_by_attributes"`
	options  *types.ServiceOptions
}

// NewRequestParams creates a new RequestParams instance and binds query parameters from the context.
// It also completes the service options with pagination, sorting, and filtering parameters.
func NewRequestParams(ctx *gin.Context, options *types.ServiceOptions, hasFieldFn func(string) bool) *RequestParams {
	request := &RequestParams{
		Page:     DefaultPage,
		PageSize: DefaultPageSize,
	}
	_ = ctx.ShouldBindQuery(request)

	request.options = options
	request.options.Page = request.Page
	request.options.PageSize = request.PageSize
	request.options.Filters = parseFilter(ctx, hasFieldFn)
	request.options.Sorts = request.parseSort(hasFieldFn)

	return request
}

// Options returns the service options, panicking if they are nil.
func (r *RequestParams) Options() *types.ServiceOptions {
	if r.options == nil {
		panic(ErrOptionsNil)
	}
	return r.options
}

// parseSort parses the sort parameters from the request and returns a slice of SortOption.
// It validates the fields using the provided hasFieldFn function.
func (r *RequestParams) parseSort(hasFieldFn func(string) bool) *[]types.SortOption {
	options := make([]types.SortOption, 0)
	if r.Sort == "" {
		return &options
	}

	params := strings.Split(r.Sort, ",")
	for _, param := range params {
		option := types.SortOption{Field: param}
		if strings.HasPrefix(param, "-") {
			option.Field = param[1:]
			option.Descending = true
		} else if strings.HasPrefix(param, "+") {
			option.Field = param[1:]
		}

		// Check if field is readable.
		if hasFieldFn(option.Field) {
			options = append(options, option)
		}
	}
	return &options
}

// parseFilter parses the filter parameters from the request and returns a slice of FilterOption.
// It validates the fields and operators using the provided hasFieldFn function.
func parseFilter(ctx *gin.Context, hasFieldFn func(string) bool) *[]types.FilterOption {
	options := make([]types.FilterOption, 0)
	if ctx.Request != nil && ctx.Request.URL != nil {
		for key, value := range ctx.Request.URL.Query() {

			// Search for filter parameter.
			if strings.HasPrefix(key, "filter[") {
				option := createFilterOption(key, value[0], hasFieldFn)
				if option != nil {
					options = append(options, *option)
				}
			}
		}
	}
	return &options
}

// createFilterOption creates a FilterOption based on the provided key and value.
// Returns a pointer to a FilterOption if the key is valid, otherwise returns nil.
func createFilterOption(key, value string, hasFieldFn func(string) bool) *types.FilterOption {
	trimmedKey := key[7:]
	endIndex := strings.IndexByte(trimmedKey, ']')
	if endIndex <= 0 || !hasFieldFn(trimmedKey[:endIndex]) {
		return nil
	}

	option := &types.FilterOption{
		Field:    trimmedKey[:endIndex],
		Value:    value,
		Operator: types.Equals,
	}

	if len(trimmedKey) > endIndex+4 && trimmedKey[endIndex+1] == '[' && trimmedKey[len(trimmedKey)-1] == ']' {
		operator := types.FilterOperator(trimmedKey[endIndex+2 : len(trimmedKey)-1])
		if _, ok := types.FilterOperators[operator]; ok {
			option.Operator = operator
		}
	}
	return option
}
