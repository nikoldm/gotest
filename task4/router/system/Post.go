package system

import "github.com/gin-gonic/gin"

type PostRouter struct{}

func (s *PostRouter) InitPostRouter(Router *gin.RouterGroup, RouterPub *gin.RouterGroup) {

	postRouter := Router.Group("post")
	postRouterPub := RouterPub.Group("post")
	{
		postRouter.POST("", postApi.CreatePost)
		postRouter.PUT("/:id", postApi.UpdatePost)
		postRouter.DELETE("/:id", postApi.DeletePost)

	}
	{
		postRouterPub.GET("/posts", postApi.GetPosts)
		postRouterPub.GET("/posts/:id", postApi.GetPost)
	}
}
