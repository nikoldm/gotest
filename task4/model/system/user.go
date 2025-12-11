package system

import (
	"task4/utils"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username  string `gorm:"type:varchar(32);uniqueIndex;not null"`
	Password  string `gorm:"type:varchar(255);not null"`
	Email     string `gorm:"type:varchar(100);uniqueIndex;not null"`
	Age       uint8  `gorm:"default:0;check:age>=0"`
	PostCount uint   `gorm:"default:0"`
	Posts     []Post `gorm:"foreignKey:UserID"` // 一对多：一个用户多篇文章
}

// BeforeCreate GORM钩子，在创建用户前自动哈希密码
func (u *User) BeforeCreate(tx *gorm.DB) error {
	u.Password = utils.BcryptHash(u.Password)
	return nil
}
