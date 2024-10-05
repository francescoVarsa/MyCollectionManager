package middleware

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"smart_modellism/pkg/models"
	"smart_modellism/pkg/throttler"
	"smart_modellism/pkg/utils"

	"github.com/gin-gonic/gin"
)

func ValidateCreateModel() gin.HandlerFunc {
	var schema models.NewModel

	return RequestValidation(&schema)
}

func RequestValidation(schema interface{}) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if err := ctx.ShouldBindJSON(&schema); err != nil {
			utils.ErrorJSON(errors.New("the request is not valid"), ctx, http.StatusBadRequest)

			ctx.Abort()

			return
		}

		bodyBytes, err := json.Marshal(schema)
		if err != nil {
			utils.ErrorJSON(errors.New("failed to re-serialize JSON"), ctx, 500)

			ctx.Abort()

			return
		}

		// Rewrap the body into bytes [] to be processed from the endpoint
		ctx.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		ctx.Next()
	}
}

// Middleware per implementare il throttling
func RateLimiterMiddleware(throttler throttler.ThrottleRequests) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ip := ctx.ClientIP()

		if isRequestvalid := throttler.Validate(ip); !isRequestvalid {
			ctx.Abort()
			utils.ErrorJSON(errors.New("too many requests"), ctx, http.StatusTooManyRequests)

			return
		}

		ctx.Next()

	}
}
