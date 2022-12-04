package core

import (
	"amazing/global"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	sysLog "log"
)

func InitViper() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		sysLog.Fatalf("读取配置文件失败: %v", err)
	}
	viper.OnConfigChange(func(e fsnotify.Event) {
		err = viper.Unmarshal(&global.Global)
		if err != nil {
			sysLog.Fatalf("加载配置文件失败: %v", err)
		}
		// todo reload config
	})
	viper.WatchConfig()
	err = viper.Unmarshal(&global.Global)

	if err != nil {
		sysLog.Fatalf("加载配置文件失败: %v", err)
	}
}
