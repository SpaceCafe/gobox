package problems

import (
	"errors"

	"github.com/gin-gonic/gin"
)

func New() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()

		lastError := ctx.Errors.Last()
		if lastError == nil {
			return
		}

		var p *Problem
		if errors.As(lastError.Err, &p) {
			// Set the request path to the instance field of the problem.
			if len(p.Instance) == 0 {
				tmp := *p
				tmp.Instance = ctx.Request.URL.Path
				p = &tmp
			}

			ctx.Writer.Header().Set("Content-Type", "application/problem+json")
			ctx.AbortWithStatusJSON(p.Status, p)
		}
	}
}
