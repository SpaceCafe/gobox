package logger

import (
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type message struct {
	Text   string `json:"text"`
	Number int    `json:"number"`
}

func (r message) String() string {
	return r.Text + strconv.Itoa(r.Number)
}

func TestItem_Marshal(t *testing.T) {
	date, _ := time.Parse("2006/01/02 15:04:05", "2024/02/05 09:15:30")

	type fields struct {
		Date    time.Time
		File    string
		Level   string
		Message any
		Line    int
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"full string", fields{
			Date:    date,
			File:    "example.go",
			Level:   "debug",
			Message: "Test message",
			Line:    123,
		}, "{\"date\":\"2024-02-05T09:15:30Z\", \"file\":\"example.go\", \"level\":\"debug\", \"message\":\"Test message\", \"line\":123}"},
		{"minimal string", fields{
			Date: date,
		}, "{\"date\":\"2024-02-05T09:15:30Z\", \"file\":\"\", \"level\":\"\", \"message\":\"\", \"line\":0}"},
		{"minimal json", fields{
			Date: date,
			Message: message{
				Text:   "Test message",
				Number: 123456,
			},
		}, "{\"date\":\"2024-02-05T09:15:30Z\", \"file\":\"\", \"level\":\"\", \"message\":{\"text\":\"Test message\", \"number\":123456}, \"line\":0}"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Item{
				Date:    tt.fields.Date,
				File:    tt.fields.File,
				Level:   tt.fields.Level,
				Message: tt.fields.Message,
				Line:    tt.fields.Line,
			}
			got, err := r.Marshal()
			assert.NoError(t, err)
			assert.JSONEq(t, tt.want, string(got))
		})
	}
}

func TestItem_String(t *testing.T) {
	date, err := time.Parse("2006/01/02 15:04:05", "2024/02/05 09:15:30")
	if err != nil {
		t.Fatalf("failed to parse date: %v", err)
	}

	type fields struct {
		Date    time.Time
		File    string
		Message any
		Line    int
		level   Level
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"full text", fields{
			Date:    date,
			File:    "example.go",
			Message: "Test message",
			Line:    123,
			level:   DebugLevel,
		}, "[DEBUG]   2024/02/05 09:15:30 example.go:123: Test message"},
		{"minimal text", fields{
			Date: date,
		}, "[DEBUG]   2024/02/05 09:15:30 :0: "},
		{"minimal json", fields{
			Date: date,
			Message: message{
				Text:   "Test message",
				Number: 123456,
			},
		}, "[DEBUG]   2024/02/05 09:15:30 :0: Test message123456"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Item{
				Date:    tt.fields.Date,
				File:    tt.fields.File,
				Message: tt.fields.Message,
				Line:    tt.fields.Line,
				level:   tt.fields.level,
			}
			assert.Equal(t, tt.want, r.String())
		})
	}
}

func TestNewItem(t *testing.T) {
	type args struct {
		level   Level
		file    string
		line    int
		message string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"Test1", args{DebugLevel, "example.go", 123, "Test message"}, "debug"},
		{"Test2", args{InfoLevel, "", 0, ""}, "info"},
		{"Test3", args{WarningLevel, "", 0, ""}, "warning"},
		{"Test4", args{ErrorLevel, "", 0, ""}, "error"},
		{"Test5", args{FatalLevel, "", 0, ""}, "fatal"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewItem(tt.args.level, tt.args.file, tt.args.line, tt.args.message)
			assert.NotNil(t, got.Date)
			assert.Equal(t, tt.args.file, got.File)
			assert.Equal(t, tt.want, got.Level)
			assert.Equal(t, tt.args.line, got.Line)
			assert.Equal(t, tt.args.message, got.Message)
			assert.Equal(t, tt.args.level, got.level)
		})
	}
}
