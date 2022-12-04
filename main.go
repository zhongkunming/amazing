package main

import (
	"amazing/core"
	"amazing/module/flows"
)

func main() {
	core.InitViper()
	core.InitZap()
	core.InitCron()

	// 注册任务
	core.Register(flows.LoadableFlows{})

	core.Server()
}
