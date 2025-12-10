package system

import "github.com/gin-gonic/gin"

type CommentApi struct{}

type CommentReq struct {
	Content string `json:"content" binding:"required,min=1,max=1000"`
}

func (comment *CommentApi) CreateComment(c *gin.Context) {

}

func (comment *CommentApi) GetComments(c *gin.Context) {

}
