package httpserver

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestNewProblem(t *testing.T) {
	type args struct {
		problemType string
		title       string
		status      int
		detail      string
	}
	tests := []struct {
		name string
		args args
		want *Problem
	}{
		{
			name: "Test1",
			args: args{problemType: "test-error", title: "Test Error", status: 400, detail: "This is a test error"},
			want: &Problem{Type: "/errors/test-error", Title: "Test Error", Status: 400, Detail: "This is a test error"},
		},
		{
			name: "Test2",
			args: args{title: "Invalid Request", status: 401},
			want: &Problem{Type: "/errors/invalid-request", Title: "Invalid Request", Status: 401, Detail: "This type of error was not specified"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.EqualValues(t, tt.want, NewProblem(tt.args.problemType, tt.args.title, tt.args.status, tt.args.detail))
		})
	}
}

func TestProblem_AbortWithError(t *testing.T) {
	type fields struct {
		Type     string
		Title    string
		Status   int
		Detail   string
		Instance string
	}
	type args struct {
		err      error
		instance string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name:   "Test1",
			fields: fields{Title: "Test Error", Status: 400, Detail: "This is a test error"},
			args: args{
				err:      errors.New("test error"),
				instance: "/hello/world",
			},
			want: "{\"type\":\"\",\"title\":\"Test Error\",\"status\":400,\"detail\":\"This is a test error\\n\\nReason: test error\",\"instance\":\"/hello/world\"}",
		},
		{
			name:   "Test2",
			fields: fields{Status: 503},
			args: args{
				err:      nil,
				instance: "",
			},
			want: "{\"type\":\"\",\"title\":\"\",\"status\":503,\"detail\":\"\",\"instance\":\"\"}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Problem{
				Type:     tt.fields.Type,
				Title:    tt.fields.Title,
				Status:   tt.fields.Status,
				Detail:   tt.fields.Detail,
				Instance: tt.fields.Instance,
			}
			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			ctx.Request, _ = http.NewRequest(http.MethodGet, tt.args.instance, nil)
			r.AbortWithError(tt.args.err, ctx)

			assert.Equal(t, "application/problem+json", w.Header().Get("Content-Type"))
			assert.Equal(t, tt.fields.Status, w.Code)
			assert.JSONEq(t, tt.want, w.Body.String())
		})
	}
}
