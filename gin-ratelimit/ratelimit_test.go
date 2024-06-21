package ratelimit

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	problems "github.com/spacecafe/gobox/gin-problems"
	"gotest.tools/assert"
)

func TestNew(t *testing.T) {
	problemQueueFull := problems.ProblemQueueFull
	problemQueueFull.Instance = "/test"
	problemRequestTimeout := problems.ProblemRequestTimeout
	problemRequestTimeout.Instance = "/test"

	type args struct {
		config *Config
	}
	tests := []struct {
		name            string
		args            args
		expectedStatus  int
		expectedProblem *problems.Problem
	}{
		{
			name: "Test successful request",
			args: args{
				config: &Config{
					MaxBurstRequests:      5,
					BurstDuration:         1 * time.Second,
					MaxConcurrentRequests: 10,
					RequestQueueSize:      10,
					RequestTimeout:        1 * time.Second,
				},
			},
			expectedStatus:  http.StatusOK,
			expectedProblem: nil,
		},
		{
			name: "Test queue full",
			args: args{
				config: &Config{
					MaxBurstRequests:      1,
					BurstDuration:         1 * time.Second,
					MaxConcurrentRequests: 1,
					RequestQueueSize:      1,
					RequestTimeout:        100 * time.Millisecond,
				},
			},
			expectedStatus:  http.StatusTooManyRequests,
			expectedProblem: problemQueueFull,
		},
		{
			name: "Test request timeout in queue",
			args: args{
				config: &Config{
					MaxBurstRequests:      1,
					BurstDuration:         1 * time.Second,
					MaxConcurrentRequests: 1,
					RequestQueueSize:      2,
					RequestTimeout:        200 * time.Millisecond,
				},
			},
			expectedStatus:  http.StatusRequestTimeout,
			expectedProblem: problemRequestTimeout,
		},
		{
			name: "Test request timeout in burst",
			args: args{
				config: &Config{
					MaxBurstRequests:      2,
					BurstDuration:         1 * time.Second,
					MaxConcurrentRequests: 1,
					RequestQueueSize:      2,
					RequestTimeout:        200 * time.Millisecond,
				},
			},
			expectedStatus:  http.StatusRequestTimeout,
			expectedProblem: problemRequestTimeout,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			r := gin.Default()
			r.Use(problems.New())
			r.Use(New(tt.args.config))
			r.GET("/test", func(c *gin.Context) {
				time.Sleep(time.Second)
				c.String(http.StatusOK, "OK")
			})

			go r.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "/test?ignore", nil))

			request := httptest.NewRequest("GET", "/test", nil)
			recorder := httptest.NewRecorder()
			time.Sleep(100 * time.Millisecond)
			r.ServeHTTP(recorder, request)

			response := recorder.Result()
			assert.Equal(t, tt.expectedStatus, response.StatusCode)

			body, err := io.ReadAll(response.Body)
			assert.NilError(t, err)
			_ = response.Body.Close()

			if tt.expectedProblem != nil {
				problem := &problems.Problem{}
				assert.NilError(t, err)
				err = json.Unmarshal(body, problem)
				assert.NilError(t, err)
				assert.DeepEqual(t, tt.expectedProblem, problem)
			}
		})
	}
}

func TestRateLimit_drainBurstChannel(t *testing.T) {
	type fields struct {
		config            *Config
		burstChannel      chan struct{}
		concurrentChannel chan struct{}
		queueChannel      chan struct{}
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Test burst channel draining",
			fields: fields{
				config: &Config{
					MaxBurstRequests: 20,
					BurstDuration:    1 * time.Second,
				},
				burstChannel: make(chan struct{}, 20),
			},
			wantErr: false,
		},
		{
			name: "Test burst channel not draining",
			fields: fields{
				config: &Config{
					MaxBurstRequests: 20,
					BurstDuration:    2 * time.Second,
				},
				burstChannel: make(chan struct{}, 20),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rl := &RateLimit{
				config:            tt.fields.config,
				burstChannel:      tt.fields.burstChannel,
				concurrentChannel: tt.fields.concurrentChannel,
				queueChannel:      tt.fields.queueChannel,
			}

			// Fill the burst channel to its capacity
			for i := 0; i < rl.config.MaxBurstRequests; i++ {
				rl.burstChannel <- struct{}{}
			}

			// Run drainBurstChannel in a separate goroutine
			go rl.drainBurstChannel()

			// Wait for a duration longer than the burst duration to ensure the channel is drained
			time.Sleep(time.Second + 100*time.Millisecond)

			if tt.wantErr {
				// Check if the burst channel is not empty
				select {
				case <-rl.burstChannel:
					// burstChannel is not empty as expected
				default:
					t.Errorf("burstChannel should not be empty but it is")
				}
			} else {
				// Check if the burst channel is empty
				select {
				case <-rl.burstChannel:
					t.Errorf("burstChannel should be empty but it is not")
				default:
					// burstChannel is empty as expected
				}
			}
		})
	}
}
