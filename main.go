package main

import (
	"mcs/core"
	"mcs/module/flows"
)

func main() {
	core.InitViper()
	core.InitZap()
	core.InitCron()

	// 注册任务
	core.Register(flows.LoadableFLows{})

	core.Server()
}
