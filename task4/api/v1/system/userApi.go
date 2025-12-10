package system

import (
	"fmt"
	"task4/model/system"

	"github.com/gin-gonic/gin"
)

type UserApi struct{}

type RegisterReq struct {
	Username string `json:"username" binding:"required,min=3,max=20"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AuthResp struct {
	Token string      `json:"token"`
	User  system.User `json:"user"`
}

func (UserApi *UserApi) Register(c *gin.Context) {
	var req RegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		return
	}
	fmt.Println(req.Username)
}

func (UserApi *UserApi) Login(c *gin.Context) {
	var req LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		return
	}
	fmt.Println(req.Username)
}

// GetUserInfo 获取用户信息
func (UserApi *UserApi) GetUserInfo(c *gin.Context) {
	userId, err := c.Get("user_id")
	if !err {
		return
	}
	fmt.Println(userId)
}
