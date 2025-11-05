package httpserver

import (
	"net/http"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
	"github.com/spacecafe/gobox/logger"
)

var (
	// ColoredStatusCodes is used to set the color and format of HTTP status codes.
	//nolint:gochecknoglobals // This is a lookup map that needs to be globally accessible.
	ColoredStatusCodes = map[int]string{
		http.StatusOK:       color.GreenString("%d", http.StatusOK),
		http.StatusCreated:  color.GreenString("%d", http.StatusCreated),
		http.StatusAccepted: color.GreenString("%d", http.StatusAccepted),
		http.StatusNonAuthoritativeInfo: color.GreenString(
			"%d",
			http.StatusNonAuthoritativeInfo,
		),
		http.StatusNoContent:    color.GreenString("%d", http.StatusNoContent),
		http.StatusResetContent: color.GreenString("%d", http.StatusResetContent),
		http.StatusPartialContent: color.GreenString(
			"%d",
			http.StatusPartialContent,
		),
		http.StatusMultiStatus: color.GreenString("%d", http.StatusMultiStatus),
		http.StatusAlreadyReported: color.GreenString(
			"%d",
			http.StatusAlreadyReported,
		),
		http.StatusIMUsed: color.GreenString("%d", http.StatusIMUsed),
		http.StatusMultipleChoices: color.WhiteString(
			"%d",
			http.StatusMultipleChoices,
		),
		http.StatusMovedPermanently: color.WhiteString(
			"%d",
			http.StatusMovedPermanently,
		),
		http.StatusFound:       color.WhiteString("%d", http.StatusFound),
		http.StatusSeeOther:    color.WhiteString("%d", http.StatusSeeOther),
		http.StatusNotModified: color.WhiteString("%d", http.StatusNotModified),
		http.StatusUseProxy:    color.WhiteString("%d", http.StatusUseProxy),
		http.StatusTemporaryRedirect: color.WhiteString(
			"%d",
			http.StatusTemporaryRedirect,
		),
		http.StatusPermanentRedirect: color.WhiteString(
			"%d",
			http.StatusPermanentRedirect,
		),
		http.StatusBadRequest:   color.YellowString("%d", http.StatusBadRequest),
		http.StatusUnauthorized: color.YellowString("%d", http.StatusUnauthorized),
		http.StatusPaymentRequired: color.YellowString(
			"%d",
			http.StatusPaymentRequired,
		),
		http.StatusForbidden: color.YellowString("%d", http.StatusForbidden),
		http.StatusNotFound:  color.YellowString("%d", http.StatusNotFound),
		http.StatusMethodNotAllowed: color.YellowString(
			"%d",
			http.StatusMethodNotAllowed,
		),
		http.StatusNotAcceptable: color.YellowString(
			"%d",
			http.StatusNotAcceptable,
		),
		http.StatusProxyAuthRequired: color.YellowString(
			"%d",
			http.StatusProxyAuthRequired,
		),
		http.StatusRequestTimeout: color.YellowString(
			"%d",
			http.StatusRequestTimeout,
		),
		http.StatusConflict: color.YellowString("%d", http.StatusConflict),
		http.StatusGone:     color.YellowString("%d", http.StatusGone),
		http.StatusLengthRequired: color.YellowString(
			"%d",
			http.StatusLengthRequired,
		),
		http.StatusPreconditionFailed: color.YellowString(
			"%d",
			http.StatusPreconditionFailed,
		),
		http.StatusRequestEntityTooLarge: color.YellowString(
			"%d",
			http.StatusRequestEntityTooLarge,
		),
		http.StatusRequestURITooLong: color.YellowString(
			"%d",
			http.StatusRequestURITooLong,
		),
		http.StatusUnsupportedMediaType: color.YellowString(
			"%d",
			http.StatusUnsupportedMediaType,
		),
		http.StatusRequestedRangeNotSatisfiable: color.YellowString(
			"%d",
			http.StatusRequestedRangeNotSatisfiable,
		),
		http.StatusExpectationFailed: color.YellowString(
			"%d",
			http.StatusExpectationFailed,
		),
		http.StatusTeapot: color.YellowString("%d", http.StatusTeapot),
		http.StatusMisdirectedRequest: color.YellowString(
			"%d",
			http.StatusMisdirectedRequest,
		),
		http.StatusUnprocessableEntity: color.YellowString(
			"%d",
			http.StatusUnprocessableEntity,
		),
		http.StatusLocked: color.YellowString("%d", http.StatusLocked),
		http.StatusFailedDependency: color.YellowString(
			"%d",
			http.StatusFailedDependency,
		),
		http.StatusTooEarly: color.YellowString("%d", http.StatusTooEarly),
		http.StatusUpgradeRequired: color.YellowString(
			"%d",
			http.StatusUpgradeRequired,
		),
		http.StatusPreconditionRequired: color.YellowString(
			"%d",
			http.StatusPreconditionRequired,
		),
		http.StatusTooManyRequests: color.YellowString(
			"%d",
			http.StatusTooManyRequests,
		),
		http.StatusRequestHeaderFieldsTooLarge: color.YellowString(
			"%d",
			http.StatusRequestHeaderFieldsTooLarge,
		),
		http.StatusUnavailableForLegalReasons: color.YellowString(
			"%d",
			http.StatusUnavailableForLegalReasons,
		),
		http.StatusInternalServerError: color.RedString(
			"%d",
			http.StatusInternalServerError,
		),
		http.StatusNotImplemented: color.RedString("%d", http.StatusNotImplemented),
		http.StatusBadGateway:     color.RedString("%d", http.StatusBadGateway),
		http.StatusServiceUnavailable: color.RedString(
			"%d",
			http.StatusServiceUnavailable,
		),
		http.StatusGatewayTimeout: color.RedString("%d", http.StatusGatewayTimeout),
		http.StatusHTTPVersionNotSupported: color.RedString(
			"%d",
			http.StatusHTTPVersionNotSupported,
		),
		http.StatusVariantAlsoNegotiates: color.RedString(
			"%d",
			http.StatusVariantAlsoNegotiates,
		),
		http.StatusInsufficientStorage: color.RedString(
			"%d",
			http.StatusInsufficientStorage,
		),
		http.StatusLoopDetected: color.RedString("%d", http.StatusLoopDetected),
		http.StatusNotExtended:  color.RedString("%d", http.StatusNotExtended),
		http.StatusNetworkAuthenticationRequired: color.RedString(
			"%d",
			http.StatusNetworkAuthenticationRequired,
		),
	}

	// ColoredMethods is used to set the color and format of HTTP methods.
	//nolint:gochecknoglobals // This is a lookup map that needs to be globally accessible.
	ColoredMethods = map[string]string{
		http.MethodGet:     color.BlueString(http.MethodGet),
		http.MethodHead:    color.MagentaString(http.MethodHead),
		http.MethodPost:    color.CyanString(http.MethodPost),
		http.MethodPut:     color.YellowString(http.MethodPut),
		http.MethodPatch:   color.GreenString(http.MethodPatch),
		http.MethodDelete:  color.RedString(http.MethodDelete),
		http.MethodConnect: color.WhiteString(http.MethodConnect),
		http.MethodOptions: color.WhiteString(http.MethodOptions),
		http.MethodTrace:   color.WhiteString(http.MethodTrace),
	}
)

// LogEntry represents a single log entry with details about an HTTP request.
type LogEntry struct {
	ClientIP   string        `json:"clientIP"`
	Errors     string        `json:"errors"`
	Method     string        `json:"method"`
	Path       string        `json:"path"`
	Latency    time.Duration `json:"latency,omitempty"`
	Size       int           `json:"size,omitempty"`
	StatusCode int           `json:"statusCode"`
}

// MethodColor returns the HTTP method of the log entry as a colored string.
func (r *LogEntry) MethodColor() string {
	if s, ok := ColoredMethods[r.Method]; ok {
		return s
	}

	return color.WhiteString(r.Method)
}

// StatusCodeColor returns the status code of the log entry as a colored string.
func (r *LogEntry) StatusCodeColor() string {
	if s, ok := ColoredStatusCodes[r.StatusCode]; ok {
		return s
	}

	return color.RedString("%d", r.StatusCode)
}

// String formats the LogEntry into a human-readable string with colored status code and method.
func (r *LogEntry) String() string {
	var builder strings.Builder

	builder.WriteString(r.StatusCodeColor())
	builder.WriteString(" | ")
	logger.LeftPadString(&builder, r.Latency.String(), 13) //nolint:mnd // Fixed string size
	builder.WriteString(" | ")
	logger.LeftPadString(&builder, r.ClientIP, 15) //nolint:mnd // Fixed string size
	builder.WriteString(" | ")
	builder.WriteString(r.MethodColor())
	builder.WriteByte(' ')
	builder.WriteString(r.Path)

	if r.Errors != "" {
		builder.WriteByte(' ')
		builder.WriteString(r.Errors)
	}

	return builder.String()
}

// NewGinLogger creates a gin.HandlerFunc that logs HTTP request details using the provided logger.
func NewGinLogger(log logger.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Start timer.
		start := time.Now()
		path := ctx.Request.URL.Path
		raw := ctx.Request.URL.RawQuery

		// Process request.
		ctx.Next()

		// Stop timer and create log item.
		entry := &LogEntry{
			ClientIP:   ctx.ClientIP(),
			Errors:     ctx.Errors.ByType(gin.ErrorTypePrivate).String(),
			Method:     ctx.Request.Method,
			Latency:    time.Since(start),
			Size:       ctx.Writer.Size(),
			StatusCode: ctx.Writer.Status(),
		}

		if raw != "" {
			path = path + "?" + raw
		}

		entry.Path = path
		log.Info(entry)
	}
}
