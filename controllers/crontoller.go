package controllers

import (
	"github.com/talkincode/logsight/controllers/dashboard"
	"github.com/talkincode/logsight/controllers/index"
	"github.com/talkincode/logsight/controllers/logs"
	"github.com/talkincode/logsight/controllers/metrics"
	"github.com/talkincode/logsight/controllers/opr"
	"github.com/talkincode/logsight/controllers/radius"
	"github.com/talkincode/logsight/controllers/settings"
)

// Init web 控制器初始化
func Init() {
	index.InitRouter()
	opr.InitRouter()
	settings.InitRouter()
	dashboard.InitRouter()
	logs.InitRouter()
	radius.InitRouter()
	metrics.InitRouter()
}
