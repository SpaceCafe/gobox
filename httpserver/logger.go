package httpserver

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
	"github.com/spacecafe/gobox/logger"
)

type LoggerItem struct {
	ClientIP   string        `json:"client-ip"`
	Errors     string        `json:"errors"`
	Method     string        `json:"method"`
	Path       string        `json:"path"`
	Latency    time.Duration `json:"latency"`
	Size       int           `json:"size"`
	StatusCode int           `json:"status-code"`
}

func (r *LoggerItem) String() string {
	return fmt.Sprintf("%s | %13v | %15s | %-7s %#v\n%s",
		r.StatusCodeColor(),
		r.Latency,
		r.ClientIP,
		r.MethodColor(),
		r.Path,
		r.Errors,
	)
}

func (r *LoggerItem) StatusCodeColor() string {
	code := strconv.Itoa(r.StatusCode)
	switch {
	case r.StatusCode >= http.StatusOK && r.StatusCode < http.StatusMultipleChoices:
		return color.GreenString(code)
	case r.StatusCode >= http.StatusMultipleChoices && r.StatusCode < http.StatusBadRequest:
		return color.WhiteString(code)
	case r.StatusCode >= http.StatusBadRequest && r.StatusCode < http.StatusInternalServerError:
		return color.YellowString(code)
	default:
		return color.RedString(code)
	}
}

func (r *LoggerItem) MethodColor() string {
	switch r.Method {
	case http.MethodGet:
		return color.BlueString(r.Method)
	case http.MethodHead:
		return color.MagentaString(r.Method)
	case http.MethodPost:
		return color.CyanString(r.Method)
	case http.MethodPut:
		return color.YellowString(r.Method)
	case http.MethodPatch:
		return color.GreenString(r.Method)
	case http.MethodDelete:
		return color.RedString(r.Method)
	default:
		return color.WhiteString(r.Method)
	}
}

func Logger(log *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer.
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request.
		c.Next()

		// Stop timer and create log item.
		item := &LoggerItem{
			ClientIP:   c.ClientIP(),
			Errors:     c.Errors.ByType(gin.ErrorTypePrivate).String(),
			Method:     c.Request.Method,
			Latency:    time.Since(start),
			Size:       c.Writer.Size(),
			StatusCode: c.Writer.Status(),
		}

		if raw != "" {
			path = path + "?" + raw
		}

		item.Path = path

		log.Info(item)
	}
}
