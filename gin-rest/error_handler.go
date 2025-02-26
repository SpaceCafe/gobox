package rest

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"strings"

	"github.com/aws/smithy-go"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/mattn/go-sqlite3"
	problems "github.com/spacecafe/gobox/gin-problems"
	"github.com/spacecafe/gobox/gin-rest/types"
	"gorm.io/gorm"
)

// HandleError handles common controller, service, and repository errors and
// maps them to appropriate HTTP responses.
// It checks the error type and calls AbortRequest with the corresponding problem.
func HandleError(ctx *gin.Context, err error) bool { //nolint:gocyclo // Complexity is acceptable for error handling.
	var (
		awsAPIError            smithy.APIError
		jsonUnmarshalTypeError *json.UnmarshalTypeError
		pgError                *pgconn.PgError
		sqlite3Error           sqlite3.Error
		validationErrors       validator.ValidationErrors
	)

	switch {
	case err == nil:
		return false

	// Controller
	case errors.Is(err, io.ErrUnexpectedEOF):
		AbortRequest(ctx, err, problems.ProblemBadRequest.WithDetail(err.Error()))
	case errors.As(err, &jsonUnmarshalTypeError):
		AbortRequest(ctx, err, problems.ProblemBadRequest.WithDetail(jsonUnmarshalTypeError.Field))
	case errors.As(err, &validationErrors):
		handleValidationError(ctx, &validationErrors)

	// Service
	case errors.Is(err, types.ErrNotAuthorized):
		AbortRequest(ctx, err, problems.ProblemInsufficientPermission)
	case errors.Is(err, types.ErrNotFound), errors.Is(err, gorm.ErrRecordNotFound):
		AbortRequest(ctx, err, problems.ProblemNoSuchKey)
	case errors.Is(err, types.ErrDuplicatedKey), errors.Is(err, gorm.ErrDuplicatedKey):
		AbortRequest(ctx, err, problems.ProblemResourceAlreadyExists)
	case errors.Is(err, gorm.ErrCheckConstraintViolated):
		AbortRequest(ctx, err, problems.ProblemMissingRequiredParameter)

	// External libraries
	case errors.As(err, &sqlite3Error):
		handleSqlite3Error(ctx, &sqlite3Error)
	case errors.As(err, &pgError):
		handlePostgreSQLError(ctx, pgError)
	case errors.As(err, &awsAPIError):
		handleAWSAPIError(ctx, awsAPIError)
	default:
		AbortRequest(ctx, err, problems.ProblemInternalError)
	}
	return true
}

// AbortRequest logs and aborts the request with the provided errors.
// It iterates over the errors, logs them to the context, and then aborts the request.
func AbortRequest(ctx *gin.Context, errs ...error) {
	for _, err := range errs {
		if err != nil {
			_ = ctx.Error(err)
		}
	}
	ctx.Abort()
}

func handleValidationError(ctx *gin.Context, errs *validator.ValidationErrors) {
	buff := bytes.NewBufferString("")
	for i := range *errs {
		buff.WriteString((*errs)[i].Field())
		buff.WriteString(" ")
	}
	AbortRequest(ctx, errs, problems.ProblemInvalidArgument.WithDetail(strings.TrimSpace(buff.String())))
}

func handleSqlite3Error(ctx *gin.Context, err *sqlite3.Error) {
	switch err.ExtendedCode {
	case types.SqliteConstraintNotNull:
		AbortRequest(ctx, err, problems.ProblemMissingRequiredParameter)
	case types.SqliteConstraintPrimaryKey:
		AbortRequest(ctx, err, problems.ProblemResourceAlreadyExists)
	default:
		AbortRequest(ctx, err, problems.ProblemInternalError)
	}
}

func handlePostgreSQLError(ctx *gin.Context, err *pgconn.PgError) {
	switch err.Code {
	case types.PostgreSQLNotNullViolation:
		AbortRequest(ctx, err, problems.ProblemMissingRequiredParameter.WithDetail(err.ColumnName))
	case types.PostgreSQLUniqueViolation:
		detail := err.ColumnName
		if detail == "" {
			detail = err.Detail
		}
		AbortRequest(ctx, err, problems.ProblemResourceAlreadyExists.WithDetail(detail))
	default:
		AbortRequest(ctx, err, problems.ProblemInternalError)
	}
}

func handleAWSAPIError(ctx *gin.Context, err smithy.APIError) {
	switch err.ErrorCode() {
	case types.AWSEntityTooLarge:
		AbortRequest(ctx, err, problems.ProblemRequestEntityTooLarge)
	default:
		AbortRequest(ctx, err, problems.ProblemInternalError)
	}
}
