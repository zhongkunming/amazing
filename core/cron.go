package core

import (
	"amazing/global"
	"github.com/robfig/cron/v3"
)

func InitCron() {
	global.Cron = cron.New(cron.WithParser(cron.NewParser(
		cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.DowOptional | cron.Descriptor,
	)))
	global.Cron.Start()
}
