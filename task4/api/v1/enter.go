package v1

import "task4/api/v1/system"

var ApiGroupApp = new(ApiGroup)

type ApiGroup struct {
	SystemApiGroup system.ApiGroup

	// …… 定义多个模块
}
