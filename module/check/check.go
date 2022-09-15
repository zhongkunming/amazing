package check

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"math/rand"
	"service-hub/config"
	"service-hub/global"
	"service-hub/model"
	"service-hub/util"
	"sync"
	"time"
)

const spec = "0 0 8 * * ?"

type DailyCheck struct {
}

func (r DailyCheck) Description() string {
	return fmt.Sprint("签到")
}

func (r DailyCheck) Load() {
	_, err := global.Cron.AddFunc(spec, func() {
		waitGroup := sync.WaitGroup{}
		global.Log.Infof("%s, 开始签到", time.Now().Format("2006-01-02"))
		waitGroup.Add(len(global.Global.Users))

		for _, user := range global.Global.Users {
			processBody := func(u config.User) {
				defer waitGroup.Done()
				r.process(u)
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

func (r DailyCheck) process(user config.User) {

	client := resty.New().R()
	rand.Seed(time.Now().UnixMicro())

	time.Sleep(time.Duration(rand.Intn(30)) * time.Second)

	err := r.login(client, &user)
	if err != nil {
		global.Log.Errorf("%s 登录异常，%s", user.Email, err)
		return
	}

	time.Sleep(time.Duration(rand.Intn(30)) * time.Second)
	result, err := r.checkin(client)
	if err != nil {
		global.Log.Errorf("%s 签到异常，%s", user.Email, err)
		return
	}

	if result.Ret == 0 {
		global.Log.Errorf("%s 签到异常: %s", user.Email, result.Msg)
	} else if result.Ret == 1 {
		global.Log.Infof("%s 签到成功: %s", user.Email, result.Msg)
		global.Log.Infof("%s 剩余未使用流量: %s", user.Email, result.TrafficInfo["unUsedTraffic"])
	}
}

func (r DailyCheck) login(client *resty.Request, user *config.User) error {
	userJsonBytes, _ := json.Marshal(user)
	body := make(map[string]string)
	_ = json.Unmarshal(userJsonBytes, &body)
	client.SetBody(body)
	loginResp, err := client.Post(global.Global.App.LoginUrl)
	if err != nil {
		return errors.New(fmt.Sprintf("请求登录接口失败, %s", err))
	}
	var loginResult = &model.LoginResult{}
	err = json.Unmarshal([]byte(util.TransByte(loginResp.Body())), loginResult)

	if 1 != loginResult.Ret {
		return errors.New(fmt.Sprintf("登录失败, %s", loginResult.Msg))
	}
	return err
}

func (r DailyCheck) checkin(client *resty.Request) (*model.CheckInResult, error) {
	checkInResp, err := client.Post(global.Global.App.CheckInUrl)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("请求签到接口失败, %s", err))
	}
	var result = &model.CheckInResult{}
	err = json.Unmarshal([]byte(util.TransByte(checkInResp.Body())), result)
	return result, err
}
