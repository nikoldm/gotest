package system

import (
	"fmt"

	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	Title         string    `gorm:"type:varchar(64);not null"`
	PostContent   string    `gorm:"type:text;not null"`
	CommentCount  uint      `gorm:"default:0"`
	Status        string    `gorm:"type:varchar(32);default:'published';not null'"`
	CommentStatus string    `gorm:"type:varchar(32);default:'无评论';not null'"`
	Comments      []Comment `gorm:"foreignKey:PostID"` // 一对多：一篇文章多个评论
	User          User      `gorm:"foreignKey:UserID"` // 属于关系
	UserID        uint      `gorm:"not null;index"`    // 外键指向User
}

// AfterCreate Post的钩子函数：创建后更新 用户的 文章数量
func (p *Post) AfterCreate(tx *gorm.DB) (err error) {
	// 后更新 用户的 文章数量
	if err = tx.Model(&User{}).Where("ID=?", p.UserID).UpdateColumn("post_count", gorm.Expr("post_count+1")).Error; err != nil {
		return fmt.Errorf("更新用户文章数量失败: %w", err.Error())
	}

	fmt.Printf("文章创建成功，已更新用户%d的文章数量\n", p.UserID)
	return nil
}

// BeforeDelete Post的钩子函数：删除前减少用户文章数量
func (p *Post) BeforeDelete(tx *gorm.DB) (err error) {
	if err = tx.First(&p, p.ID).Error; err != nil {
		return fmt.Errorf("文章不存在：%w", err.Error())
	}
	if err = tx.Model(&User{}).Where("ID=?", p.UserID).UpdateColumn("post_count", gorm.Expr("post_count-1")).Error; err != nil {
		return fmt.Errorf("更新用户文章数量失败: %w", err.Error())
	}
	fmt.Printf("文章删除前，已更新用户%d的文章数量\n", p.UserID)
	return nil
}
