package main

import (
	"service-hub/core"
	"service-hub/module/check"
)

func main() {
	core.InitViper()
	core.InitZap()
	core.InitCron()

	// 注册签到任务
	core.Register(check.DailyCheck{})

	core.Server()
}
