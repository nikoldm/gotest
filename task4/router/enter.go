package router

import "task4/router/system"

var RouterGroupApp = new(RouterGroup)

type RouterGroup struct {
	// 系统模块路由组
	System system.RouterGroup
	// …… 各个模块
}
