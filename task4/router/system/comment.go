package system

import "github.com/gin-gonic/gin"

type CommentRouter struct{}

func (s *CommentRouter) InitCommentRouter(Router *gin.RouterGroup, RouterPub *gin.RouterGroup) {

	commentRouter := Router.Group("comment")
	commentRouterPub := RouterPub.Group("comment")
	{
		commentRouter.POST("/post/:post_id", commentApi.CreateComment)

	}
	{
		commentRouterPub.GET("/post/:post_id", commentApi.GetComments)
	}
}
