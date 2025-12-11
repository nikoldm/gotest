package system

import (
	"task4/global"
	"task4/model/system"
	"task4/utils"

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
	var r RegisterReq
	if err := c.ShouldBindJSON(&r); err != nil {
		global.BadRequest(c, err.Error())
		return
	}
	user := &system.User{Username: r.Username, Password: r.Password, Email: r.Email}
	userReturn, err := userService.Register(user)
	if err != nil {
		global.InternalServerError(c, err.Error())
		return
	}
	token, err := utils.CreateToken(userReturn.ID, userReturn.Username)
	if err != nil {
		global.InternalServerError(c, "Failed to generate token")
		return
	}
	global.Success(c, AuthResp{
		Token: token,
		User:  *userReturn,
	})
}

func (UserApi *UserApi) Login(c *gin.Context) {
	var req LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		global.BadRequest(c, err.Error())
		return
	}
	u := &system.User{Username: req.Username, Password: req.Password}
	user, err := userService.Login(u)
	if err != nil {
		global.Unauthorized(c, "Invalid username or password")
		return
	}
	// 生成JWT token
	token, err := utils.CreateToken(user.ID, user.Username)
	if err != nil {
		global.InternalServerError(c, "Failed to generate token")
		return
	}
	// 放入cookie
	utils.SetToken(c, "Bearer "+token, 60*60*1000)

	global.Success(c, AuthResp{
		Token: token,
		User:  *user,
	})
}

// GetUserInfo 获取用户信息
func (UserApi *UserApi) GetUserInfo(c *gin.Context) {
	userId, exists := c.Get("user_id")
	if !exists {
		return
	}
	user, err := userService.GetUserInfo(userId)
	if err != nil {
		global.NotFound(c, "User not found")
		return
	}
	global.Success(c, user)
}
