package system

import "github.com/gin-gonic/gin"

type CommentRouter struct{}

func (s *CommentRouter) InitCommentRouter(Router *gin.RouterGroup, RouterPub *gin.RouterGroup) {

	commentRouter := Router.Group("user")
	commentRouterPub := RouterPub.Group("user")
	{
		commentRouter.POST("", commentApi.CreateComment)

	}
	{
		commentRouterPub.GET("/post/:post_id", commentApi.GetComments)
	}
}
