package system

import "task4/service"

type ApiGroup struct {
	UserApi
	PostApi
	CommentApi
}

var (
	// 定义service层的接口
	userService    = service.ServiceGroupApp.SystemServiceGroup.UserService
	postService    = service.ServiceGroupApp.SystemServiceGroup.PostService
	commentService = service.ServiceGroupApp.SystemServiceGroup.CommentService
)
