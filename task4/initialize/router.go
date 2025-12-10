package initialize

import (
	"task4/api/v1/system"
	"task4/middleware"

	"github.com/gin-gonic/gin"
)

// 初始化总路由
// 路由管理器
type RouterManager struct {
	router     *gin.Engine
	userApi    *system.UserApi
	postApi    *system.PostApi
	commentApi *system.CommentApi
}

func NewRouterManager() *RouterManager {
	return &RouterManager{
		router:     gin.New(),
		userApi:    &system.UserApi{},
		postApi:    &system.PostApi{},
		commentApi: &system.CommentApi{},
	}
}

func (rm *RouterManager) Setup() *gin.Engine {
	// 全局中间件
	rm.router.Use(middleware.LoggerMiddleware())
	rm.router.Use(middleware.ErrorLoggerMiddleware())
	rm.router.Use(gin.Recovery())

	// 设置路由
	rm.setupRoutes()

	return rm.router
}

func (rm *RouterManager) setupRoutes() {
	// API v1
	apiV1 := rm.router.Group("/api/v1")

	// 各模块路由
	rm.setupAuthRoutes(apiV1)
	rm.setupAuthenticatedRoutes(apiV1)
	rm.setupPublicRoutes(apiV1)

	// 健康检查
	rm.setupHealthCheck()
}

func (rm *RouterManager) setupAuthRoutes(api *gin.RouterGroup) {
	auth := api.Group("/auth")
	auth.POST("/register", rm.userApi.Register)
	auth.POST("/login", rm.userApi.Login)
}

func (rm *RouterManager) setupAuthenticatedRoutes(api *gin.RouterGroup) {
	needAuth := api.Group("")
	needAuth.Use(middleware.AuthMiddleware())

	// 用户信息
	needAuth.GET("/profile", rm.userApi.GetUserInfo)

	// 文章管理
	rm.setupPostManagementRoutes(needAuth)

	// 评论管理
	rm.setupCommentManagementRoutes(needAuth)
}

func (rm *RouterManager) setupPostManagementRoutes(authGroup *gin.RouterGroup) {
	posts := authGroup.Group("/posts")
	posts.POST("", rm.postApi.CreatePost)
	posts.PUT("/:id", rm.postApi.UpdatePost)
	posts.DELETE("/:id", rm.postApi.DeletePost)
}

func (rm *RouterManager) setupCommentManagementRoutes(authGroup *gin.RouterGroup) {
	comments := authGroup.Group("/posts/:post_id/comments")
	comments.POST("", rm.commentApi.CreateComment)
}

func (rm *RouterManager) setupPublicRoutes(api *gin.RouterGroup) {
	public := api.Group("")
	public.GET("/posts", rm.postApi.GetPosts)
	public.GET("/posts/:id", rm.postApi.GetPost)

	// 评论公开路由（单独分组）
	comments := api.Group("/comments")
	comments.GET("/post/:post_id", rm.commentApi.GetComments)
}

func (rm *RouterManager) setupHealthCheck() {
	rm.router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "pong",
		})
	})
}

func Routers() *gin.Engine {
	rm := NewRouterManager()
	return rm.Setup()
}
