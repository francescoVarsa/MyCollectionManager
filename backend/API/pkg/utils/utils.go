package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type SuccessResponse struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}

type FailResponse struct {
	Status  string `json:"status"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func ErrorJSON(err error, ctx *gin.Context, status int) {
	errorCode := map[int]int{
		200: http.StatusOK,
		400: http.StatusBadRequest,
		404: http.StatusNotFound,
		401: http.StatusUnauthorized,
		500: http.StatusInternalServerError,
	}

	errorResponse := FailResponse{
		Status:  "Failed",
		Code:    errorCode[status],
		Message: err.Error(),
	}

	ctx.JSON(errorCode[status], errorResponse)
}
