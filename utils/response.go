package utils

import (
	"github.com/MegeKaplan/megebase-identity-service/utils/response"
	"github.com/gin-gonic/gin"
)

func JSONError(c *gin.Context, err *response.AppError, details string) {
	if details != "" {
		err.Details = details
	}
	c.JSON(err.HTTPStatus, err)
	c.Abort()
}

func JSONSuccess(c *gin.Context, resp *response.AppSuccess, data interface{}) {
	if data != nil {
		resp.Data = data
	}
	c.JSON(resp.HTTPStatus, resp)
}
