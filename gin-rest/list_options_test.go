package rest

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/spacecafe/gobox/gin-rest/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type MockModel1 struct{}

func (m *MockModel1) Sortable(_ *gin.Context) map[string]struct{} {
	return map[string]struct{}{"name": {}, "created_at": {}, "description": {}}
}

func (m *MockModel1) Filterable(_ *gin.Context) map[string]struct{} {
	return map[string]struct{}{"name": {}, "created_at": {}, "description": {}}
}

type MockModel2 struct{}

func TestListOptions(t *testing.T) {
	err := InitializeValidators()
	require.NoError(t, err)

	type args struct {
		request *http.Request
		entity  any
	}
	type wants struct {
		err     error
		options *ListOptions
	}
	tests := []struct {
		name  string
		args  args
		wants wants
	}{
		{
			name: "default values",
			args: args{
				request: &http.Request{URL: &url.URL{}},
				entity:  &MockModel1{},
			},
			wants: wants{
				options: &ListOptions{
					Page:     1,
					PageSize: 10,
					Sorts:    &[]types.SortOption{},
					Filters:  &[]types.FilterOption{},
				},
			},
		},
		{
			name: "custom values",
			args: args{
				request: &http.Request{
					URL: &url.URL{
						RawQuery: "page[number]=2&page[size]=5&sort=%2Bname,-created_at,DescriPtion,%2Binvalid&filter[name]=test&filter[DescriPtion][ne]=lorem&filter[inval%00id][ne]=7&filter[unknown]=3",
					},
				},
				entity: &MockModel1{},
			},
			wants: wants{
				options: &ListOptions{
					Page:     2,
					PageSize: 5,
					Sorts: &[]types.SortOption{
						{Field: "name", Descending: false},
						{Field: "created_at", Descending: true},
						{Field: "description", Descending: false},
					},
					Filters: &[]types.FilterOption{
						{Field: "name", Value: "test", Operator: types.Equals},
						{Field: "description", Value: "lorem", Operator: types.NotEquals},
					},
				},
			},
		},
		{
			name: "invalid page",
			args: args{
				request: &http.Request{
					URL: &url.URL{
						RawQuery: "page[number]=-1",
					},
				},
				entity: &MockModel1{},
			},
			wants: wants{
				err: validator.ValidationErrors{},
				options: &ListOptions{
					Page:     1,
					PageSize: 10,
					Sorts:    &[]types.SortOption{},
					Filters:  &[]types.FilterOption{},
				},
			},
		},
		{
			name: "simple model",
			args: args{
				request: &http.Request{URL: &url.URL{}},
				entity:  &MockModel2{},
			},
			wants: wants{
				options: &ListOptions{
					Page:     1,
					PageSize: 10,
					Sorts:    &[]types.SortOption{},
					Filters:  &[]types.FilterOption{},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
			ctx.Request = tt.args.request

			err = ParseListOptions(ctx, tt.args.entity)
			if tt.wants.err != nil {
				assert.ErrorAs(t, err, &tt.wants.err)
			} else {
				assert.NoError(t, err)
				listOptions := GetListOptions(ctx)
				assert.Equal(t, tt.wants.options.Page, listOptions.Page)
				assert.Equal(t, tt.wants.options.PageSize, listOptions.PageSize)
				assert.Equal(t, *tt.wants.options.Sorts, *listOptions.Sorts)
				assert.Len(t, *listOptions.Filters, len(*tt.wants.options.Filters))
				for _, wantFilter := range *tt.wants.options.Filters {
					found := false
					for _, actualFilter := range *listOptions.Filters {
						if wantFilter.Field == actualFilter.Field &&
							wantFilter.Value == actualFilter.Value &&
							wantFilter.Operator == actualFilter.Operator {
							found = true
							break
						}
					}
					assert.True(t, found, "filter %v not found in actual filters", wantFilter)
				}
			}
		})
	}
}
