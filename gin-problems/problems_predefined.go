package problems

import (
	"net/http"
)

var (
	ProblemUnauthorized = NewProblem(
		"",
		http.StatusText(http.StatusUnauthorized),
		http.StatusUnauthorized,
		"The request does not contains any valid authentication data.",
	)
	ProblemInternalError = NewProblem(
		"",
		http.StatusText(http.StatusInternalServerError),
		http.StatusInternalServerError,
		"An internal error occurred. Please, try again later.",
	)
	ProblemMethodNotAllowed = NewProblem(
		"",
		http.StatusText(http.StatusMethodNotAllowed),
		http.StatusMethodNotAllowed,
		"The specified method is not allowed against this resource.",
	)
	ProblemUnsupportedMediaType = NewProblem(
		"",
		http.StatusText(http.StatusUnsupportedMediaType),
		http.StatusUnsupportedMediaType,
		"The specified media type is not allowed against this resource.",
	)
	ProblemBadRequest = NewProblem(
		"",
		http.StatusText(http.StatusBadRequest),
		http.StatusBadRequest,
		"The received content is not valid for this resource.",
	)
	ProblemNoSuchAccessPoint = NewProblem(
		"",
		"No such access point",
		http.StatusNotFound,
		"The specified access point does not exist.",
	)
	ProblemSignatureDoesNotMatch = NewProblem(
		"",
		"Signature does not match",
		http.StatusBadRequest,
		"The request signature that the server calculated does not match the signature that you provided.",
	)
	ProblemNoSuchKey = NewProblem(
		"",
		"No such key",
		http.StatusNotFound,
		"The specified key does not exist.",
	)
	ProblemInvalidFilterKey = NewProblem(
		"",
		"Invalid filter key",
		http.StatusBadRequest,
		"The specified key cannot be used to filter.",
	)
	ProblemCSRFMissing = NewProblem(
		"",
		"CSRF missing",
		http.StatusForbidden,
		"The CSRF token is missing.",
	)
	ProblemCSRFMalfunction = NewProblem(
		"",
		"CSRF malfunction",
		http.StatusInternalServerError,
		"The server could not create a CSRF token.",
	)
	ProblemCSRFInvalid = NewProblem(
		"",
		"CSRF invalid",
		http.StatusForbidden,
		"The provided CSRF token is malformed or otherwise not valid.",
	)
	ProblemJWTMissing = NewProblem(
		"",
		"JWT missing",
		http.StatusUnauthorized,
		"The JSON Web Token (JWT) is missing.",
	)
	ProblemJWTInvalid = NewProblem(
		"",
		"JWT invalid",
		http.StatusForbidden,
		"The provided JSON Web Token (JWT) is malformed or otherwise not valid.",
	)
	ProblemInsufficientPermission = NewProblem(
		"",
		"Insufficient permission",
		http.StatusForbidden,
		"Insufficient permission",
	)
	ProblemRequestTimeout = NewProblem(
		"",
		http.StatusText(http.StatusRequestTimeout),
		http.StatusRequestTimeout,
		"The request timed out due to rate limiting. Please try again later.",
	)
	ProblemQueueFull = NewProblem(
		"",
		http.StatusText(http.StatusTooManyRequests),
		http.StatusTooManyRequests,
		"The waiting queue is full due to high traffic. Please try again later.",
	)
	ProblemResourceAlreadyExists = NewProblem(
		"",
		http.StatusText(http.StatusConflict),
		http.StatusConflict,
		"The resource or the value of an identifier already exists. Please specify a different identifier and try again.",
	)
	ProblemMissingRequiredParameter = NewProblem(
		"",
		http.StatusText(http.StatusBadRequest),
		http.StatusBadRequest,
		"The resource is missing a required parameter. Check the service documentation and try again.",
	)
)
