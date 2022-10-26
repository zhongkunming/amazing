package global

import (
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"mcs/config"
)

var (
	Global *config.Config
	Log    *zap.SugaredLogger
	Cron   *cron.Cron
)
