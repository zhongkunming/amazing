package flows

import (
	"amazing/config"
	"amazing/global"
	"github.com/go-resty/resty/v2"
	"sync"
	"time"
)

const spec = "0 0 8 * * ?"

//const spec = "0 */1 * * * ?"

type LoadableFlows struct{}

func (r LoadableFlows) CanLoad() bool {
	return true
}

func (r LoadableFlows) Load() {
	_, err := global.Cron.AddFunc(spec, func() {
		waitGroup := sync.WaitGroup{}
		global.Log.Infof("%s, 开始签到", time.Now().Format("2006-01-02"))

		users := global.Global.Flows.Users
		loginUrl := global.Global.Flows.LoginUrl
		flowsUrl := global.Global.Flows.FlowsUrl

		waitGroup.Add(len(users))
		for _, user := range users {
			processBody := func(u config.FLowsUser) {
				defer waitGroup.Done()
				body{user: u,
					loginUrl: loginUrl,
					flowsUrl: flowsUrl,
					client:   resty.New().R()}.do()
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
