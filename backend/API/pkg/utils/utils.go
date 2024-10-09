package utils

import (
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

	errorResponse := FailResponse{
		Status:  "Failed",
		Code:    status,
		Message: err.Error(),
	}

	ctx.JSON(status, errorResponse)
}
