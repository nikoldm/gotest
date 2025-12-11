package system

import (
	"errors"
	"fmt"
	"task4/config"
	"task4/model/system"
	"task4/utils"

	"gorm.io/gorm"
)

type UserService struct{}

func (UserService *UserService) Register(u *system.User) (userInter *system.User, err error) {
	var user system.User
	if !errors.Is(config.DB.Where("username = ?", u.Username).Or("email=?", u.Email).First(&user).Error, gorm.ErrRecordNotFound) { // 判断用户名是否注册
		return userInter, errors.New("用户名或邮箱已注册")
	}
	// 钩子函数对密码hash加密 注册
	err = config.DB.Create(&u).Error
	return u, err
}

func (UserService *UserService) Login(u *system.User) (userInter *system.User, err error) {
	if nil == config.DB {
		return nil, fmt.Errorf("db not init")
	}

	var user system.User
	err = config.DB.Where("username = ?", u.Username).First(&user).Error
	if err == nil {
		if ok := utils.BcryptCheck(u.Password, user.Password); !ok {
			return nil, errors.New("密码错误")
		}
	}
	return &user, err
}

func (UserService *UserService) GetUserInfo(userId any) (u *system.User, err error) {
	var user system.User
	if err = config.DB.First(&user, userId).Error; err != nil {
		return &user, err
	}
	return &user, err
}
