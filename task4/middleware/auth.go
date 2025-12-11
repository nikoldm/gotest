package middleware

import (
	"strings"
	"task4/global"
	"task4/utils"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware 认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取token
		authHeader := utils.GetToken(c)
		if authHeader == "" {
			global.Unauthorized(c, "未登录或非法访问，请登录")
			c.Abort()
			return
		}

		// 检查Bearer前缀
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			global.Unauthorized(c, "Bearer token is required")
			c.Abort()
			return
		}

		// 解析token
		claims, err := utils.ParseToken(tokenString)
		if err != nil {
			global.Unauthorized(c, "Invalid token")
			c.Abort()
			return
		}

		// 将用户信息存储到上下文中
		c.Set("user_id", claims.UserId)
		c.Set("username", claims.Username)
		c.Next()
	}
}
