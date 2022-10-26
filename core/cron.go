package core

import (
	"github.com/robfig/cron/v3"
	"mcs/global"
)

func InitCron() {
	global.Cron = cron.New(cron.WithParser(cron.NewParser(
		cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.DowOptional | cron.Descriptor,
	)))
	global.Cron.Start()
}
