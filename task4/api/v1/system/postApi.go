package system

import (
	"strconv"
	"task4/global"
	"task4/model/system"

	"github.com/gin-gonic/gin"
)

type PostApi struct{}

type PostReq struct {
	Title   string `json:"title" binding:"required,min=1,max=200"`
	Content string `json:"content" binding:"required,min=1"`
}

// CreatePost 创建文章
func (p *PostApi) CreatePost(c *gin.Context) {
	var req PostReq
	if err := c.ShouldBindJSON(&req); err != nil {
		global.BadRequest(c, err.Error())
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		global.Unauthorized(c, "用户未登录")
		return
	}
	post := system.Post{
		Title:       req.Title,
		PostContent: req.Content,
		UserID:      userID.(uint),
	}
	createPost, err := postService.CreatePost(&post)
	if err != nil {
		global.BadRequest(c, "创建失败："+err.Error())
		return
	}

	global.Success(c, createPost)
}

// UpdatePost 更新文章
func (p *PostApi) UpdatePost(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		global.BadRequest(c, "postID必传："+err.Error())
		return
	}
	var req PostReq
	if err = c.ShouldBindJSON(&req); err != nil {
		global.BadRequest(c, err.Error())
		return
	}

	userID := c.MustGet("user_id")
	post := system.Post{
		Title:       req.Title,
		PostContent: req.Content,
		UserID:      userID.(uint),
	}
	post.ID = uint(postID)

	err = postService.UpdatePost(&post)
	if err != nil {
		global.BadRequest(c, "post更新失败："+err.Error())
		return
	}

	global.Success(c, post)
}

// DeletePost 删除文章
func (p *PostApi) DeletePost(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		global.BadRequest(c, "postID必传："+err.Error())
		return
	}
	userID := c.MustGet("user_id")
	postService.DeletePost(postID, userID)
	global.Success(c, gin.H{"message": "文章删除成功"})
}

// GetPosts 获取文章列表
func (p *PostApi) GetPosts(c *gin.Context) {
	// 分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	// 获取总数
	var total int64
	// 查询文章列表，预加载用户信息
	posts, err := postService.GetPosts(page, pageSize, &total)
	if err != nil {
		global.BadRequest(c, "查询出错："+err.Error())
	}
	global.Success(c, gin.H{
		"posts":     posts,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// GetPost 获取文章详情
func (p *PostApi) GetPost(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		global.BadRequest(c, "postID必传："+err.Error())
		return
	}

	post, err := postService.GetPost(postID)
	if err != nil {
		global.BadRequest(c, err.Error())
		return
	}
	global.Success(c, post)
}
