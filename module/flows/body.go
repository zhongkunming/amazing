package flows

import (
	"amazing/config"
	"amazing/global"
	"amazing/model"
	"amazing/util"
	"encoding/json"
	"errors"
	"fmt"
	huge "github.com/dablelv/go-huge-util"
	"github.com/go-resty/resty/v2"
	"math/rand"
	"time"
)

type body struct {
	user     config.FLowsUser
	loginUrl string
	flowsUrl string
	client   *resty.Request
}

func (r body) do() {

	r.client.SetHeader("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36")
	r.client.SetHeader("Content-Type", "application/json")

	rand.Seed(time.Now().UnixMicro())

	time.Sleep(time.Duration(rand.Intn(30)) * time.Second)
	err := r.login()
	if err != nil {
		global.Log.Errorf("%s 登录异常，%s", r.user.Email, err)
		return
	}

	global.Log.Infof("%s 登录成功", r.user.Email)

	time.Sleep(time.Duration(rand.Intn(30)) * time.Second)
	r.checkin()
}

func (r body) login() error {
	userJsonBytes, _ := json.Marshal(r.user)
	body := make(map[string]string)
	_ = json.Unmarshal(userJsonBytes, &body)
	r.client.SetBody(body)
	loginResp, err := r.client.Post(r.loginUrl)
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

func (r body) checkin() {
	r.client.SetHeader("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36")
	r.client.SetHeader("Content-Type", "application/json")

	checkInResp, err := r.client.Post(r.flowsUrl)
	if err != nil {
		global.Log.Errorf("请求签到接口失败, %s", err)
		return
	}

	var result = &model.CheckInResult{}
	err = json.Unmarshal([]byte(util.TransByte(checkInResp.Body())), result)

	if err != nil {
		global.Log.Errorf("%s 请求结果转换异常, %s", r.user.Email, err)
		return
	}

	indentJSON, _ := huge.ToIndentJSON(result)
	global.Log.Infof("%s 请求结果: %s", r.user.Email, indentJSON)

	if result.Ret == 1 {
		global.Log.Infof("%s 签到成功: %s", r.user.Email, result.Msg)
		global.Log.Infof("%s 剩余未使用流量: %s", r.user.Email, result.TrafficInfo["unUsedTraffic"])
	}

}
