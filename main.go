package main

import (
	"service-hub/core"
	"service-hub/module/daily_check"
	"service-hub/module/db_backup"
)

func main() {
	core.InitViper()
	core.InitZap()
	core.InitCron()

	// 注册任务
	core.Register(daily_check.DailyCheck{})
	core.Register(db_backup.DbBackup{})

	core.Server()
}
