package middleware

import (
	"api_server/handler"
	"api_server/pkg/errno"
	"api_server/pkg/token"
	"github.com/gin-gonic/gin"
)

// 验证 token
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Parse the json web token.
		if _, err := token.ParseRequest(c); err != nil {
			handler.SendResponse(c, errno.ErrTokenInvalid, nil)
			c.Abort()
			return
		}

		c.Next()
	}
}
