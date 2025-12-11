package system

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Comment struct {
	gorm.Model
	Content     string    `gorm:"type:text;not null"`
	CommentTime time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP;not null"`
	UserID      uint      `gorm:"not null;index"` //评论人的id
	User        User      `gorm:"foreignKey:UserID"`
	PostID      uint      `gorm:"not null;index"`    // 外键，指向Post
	Post        Post      `gorm:"foreignKey:PostID"` //属于关系
}

// AfterCreate 更新Comment模型定义，添加钩子函数
// Comment的钩子函数：创建后更新文章评论数量
func (c *Comment) AfterCreate(tx *gorm.DB) (err error) {
	if err = tx.Model(&Post{}).Where("id=?", c.PostID).UpdateColumn("comment_count", gorm.Expr("comment_count+1")).UpdateColumn("comment_status", "有评论").Error; err != nil {
		return fmt.Errorf("更新文章评论数量失败: %w", err.Error)
	}
	fmt.Printf("评论创建成功，已更新文章%d的评论数量\n", c.PostID)
	return nil
}

// AfterDelete 钩子函数：删除时检查文章评论数量
func (c *Comment) AfterDelete(tx *gorm.DB) (err error) {
	// 删除后，获取评论数量
	var commentCount int64
	commentStatus := "有评论"
	tx.Model(&Comment{}).Where("Post_id=?", c.PostID).Count(&commentCount)
	if commentCount <= 0 {
		// 更新文章的评论数量
		commentStatus = "无评论"
	}
	tx.Model(&Post{}).Where("ID=?", c.PostID).UpdateColumn("comment_count", commentCount).UpdateColumn("comment_status", commentStatus)
	fmt.Printf("评论删除成功，文章%d还有%d条评论\n", c.PostID, commentCount)
	return nil
}
