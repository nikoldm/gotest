package system

import (
	"fmt"
	"task4/config"
	"task4/model/system"
)

type CommentService struct{}

func (comm *CommentService) CreateComment(comment *system.Comment) error {
	// 检查文章是否存在
	var post system.Post
	if err := config.DB.First(&post, comment.PostID).Error; err != nil {
		return fmt.Errorf("文章不存在：" + err.Error())
	}

	if err := config.DB.Create(&comment).Error; err != nil {
		return err
	}

	// 预加载用户信息
	config.DB.Preload("User").First(comment, comment.ID)
	return nil
}

// GetComments 获取文章评论
func (comm *CommentService) GetComments(postID uint64, page int, pageSize int, total *int64) (comments []system.Comment, err error) {
	// 检查文章是否存在
	var post system.Post
	if err := config.DB.First(&post, postID).Error; err != nil {
		return comments, fmt.Errorf("文章不存在")
	}

	// 查询评论列表，预加载用户信息
	offset := (page - 1) * pageSize
	if err := config.DB.Preload("User").
		Where("post_id = ?", postID).
		Order("created_at ASC").
		Limit(pageSize).
		Offset(offset).
		Find(&comments).Error; err != nil {
		return comments, fmt.Errorf("查询失败：" + err.Error())
	}
	config.DB.Model(&system.Comment{}).Where("post_id = ?", postID).Count(total)
	return comments, nil
}
