package system

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username  string `gorm:"type:varchar(32);not null"`
	Password  string `gorm:"type:varchar(255);not null"`
	Email     string `gorm:"type:varchar(100);uniqueIndex;not null"`
	Age       uint8  `gorm:"check:age>0"`
	PostCount uint   `gorm:"default:0"`
	Posts     []Post `gorm:"foreignKey:UserID"` // 一对多：一个用户多篇文章
}
