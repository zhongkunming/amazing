package push

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"service-hub/global"
)

type body struct {
	appId, appSecret string
	client           *resty.Request
}

const getAccessTokenURL = "https://api.weixin.qq.com/cgi-bin/token"

const pushMessageURL = "https://api.weixin.qq.com/cgi-bin/message/template/send"

func (r body) refresh() {
	refreshUrl := fmt.Sprintf("%s?grant_type=client_credential&appid=%s&secret=%s", getAccessTokenURL, r.appId, r.appSecret)
	response, _ := r.client.Get(refreshUrl)
	responseData := make(map[string]string)
	json.Unmarshal(response.Body(), &responseData)
	global.Global.Push.AccessToken = responseData["access_token"]
}

func MessagePushFlow(account, getFlow, remainFlow string) {
	pushMessageUrl := fmt.Sprintf("%s?access_token=%s", pushMessageURL, global.Global.Push.AccessToken)
	client := resty.New().R()
	messageBody := make(map[string]interface{})

	messageBody["touser"] = "o_NxN5p04g6LRiI94g6mih90x3OE"
	messageBody["template_id"] = "wI5Dez-w0st1HzyT_fMZc1iPmi5e-rCjQz8ZteROGmY"
	messageBody["url"] = "https://www.baidu.com/"

	messageBodyData := make(map[string]interface{})

	messageBodyDataName := make(map[string]interface{})
	messageBodyDataName["value"] = "钟坤明"
	messageBodyDataAccount := make(map[string]interface{})
	messageBodyDataAccount["value"] = account
	messageBodyDataGetFlow := make(map[string]interface{})
	messageBodyDataGetFlow["value"] = getFlow
	messageBodyDataRemainFlow := make(map[string]interface{})
	messageBodyDataRemainFlow["value"] = remainFlow

	messageBodyData["name"] = messageBodyDataName
	messageBodyData["account"] = messageBodyDataAccount
	messageBodyData["getFlow"] = messageBodyDataGetFlow
	messageBodyData["remainFlow"] = messageBodyDataRemainFlow
	messageBody["data"] = messageBodyData

	marshal, _ := json.Marshal(messageBody)
	client.SetBody(marshal)
	client.Post(pushMessageUrl)
}
