package problems

import (
	"net/http"
	"strings"
	"unicode"

	"github.com/gin-gonic/gin"
)

type IProblem interface {
	Error() (msg string)
	WithError(err error) (newProblem *Problem)
	Abort(ctx *gin.Context)
	appendDetail(err error)
}

// Problem is used for a standardised error handling for the REST API that the IETF
// has worked out in [RFC 7807](https://datatracker.ietf.org/doc/html/rfc7807).
type Problem struct {

	// Type is a URI reference that identifies the problem type.
	Type string `json:"type"`

	// Title is a short, human-readable summary of the problem.
	Title string `json:"title"`

	// Status equals the HTTP status code generated by the origin server for this occurrence of the problem.
	// A list with status codes and their purpose:
	//   - 400 Bad Request: the input parameters are incorrect or missing, or the request itself is incomplete.
	//   - 401 Unauthorized: the request is unauthenticated.
	//   - 403 Forbidden: the client is not authorized to perform this request.
	//   - 404 Not Found: the resource does not exist.
	//   - 405 Method Not Allowed: the HTTP method is not allowed for the requested resource.
	//   - 406 Not Acceptable: the Accept header does not match. Also, can be used to refuse request.
	//   - 409 Conflict: an attempt is made for a duplicate create operation.
	//   - 429 Too Many Requests: a user sends too many requests in a given amount of time
	//   - 500 Internal Server Error: a generic server error
	//   - 502 Bad Gateway: the upstream server or third-party service calls fail
	//   - 503 Service Unavailable: something unexpected happened at the server
	Status int `json:"status"`

	// Detail contains a human-readable explanation specific to this occurrence of the problem.
	Detail string `json:"detail"`

	// Instance is a URI reference that identifies the specific occurrence of the problem.
	// Example: /item/list
	Instance string `json:"instance"`
}

// NewProblem creates a new instance of Problem with default values, if not specified by function parameters.
func NewProblem(problemType string, title string, status int, detail string) *Problem {
	p := &Problem{
		Type:   "/errors/unspecified-error",
		Title:  "Unspecified error",
		Status: http.StatusInternalServerError,
		Detail: "This type of error was not specified",
	}

	if len(problemType) > 0 {
		p.Type = "/errors/" + problemType
	}

	if len(problemType) == 0 && len(title) > 0 {
		// Convert title to kebab case and set it as type.
		words := strings.FieldsFunc(title, func(r rune) bool {
			return !unicode.IsLetter(r) && !unicode.IsNumber(r)
		})
		p.Type = "/errors/" + strings.ToLower(strings.Join(words, "-"))
	}

	if len(title) > 0 {
		p.Title = title
	}

	if status >= 400 && status < 600 {
		p.Status = status
	}

	if len(detail) > 0 {
		p.Detail = detail
	}

	return p
}

// NewProblemWithError creates a new instance of Problem and set its error field using an existing error.
func NewProblemWithError(problemType string, title string, status int, detail string, err error) *Problem {
	p := NewProblem(problemType, title, status, detail)
	p.appendDetail(err.Error())
	return p
}

// Error returns a formatted error message for the Problem struct type.
func (r *Problem) Error() (msg string) {
	return r.Title + ": " + r.Detail
}

// WithError creates a copy of the Problem instance and adds the error message to details.
func (r *Problem) WithError(err error) (newProblem *Problem) {
	p := *r
	p.appendDetail(err.Error())
	return &p
}

// WithDetail creates a copy of the Problem instance and adds text to details.
func (r *Problem) WithDetail(text string) (newProblem *Problem) {
	p := *r
	p.appendDetail(text)
	return &p
}

// Abort aborts the current HTTP request and attaches the Problem instance to the context's error.
func (r *Problem) Abort(ctx *gin.Context) {
	_ = ctx.Error(r)
	ctx.Abort()
}

// appendDetail adds a text or the human-readable error message to the detail field of the problem.
func (r *Problem) appendDetail(text string) {
	if text != "" {
		r.Detail = r.Detail + "\n\nReason: " + text
	}
}
