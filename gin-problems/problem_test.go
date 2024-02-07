package problems

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewProblemWithError(t *testing.T) {
	type args struct {
		problemType string
		title       string
		status      int
		detail      string
		err         error
	}
	tests := []struct {
		name string
		args args
		want *Problem
	}{
		{
			name: "Test1",
			args: args{problemType: "test-error", title: "Test Error", status: 400, detail: "This is a test error", err: errors.New("it is just a test")},
			want: &Problem{Type: "/errors/test-error", Title: "Test Error", Status: 400, Detail: "This is a test error\n\nReason: it is just a test"},
		},
		{
			name: "Test2",
			args: args{title: "Invalid Request", status: 401},
			want: &Problem{Type: "/errors/invalid-request", Title: "Invalid Request", Status: 401, Detail: "This type of error was not specified"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.EqualValues(t, tt.want, NewProblemWithError(tt.args.problemType, tt.args.title, tt.args.status, tt.args.detail, tt.args.err))
		})
	}
}

func TestProblem_Error(t *testing.T) {
	type fields struct {
		Type     string
		Title    string
		Status   int
		Detail   string
		Instance string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Test1",
			fields: fields{
				Type:     "Type",
				Title:    "Test Title",
				Status:   400,
				Detail:   "This is a test error",
				Instance: "",
			},
			want: "Test Title: This is a test error",
		},
		{
			name:   "Test2",
			fields: fields{},
			want:   ": ",
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
			assert.Equal(t, tt.want, r.Error())
		})
	}
}

func TestProblem_WithError(t *testing.T) {
	type args struct {
		problem *Problem
		err     error
	}
	tests := []struct {
		name string
		args args
		want *Problem
	}{
		{
			name: "Test1",
			args: args{
				problem: &Problem{
					Detail: "This is a test error",
				},
				err: errors.New("it is just a test"),
			},
			want: &Problem{
				Detail: "This is a test error\n\nReason: it is just a test",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.args.problem.WithError(tt.args.err)
			assert.EqualValues(t, tt.want, got)
			assert.NotEqual(t, tt.args.problem, got)
		})
	}
}
