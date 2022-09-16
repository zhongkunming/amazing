package daily_check

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"service-hub/config"
	"service-hub/global"
	"sync"
	"time"
)

const spec = "0 0 8 * * ?"

type DailyCheck struct{}

func (r DailyCheck) Description() string {
	return fmt.Sprint("签到")
}

func (r DailyCheck) Load() {
	_, err := global.Cron.AddFunc(spec, func() {
		waitGroup := sync.WaitGroup{}
		global.Log.Infof("%s, 开始签到", time.Now().Format("2006-01-02"))
		waitGroup.Add(len(global.Global.DailyCheckConfig.Users))
		for _, user := range global.Global.DailyCheckConfig.Users {
			processBody := func(u config.DailyCheckUsersConfig) {
				defer waitGroup.Done()
				body{user: u, client: resty.New().R()}.do()
			}
			go processBody(user)
		}
		waitGroup.Wait()
		global.Log.Infof("%s, 签到完成", time.Now().Format("2006-01-02"))
	})
	if err != nil {
		global.Log.Fatal(err)
	}
}
