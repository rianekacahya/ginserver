package middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/rianekacahya/ginserver/response"
)

var(
	errorString = errors.New("Internal Server Error")
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				response.Error(c, errorString)
			}
		}()

		c.Next()
	}
}
