package system

import "gorm.io/gorm"

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
