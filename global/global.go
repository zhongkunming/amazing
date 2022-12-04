package global

import (
	"amazing/config"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

var (
	Global *config.Config
	Log    *zap.SugaredLogger
	Cron   *cron.Cron
)
