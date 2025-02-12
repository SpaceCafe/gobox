package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"mime"
	"sort"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/mattn/go-sqlite3"
	"github.com/spacecafe/gobox/gin-authorization"
	"github.com/spacecafe/gobox/gin-jwt"
	"github.com/spacecafe/gobox/gin-problems"
	"github.com/spacecafe/gobox/gin-rest/types"
	"gorm.io/gorm"
)

const (
	// AcceptSliceCapacity contains the initial capacity for slices to store MIME types and weights.
	AcceptSliceCapacity = 10

	// AcceptQualityWeight corresponds the default quality weight for MIME types.
	AcceptQualityWeight = 1.0

	// AcceptQualityParameter is used to specify quality weight parameter in the header.
	AcceptQualityParameter = ";q="

	// AcceptSeparator is the used Separator for multiple MIME types in the header.
	AcceptSeparator = ','

	// AcceptHeader is the HTTP header key for the Accept header.
	AcceptHeader = "Accept"
)

// GetView retrieves the appropriate view for a given resource based on the client's Accept header.
// It iterates through the MIME types specified in the Accept header and returns the corresponding view if found.
// If no matching view is found, it handles the error by responding with an unsupported media type problem.
func GetView[T any](ctx *gin.Context, resource types.Resource[T]) any {
	for _, mimeType := range ParseAcceptHeader(ctx.GetHeader(AcceptHeader)) {
		if view, ok := resource.GetViews()[mimeType]; ok {
			return view
		}
	}
	HandleError(ctx, problems.ProblemUnsupportedMediaType)
	return nil
}

// NewServiceOptions creates a new instance of ServiceOptions.
// It extracts the username from JWT claims and retrieves authorizations
// from the context, then returns a pointer to a ServiceOptions struct.
func NewServiceOptions(ctx *gin.Context) *types.ServiceOptions {
	subject, err := jwt.GetClaims(ctx).GetSubject()
	if err != nil {
		panic(err)
	}
	return &types.ServiceOptions{
		UserID:         subject,
		Authorizations: authorization.GetAuthorizations(ctx),
	}
}

// ParseAcceptHeader parses the Accept header from an HTTP request.
// It returns a list of supported MIME types sorted by their quality weights in descending order.
func ParseAcceptHeader(acceptHeader string) []string {
	// Preallocate slices with a reasonable initial capacity.
	mimeTypes := make([]string, 0, AcceptSliceCapacity)
	weights := make([]float64, 0, AcceptSliceCapacity)
	startIndex := 0

	// Loop multiple accepted mimetypes.
	for i := 0; i <= len(acceptHeader); i++ {
		if i == len(acceptHeader) || acceptHeader[i] == AcceptSeparator {
			headerSegment := strings.TrimSpace(acceptHeader[startIndex:i])
			mimeType := headerSegment
			weight := AcceptQualityWeight

			// Receive quality weight, if set.
			if k := strings.Index(headerSegment, AcceptQualityParameter); k != -1 {
				mimeType = headerSegment[:k]
				if value, err := strconv.ParseFloat(headerSegment[k+3:], 64); err == nil {
					weight = value
				}
			}

			// Add mimetype to list, if supported by caller.
			if value, _, err := mime.ParseMediaType(mimeType); err == nil {
				mimeTypes = append(mimeTypes, value)
				weights = append(weights, weight)
			}

			startIndex = i + 1
		}
	}

	// Sort by weight in descending order using a custom sorting algorithm.
	sort.Slice(mimeTypes, func(i, j int) bool {
		return weights[i] > weights[j]
	})

	return mimeTypes
}

// HandleError logs and aborts the request with the provided errors.
// It iterates over the errors, logs them to the context, and then aborts the request.
func HandleError(ctx *gin.Context, errs ...error) {
	for _, err := range errs {
		if err != nil {
			_ = ctx.Error(err)
		}
	}
	ctx.Abort()
}

// HandleControllerError handles common controller errors and maps them to appropriate HTTP responses.
// It checks the error type and calls HandleError with the corresponding problem.
func HandleControllerError(ctx *gin.Context, err error) bool {
	var jsonUnmarshalTypeError *json.UnmarshalTypeError
	var validationErrors validator.ValidationErrors

	switch {
	case err == nil:
		return false
	case errors.As(err, &jsonUnmarshalTypeError):
		HandleError(ctx, err, problems.ProblemBadRequest.WithDetail(jsonUnmarshalTypeError.Field))
	case errors.As(err, &validationErrors):
		buff := bytes.NewBufferString("")
		for i, _ := range validationErrors {
			buff.WriteString(validationErrors[i].Field())
			buff.WriteString(" ")
		}
		HandleError(ctx, err, problems.ProblemInvalidArgument.WithDetail(strings.TrimSpace(buff.String())))
	default:
		HandleError(ctx, err, problems.ProblemBadRequest)
	}
	return true
}

// HandleServiceError handles common service errors and maps them to appropriate HTTP responses.
// It checks the error type and calls HandleError with the corresponding problem.
func HandleServiceError(ctx *gin.Context, err error) bool {
	var pgError *pgconn.PgError
	var sqliteError sqlite3.Error

	switch {
	case err == nil:
		return false
	case errors.Is(err, types.ErrNotAuthorized):
		HandleError(ctx, err, problems.ProblemInsufficientPermission)
	case errors.Is(err, types.ErrNotFound), errors.Is(err, gorm.ErrRecordNotFound):
		HandleError(ctx, err, problems.ProblemNoSuchKey)
	case errors.Is(err, types.ErrDuplicatedKey), errors.Is(err, gorm.ErrDuplicatedKey):
		HandleError(ctx, err, problems.ProblemResourceAlreadyExists)
	case errors.Is(err, gorm.ErrCheckConstraintViolated):
		HandleError(ctx, err, problems.ProblemMissingRequiredParameter)

	// Sqlite 3
	case errors.As(err, &sqliteError):
		switch sqliteError.ExtendedCode {
		case types.SqliteConstraintNotNull:
			HandleError(ctx, err, problems.ProblemMissingRequiredParameter)
		case types.SqliteConstraintPrimaryKey:
			HandleError(ctx, err, problems.ProblemResourceAlreadyExists)
		default:
			HandleError(ctx, err, problems.ProblemInternalError)
		}

	// PostgreSQL
	case errors.As(err, &pgError):
		switch pgError.Code {
		case "23502":
			HandleError(ctx, err, problems.ProblemMissingRequiredParameter.WithDetail(pgError.ColumnName))
		case "23505":
			HandleError(ctx, err, problems.ProblemResourceAlreadyExists.WithDetail(pgError.ColumnName))
		default:
			HandleError(ctx, err, problems.ProblemInternalError)
		}

	default:
		HandleError(ctx, err, problems.ProblemInternalError)
	}
	return true
}
