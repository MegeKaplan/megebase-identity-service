package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Register() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.String(http.StatusOK, "register endpoint")
	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.String(http.StatusOK, "login endpoint")
	}
}