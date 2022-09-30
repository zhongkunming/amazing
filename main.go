package main

import (
	"service-hub/core"
	"service-hub/module/push"
)

func main() {
	core.InitViper()
	core.InitZap()
	core.InitCron()

	// 注册任务
	//core.Register(daily_check.DailyCheck{})
	//core.Register(db_backup.DbBackup{})
	core.Register(push.Push{})

	core.Server()
}
