package system

import (
	"errors"
	"task4/config"
	"task4/model/system"
	"task4/utils"

	"gorm.io/gorm"
)

type UserService struct{}

func (userService *UserService) Register(u system.User) (userInter system.User, err error) {
	var user system.User
	if !errors.Is(config.DB.Where("username = ?", u.Username).First(&user).Error, gorm.ErrRecordNotFound) { // 判断用户名是否注册
		return userInter, errors.New("用户名已注册")
	}
	// 否则 附加uuid 密码hash加密 注册
	u.Password = utils.BcryptHash(u.Password)
	//u.UUID = uuid.New()
	err = config.DB.Create(&u).Error
	return u, err
}
