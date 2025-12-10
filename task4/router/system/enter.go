package system

import api "task4/api/v1"

type RouterGroup struct {
	UserRouter
	// 扩展多个
}

var (
	// 定义api模块的接口
	userApi    = api.ApiGroupApp.SystemApiGroup.UserApi
	postApi    = api.ApiGroupApp.SystemApiGroup.PostApi
	commentApi = api.ApiGroupApp.SystemApiGroup.CommentApi
)
