package controller

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/spacecafe/gobox/gin-rest/types"
	"github.com/stretchr/testify/require"
)

func TestNewRequestParams(t *testing.T) {
	err := InitializeValidators()
	require.NoError(t, err)

	type args struct {
		ctx        *gin.Context
		hasFieldFn func(string) bool
	}
	tests := []struct {
		name string
		args args
		want *RequestParams
	}{
		{
			name: "Default values",
			args: args{
				ctx: func() *gin.Context {
					c, _ := gin.CreateTestContext(httptest.NewRecorder())
					c.Request = &http.Request{URL: &url.URL{}}
					return c
				}(),
				hasFieldFn: func(string) bool { return true },
			},
			want: &RequestParams{
				Page:     DefaultPage,
				PageSize: DefaultPageSize,
				options: &types.ServiceOptions{
					Page:     DefaultPage,
					PageSize: DefaultPageSize,
					Filters:  &[]types.FilterOption{},
					Sorts:    &[]types.SortOption{},
				},
			},
		},
		{
			name: "Custom values",
			args: args{
				ctx: func() *gin.Context {
					c, _ := gin.CreateTestContext(httptest.NewRecorder())
					c.Request = &http.Request{
						URL: &url.URL{
							RawQuery: "page[number]=2&page[size]=5&sort=%2Bname,-date,description&filter[name]=test&filter[description][ne]=lorem",
						},
					}
					return c
				}(),
				hasFieldFn: func(field string) bool {
					return field == "name" || field == "date" || field == "description"
				},
			},
			want: &RequestParams{
				Page:     2,
				PageSize: 5,
				Sort:     "+name,-date,description",
				options: &types.ServiceOptions{
					Page:     2,
					PageSize: 5,
					Filters: &[]types.FilterOption{
						{Field: "name", Value: "test", Operator: types.Equals},
						{Field: "description", Value: "lorem", Operator: types.NotEquals},
					},
					Sorts: &[]types.SortOption{
						{Field: "name", Descending: false},
						{Field: "date", Descending: true},
						{Field: "description", Descending: false},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewRequestParams(tt.args.ctx, tt.args.hasFieldFn); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRequestParams() = %+v, want %+v", got, tt.want)
			}
		})
	}
}
