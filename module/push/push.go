package push

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"service-hub/global"
	"sync"
)

const spec = "0 */1 * * * ?"

type Push struct{}

func (r Push) Description() string {
	return fmt.Sprint("access_token刷新")
}

func (r Push) Load() {
	go processBody(global.Global.Push.AppId, global.Global.Push.AppSecret)

	_, err := global.Cron.AddFunc(spec, func() {
		global.Log.Info("access_token开始刷新")
		var waitGroup = sync.WaitGroup{}
		waitGroup.Add(1)
		process := func() {
			defer waitGroup.Done()
			processBody(global.Global.Push.AppId, global.Global.Push.AppSecret)
		}
		go process()
		waitGroup.Wait()
		global.Log.Info("access_token刷新完成")
	})
	if err != nil {
		global.Log.Fatal(err)
	}
}

func processBody(appId, appSecret string) {
	body{appId: appId,
		appSecret: appSecret,
		client:    resty.New().R()}.refresh()
}
