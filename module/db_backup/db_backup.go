package db_backup

import (
	"fmt"
	"service-hub/global"
	"sync"
	"time"
)

const spec = "*/10 * * * * ?"

type DbBackup struct {
}

func (r DbBackup) Description() string {
	return fmt.Sprint("数据库备份")
}

func (r DbBackup) Load() {
	_, err := global.Cron.AddFunc(spec, func() {
		waitGroup := sync.WaitGroup{}
		global.Log.Infof("%s, 开始备份数据库", time.Now().Format("2006-01-02"))
		waitGroup.Add(1)
		processBody := func() {
			defer waitGroup.Done()
			body{}.run()
		}
		go processBody()
		waitGroup.Wait()
		global.Log.Infof("%s, 备份完成", time.Now().Format("2006-01-02"))
	})
	if err != nil {
		global.Log.Fatal(err)
	}
}
