package system

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type PostApi struct{}

type PostReq struct {
	Title   string `json:"title" binding:"required,min=1,max=200"`
	Content string `json:"content" binding:"required,min=1"`
}

func (p *PostApi) CreatePost(c *gin.Context) {
	var postReq PostReq
	if err := c.ShouldBindJSON(&postReq); err != nil {
		return
	}
}

func (p *PostApi) UpdatePost(c *gin.Context) {
	fmt.Println("post update")
}

func (p *PostApi) DeletePost(c *gin.Context) {

}

func (p *PostApi) GetPosts(c *gin.Context) {
	if err := c.ShouldBindJSON(&PostReq{}); err != nil {
		return
	}
	fmt.Println("sdfsdfsdffff======")
}

func (p *PostApi) GetPost(c *gin.Context) {

}
