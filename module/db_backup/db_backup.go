package db_backup

import (
	"database/sql"
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

			bDbConfig := global.Global.DbBackup.BDb
			bDbUrl := fmt.Sprintf("%s:%s@tcp(%s)/%s", bDbConfig.Username, bDbConfig.Passwd, bDbConfig.Host, bDbConfig.Database)
			bDb, _ := sql.Open("mysql", bDbUrl)
			defer bDb.Close()

			sDbConfig := global.Global.DbBackup.SDb
			sDbUrl := fmt.Sprintf("%s:%s@tcp(%s)/%s", sDbConfig.Username, sDbConfig.Passwd, sDbConfig.Host, sDbConfig.Database)
			sDb, _ := sql.Open("mysql", sDbUrl)
			defer sDb.Close()

			body{bDb: bDb, sDb: sDb}.run()
		}
		go processBody()
		waitGroup.Wait()
		global.Log.Infof("%s, 备份完成", time.Now().Format("2006-01-02"))
	})
	if err != nil {
		global.Log.Fatal(err)
	}
}
