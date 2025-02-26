package rest

import (
	"mime"
	"sort"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	problems "github.com/spacecafe/gobox/gin-problems"
	"github.com/spacecafe/gobox/gin-rest/types"
)

// AcceptMiddleware is a Gin middleware that checks if the request's Accept header contains any of the supported MIME types.
// If no supported MIME type is found, it aborts the request with a ProblemUnsupportedMediaType error.
func AcceptMiddleware(supportedMimetypes map[string]any) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		for _, mimetype := range ParseAcceptHeader(ctx.GetHeader(types.AcceptHeader)) {
			if _, ok := supportedMimetypes[mimetype]; ok {
				ctx.Set(types.ContextDataRenderMimetype, mimetype)
				ctx.Next()
				return
			}
		}
		problems.ProblemUnsupportedMediaType.Abort(ctx)
	}
}

// ParseAcceptHeader parses the Accept header from an HTTP request.
// It returns a list of supported MIME types sorted by their quality weights in descending order.
func ParseAcceptHeader(acceptHeader string) []string {
	// Preallocate slices with a reasonable initial capacity.
	mimeTypes := make([]string, 0, types.AcceptSliceCapacity)
	weights := make([]float64, 0, types.AcceptSliceCapacity)

	for _, segment := range strings.Split(acceptHeader, types.AcceptSeparator) {
		if mediaType, params, err := mime.ParseMediaType(segment); err == nil {
			mimeTypes = append(mimeTypes, mediaType)

			if weightStr, ok := params[types.AcceptQualityParameter]; ok {
				if weight, err := strconv.ParseFloat(weightStr, 64); err == nil {
					weights = append(weights, weight)
					continue
				}
			}
			weights = append(weights, types.AcceptQualityWeight)
		}
	}

	// Sort by weight in descending order using a custom sorting algorithm.
	//nolint:gocritic // It's not necessary to use the mimeTypes slice for sorting.
	sort.Slice(mimeTypes, func(i, j int) bool {
		return weights[i] > weights[j]
	})

	return mimeTypes
}
