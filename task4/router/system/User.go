package system

import "github.com/gin-gonic/gin"

type UserRouter struct{}

func (s *UserRouter) InitUserRouter(Router *gin.RouterGroup, RouterPub *gin.RouterGroup) {

	userRouter := Router.Group("user")
	userRouterPub := RouterPub.Group("user")
	{
		userRouter.GET("/getUserInfo", userApi.GetUserInfo)

	}
	{
		userRouterPub.POST("/register", userApi.Register)
		userRouterPub.POST("/login", userApi.Login)
	}
}
