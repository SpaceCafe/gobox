package logger

import (
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type message struct {
	Text   string `json:"text"`
	Number int    `json:"number"`
}

func (r message) String() string {
	return r.Text + strconv.Itoa(r.Number)
}

func TestItem_Marshal(t *testing.T) {
	mockDate, err := time.Parse("2006/01/02 15:04:05", "2024/02/05 09:15:30")
	require.NoError(t, err)

	tests := []struct {
		name  string
		entry *Entry
		want  string
	}{
		{
			"full string",
			&Entry{
				Date:    mockDate,
				File:    "example.go",
				Level:   InfoLevel,
				Message: "Test message",
				Line:    123,
			},
			`{"date":"2024-02-05T09:15:30Z", "file":"example.go", "level":"info", "message":"Test message", "line":123}`,
		},
		{
			"minimal string",
			&Entry{
				Date: mockDate,
			},
			`{"date":"2024-02-05T09:15:30Z", "level":"debug", "message":""}`,
		},
		{
			"minimal json",
			&Entry{
				Date: mockDate,
				Message: message{
					Text:   "Test message",
					Number: 123456,
				},
			},
			`{"date":"2024-02-05T09:15:30Z", "level":"debug", "message":{"text":"Test message", "number":123456}}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.entry.Marshal()
			assert.NoError(t, err)
			assert.JSONEq(t, tt.want, string(got))
		})
	}
}

func TestItem_String(t *testing.T) {
	mockDate, err := time.Parse("2006/01/02 15:04:05", "2024/02/05 09:15:30")
	require.NoError(t, err)

	tests := []struct {
		name  string
		entry *Entry
		want  string
	}{
		{
			"full text",
			&Entry{
				Date:    mockDate,
				File:    "example.go",
				Level:   DebugLevel,
				Message: "Test message",
				Line:    123,
			},
			"Test message",
		},
		{
			"minimal text",
			&Entry{
				Date: mockDate,
			},
			"",
		},
		{
			"minimal json",
			&Entry{
				Date: mockDate,
				Message: message{
					Text:   "Test message",
					Number: 123456,
				},
			},
			"Test message123456",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.entry.String())
		})
	}
}
