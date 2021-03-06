package dtos

import (
	"github.com/gin-gonic/gin"

	"github.com/steppbol/activity-manager/utils/exception"
)

type response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func CreateResponse(g *gin.Context, httpCode, code int, data interface{}) {
	g.JSON(httpCode, response{
		Code:    code,
		Message: exception.GetMessage(code),
		Data:    data,
	})
}
