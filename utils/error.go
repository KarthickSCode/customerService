package utils

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Message string `"json:"error"`
}

// NewError example
func NewError(ctx *gin.Context, status int, err error) {
	fmt.Println("New Error")
	er := HTTPError{
		Code:    status,
		Message: err.Error(),
	}
	ctx.JSON(status, er)
}

// HTTPError example
type HTTPError struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"status bad request"`
}
