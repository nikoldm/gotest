package system

import (
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
