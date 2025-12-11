package initialize

import (
	"task4/middleware"
	"task4/router"

	"github.com/gin-gonic/gin"
)

// InitRouters 初始化总路由
func InitRouters() *gin.Engine {
	r := gin.Default()

	// 1. 公共路由组（无需认证）
	publicGroup := r.Group("/api/v1")

	// 2. 私有路由组（需要认证）
	privateGroup := r.Group("/api/v1")
	privateGroup.Use(middleware.AuthMiddleware())

	// 3. 初始化所有模块路由（关键步骤）
	// 文件3的 RouterGroupApp.System 对应文件1的 RouterGroup
	systemRouter := router.RouterGroupApp.System

	// 调用 System 模块下的各个路由组初始化方法
	systemRouter.InitUserRouter(privateGroup, publicGroup)
	systemRouter.InitPostRouter(privateGroup, publicGroup)
	systemRouter.InitCommentRouter(privateGroup, publicGroup)

	return r
}
