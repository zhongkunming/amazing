package daily_check

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
	"time"
)

type body struct {
	user   config.DailyCheckUsersConfig
	client *resty.Request
}

func (r body) do() {
	rand.Seed(time.Now().UnixMicro())

	time.Sleep(time.Duration(rand.Intn(30)) * time.Second)
	err := r.login()
	if err != nil {
		global.Log.Errorf("%s 登录异常，%s", r.user.Email, err)
		return
	}

	time.Sleep(time.Duration(rand.Intn(30)) * time.Second)
	result, err := r.checkin()
	if err != nil {
		global.Log.Errorf("%s 签到异常，%s", r.user.Email, err)
		return
	}

	if result.Ret == 0 {
		global.Log.Errorf("%s 签到异常: %s", r.user.Email, result.Msg)
	} else if result.Ret == 1 {
		global.Log.Infof("%s 签到成功: %s", r.user.Email, result.Msg)
		global.Log.Infof("%s 剩余未使用流量: %s", r.user.Email, result.TrafficInfo["unUsedTraffic"])
	}
}

func (r body) login() error {
	userJsonBytes, _ := json.Marshal(r.user)
	body := make(map[string]string)
	_ = json.Unmarshal(userJsonBytes, &body)
	r.client.SetBody(body)
	loginResp, err := r.client.Post(global.Global.DailyCheckConfig.LoginUrl)
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

func (r body) checkin() (*model.CheckInResult, error) {
	checkInResp, err := r.client.Post(global.Global.DailyCheckConfig.CheckInUrl)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("请求签到接口失败, %s", err))
	}
	var result = &model.CheckInResult{}
	err = json.Unmarshal([]byte(util.TransByte(checkInResp.Body())), result)
	return result, err
}
