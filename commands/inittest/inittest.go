package main

import (
	"github.com/talkincode/logsight/app"
	"github.com/talkincode/logsight/config"
)

func main() {
	app.InitGlobalApplication(config.LoadConfig("../logsight.yml"))
	app.GApp().InitTest()
}
