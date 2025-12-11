package system

import (
	"strconv"
	"task4/global"
	"task4/model/system"

	"github.com/gin-gonic/gin"
)

type CommentApi struct{}

type CommentReq struct {
	Content string `json:"content" binding:"required,min=1,max=1000"`
}

func (commentApi *CommentApi) CreateComment(c *gin.Context) {
	postID, reqErr := strconv.ParseUint(c.Param("post_id"), 10, 32)
	if reqErr != nil {
		global.BadRequest(c, "文章ID必传")
		return
	}
	var req CommentReq
	if err := c.ShouldBindJSON(&req); err != nil {
		global.BadRequest(c, err.Error())
		return
	}
	userID := c.MustGet("user_id")
	comment := system.Comment{
		Content: req.Content,
		UserID:  userID.(uint),
		PostID:  uint(postID),
	}
	if err := commentService.CreateComment(&comment); err != nil {
		global.BadRequest(c, "评论创建失败:"+err.Error())
	}

	global.Success(c, comment)
}

func (commentApi *CommentApi) GetComments(c *gin.Context) {
	postID, reqErr := strconv.ParseUint(c.Param("post_id"), 10, 32)
	if reqErr != nil {
		global.BadRequest(c, "postID必传："+reqErr.Error())
		return
	}

	// 分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var total int64
	comments, err := commentService.GetComments(postID, page, pageSize, &total)

	if err != nil {
		global.NotFound(c, "查询失败："+err.Error())
	}

	global.Success(c, gin.H{
		"comments":  comments,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}
