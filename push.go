package GeTuiGo

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Client struct {
	appId               string
	appKey              string
	masterSecret        string
	authToken           string
	authTokenExpireTime int
}

// 推送消息体
type Push struct {
	Message       *Message           // 消息内容
	Notification  *TmplNotification  // 点开通知打开应用模板
	Link          *TmplLink          // 点开通知打开网页模板
	NotifyPopLoad *TmplNotifyPopLoad // 点击通知弹窗下载模板
	StartActivity *TmplStartActivity // 点开通知打开应用内特定页面模板
	Transmission  *TmplTransmission  // 透传消息模板
	PushInfo      *ApnPushInfo       // apns推送消息, json串，当手机为ios，并且为离线的时候
	Cid           string             // 与alias二选一
	Alias         string             // 与cid二选一
	RequestId     string             // 必传: 请求唯一标识
}

type PushResult struct {
	Result string `json:"result"`
	TaskId string `json:"taskid"`
	Desc   string `json:"desc"`
	Status string `json:"status"`
}

func (push *Push) Content(appKey string) string {
	// 构造要发送的数据
	push.Message.AppKey = appKey

	data := map[string]interface{}{
		"message": push.Message,
	}

	switch push.Message.MsgType {
	case "notification":
		data["notification"] = push.Notification
		break
	case "link":
		data["link"] = push.Link
		break
	case "notypopload":
		data["notypopload"] = push.NotifyPopLoad
		break
	case "startactivity":
		data["startactivity"] = push.StartActivity
		break
	case "transmission":
		data["transmission"] = push.Transmission
		break
	}

	if push.PushInfo != nil {
		data["push_info"] = push.PushInfo
	}

	if push.Cid != "" {
		data["cid"] = push.Cid
	} else if push.Alias != "" {
		data["alias"] = push.Alias
	}

	// 请求唯一标识为空时，创建一个
	if push.RequestId == "" {
		push.RequestId = fmt.Sprintf("%d", time.Now().UnixNano())
		data["requestid"] = push.RequestId
	}

	res, _ := json.Marshal(data)
	return string(res)
}

type ListBody struct {
	Message       Message           // 消息内容
	Notification  TmplNotification  // 点开通知打开应用模板
	Link          TmplLink          // 点开通知打开网页模板
	NotifyPopLoad TmplNotifyPopLoad // 点击通知弹窗下载模板
	StartActivity TmplStartActivity // 点开通知打开应用内特定页面模板
	Transmission  TmplTransmission  // 透传消息模板
	PushInfo      ApnPushInfo       // 	apns推送消息, json串，当手机为ios，并且为离线的时候
	TaskName      string            // 任务名称,可以给多个任务指定相同的task_name，后面用task_name查询推送结果能得到多个任务的结果
}

// 用户身份验证通过获得auth_token权限令牌，后面的请求都需要带上auth_token
func (c *Client) getAutoToken(appId, appKey, masterSecret string) (authToken string, expTime int, err error) {
	timestamp := time.Now().UnixNano() / 1000000
	sign := sha256.Sum256([]byte(fmt.Sprintf("%s%d%s", appKey, timestamp, masterSecret)))
	data := fmt.Sprintf(`{"sign":"%x","timestamp":"%d","appkey":"%s"}`, sign, timestamp, appKey)
	url := fmt.Sprintf("https://restapi.getui.com/v1/%s/auth_sign", appId)
	req, err := http.NewRequest("POST", url, strings.NewReader(data))
	if err != nil {
		return
	}

	req.Header.Add("Content-Type", "application/json")
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}

	respBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}

	var respData struct {
		Result     string `json:"result"`
		ExpireTime string `json:"expire_time"`
		AuthToken  string `json:"auth_token"`
	}

	err = json.Unmarshal(respBody, &respData)
	if err != nil {
		return
	}

	if respData.Result != ResultOk {
		err = errors.New(respData.Result)
	}

	expTime, _ = strconv.Atoi(respData.ExpireTime)
	authToken = respData.AuthToken
	return
}

func NewClient(appId, appKey, masterSecret string) (*Client, error) {
	client := &Client{
		appId:        appId,
		appKey:       appKey,
		masterSecret: masterSecret,
	}

	token, expTime, err := client.getAutoToken(appId, appKey, masterSecret)
	if err != nil {
		return nil, err
	}

	client.authToken = token
	client.authTokenExpireTime = expTime
	return client, nil
}

// 对使用App的某个用户，单独推送消息
//  push 要推送的消息
//
//  返回：
//  result  推送结果，ok 鉴权成功，使用 ResultXXX 常量进行判断
//  taskId  任务编号
//  desc    错误信息描述
//  status  推送结果:
//  - successed_offline 离线下发
//  - successed_online  在线下发
//  - successed_ignore  非活跃用户不下发
func (c *Client) SinglePush(push *Push) (result PushResult, err error) {
	url := fmt.Sprintf("https://restapi.getui.com/v1/%s/push_single", c.appId)

	req, err := http.NewRequest("POST", url, strings.NewReader(push.Content(c.appKey)))
	if err != nil {
		return
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("authtoken", c.authToken)
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}

	respBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}

	var respData PushResult
	err = json.Unmarshal(respBody, &respData)
	if err != nil {
		return
	}

	return respData, nil
}
