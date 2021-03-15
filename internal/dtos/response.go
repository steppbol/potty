package dtos

import (
	"github.com/gin-gonic/gin"

	"github.com/steppbol/activity-manager/internal/utils/exception"
)

type response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func CreateJSONResponse(g *gin.Context, httpCode, code int, data interface{}) {
	g.JSON(httpCode, response{
		Code:    code,
		Message: exception.GetMessage(code),
		Data:    data,
	})
}

func CreateBinResponse(g *gin.Context, path string) {
	g.File(path)
}
