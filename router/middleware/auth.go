package middleware

import (
	"github.com/Away0x/api_server/handler"
	"github.com/Away0x/api_server/pkg/errno"
	"github.com/Away0x/api_server/pkg/token"
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
