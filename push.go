package GeTuiGo

import (
	"bytes"
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
	conditions    []Condition        // 筛选目标用户条件
	speed         int                // 可选字段 推送速度控制
	pushTime      time.Time          // 定时下发时间
	taskName      string             // 可选字段 任务名称 可以给多个任务指定相同的task_name，后面用task_name查询推送结果能得到多个任务的结果
	durationBegin time.Time          // 可选字段 设定展示开始时间，格式为yyyy-MM-dd HH:mm:ss
	durationEnd   time.Time          // 可选字段 设定展示结束时间，格式为yyyy-MM-dd HH:mm:ss
}

type PushResult struct {
	Result string `json:"result"`
	TaskId string `json:"taskid"`
	Desc   string `json:"desc"`
	Status string `json:"status"`
}

type PushList struct {
	Cid        []string `json:"cid"`         // cid为cid list，与alias list二选一
	TaskId     string   `json:"taskid"`      // 必传: 任务号，取save_list_body返回的taskid
	Alias      []string `json:"alias"`       // alias为alias list，与cid list二选一
	NeedDetail bool     `json:"need_detail"` // 是否需要返回每个CID的状态
}

type PushListResult struct {
	Result       string            `json:"result"`        // 响应结果，见详情
	TaskId       string            `json:"taskid"`        // 任务标识号，用于tolist接口
	Desc         string            `json:"desc"`          // 错误信息描述
	CidDetails   map[string]string `json:"cid_details"`   // 目标cid用户推送结果详情
	AliasDetails map[string]string `json:"alias_details"` // 目标别名用户推送结果详情
}

// 筛选目标用户条件
type Condition struct {
	Key     string   `json:"key"`      // 必传: 筛选条件类型名称(省市region,手机类型phonetype,用户标签tag)
	Values  []string `json:"values"`   // 必传: 筛选参数
	OptType int      `json:"opt_type"` // 必传: 筛选参数的组合，0:取参数并集or，1：交集and，2：相当与not in {参数1，参数2，....}
}

// 为推送消息添加筛选条件
func (push *Push) AppendCondition(cond Condition) {
	push.conditions = append(push.conditions, cond)
}

// 推送速度控制
func (push *Push) SetSpeed(speed int) {
	push.speed = speed
}

// 设定展示开始时间
//  begin 	开始时间
//  end 	结束时间
func (push *Push) SetDuration(begin, end time.Time) {
	push.durationBegin = begin
	push.durationEnd = end
}

// 定时下发时间
func (push *Push) SetPushTime(pushTime time.Time) {
	push.pushTime = pushTime
}

// 转为json字符
func (push *Push) ToJsonString(appKey string) string {
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

	// 筛选条件
	if len(push.conditions) > 0 {
		data["condition"] = push.conditions
	}

	// 速度
	if push.speed > 0 {
		data["speed"] = push.speed
	}

	// 展示时间
	if !push.durationBegin.IsZero() {
		data["duration_begin"] = push.durationBegin.Format("2020-03-21 14:01:03")
	}

	if !push.durationEnd.IsZero() {
		data["duration_end"] = push.durationEnd.Format("2020-03-21 14:01:03")
	}

	// 定时下发时间
	if !push.pushTime.IsZero() {
		data["push_time"] = push.pushTime.Format("2020-03-21 14:01:03")
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

func authRequest(method, url, data string) (*http.Request, error) {

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

	req, err := http.NewRequest("POST", url, strings.NewReader(push.ToJsonString(c.appKey)))
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

// 批量单推接口
//  在给每个用户的推送内容都不同的情况下，又因为单推消息发送较慢，可以使用此接口。
func (c *Client) SinglePushBatch(pushList []*Push, needDetail bool) (result, taskId, desc string, err error) {
	url := fmt.Sprintf("https://restapi.getui.com/v1/%s/push_single_batch", c.appId)

	list := make([]string, len(pushList))
	for _, push := range pushList {
		str := push.ToJsonString(c.appKey)
		list = append(list, str)
	}
	body := fmt.Sprintf(`{"msg_list":[%s],"need_detail":%s}`, strings.Join(list, ","), strconv.FormatBool(needDetail))
	req, err := http.NewRequest("POST", url, strings.NewReader(body))
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

	var resultData map[string]string

	err = json.Unmarshal(respBody, &resultData)
	if err != nil {
		return
	}

	result = resultData["result"]
	taskId = resultData["taskid"]
	desc = resultData["desc"]
	return
}

// 在执行群推任务的时候，需首先执行save_list_body接口，将推送消息保存在服务器上，后面可以重复调用toList接口将保存的消息发送给不同的目标用户。
//  push 要推送的消息
//
//  返回：
//  result  推送结果，ok 鉴权成功，使用 ResultXXX 常量进行判断
//  taskId  任务编号
//  desc    错误信息描述
func (c *Client) SaveListBody(push *Push) (result, taskId, desc string, err error) {
	url := fmt.Sprintf("https://restapi.getui.com/v1/%s/save_list_body", c.appId)

	req, err := http.NewRequest("POST", url, strings.NewReader(push.ToJsonString(c.appKey)))
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

	var respData struct {
		Result string `json:"result"` // 响应结果，见详情
		TaskId string `json:"taskid"` // 任务标识号，用于tolist接口的taskid
		Desc   string `json:"desc"`   // 错误信息描述
	}
	err = json.Unmarshal(respBody, &respData)
	if err != nil {
		return
	}

	return respData.Result, respData.TaskId, respData.Desc, nil
}

// 消息群发给cid list或者alias list列表对应的客户群，当两者并存的时候，以cid为准；并使用save_list_body返回的taskId，调用toList接口，完成群推推送。
//  pushList	与alias list二选一
//
//  result		推送结果
func (c *Client) PushList(pushList *PushList) (result PushListResult, err error) {
	url := fmt.Sprintf("https://restapi.getui.com/v1/%s/push_list", c.appId)

	body, _ := json.Marshal(pushList)

	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
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

	err = json.Unmarshal(respBody, &result)
	if err != nil {
		return
	}

	return
}

// 群推
//  针对某个，根据筛选条件，将消息群发给符合条件客户群
func (c *Client) PushToApp(push *Push) (result, taskId, desc string) {
	url := fmt.Sprintf("https://restapi.getui.com/v1/%s/push_app", c.appId)
	req, err := http.NewRequest("POST", url, strings.NewReader(push.ToJsonString(c.appKey)))
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

	var resultData map[string]string

	err = json.Unmarshal(respBody, &resultData)
	if err != nil {
		return
	}

	result = resultData["result"]
	taskId = resultData["taskid"]
	desc = resultData["desc"]
	return
}

// stop群推任务
//  在有效期内的消息进行停止
func (c *Client) StopTask(taskId string) (result, respTaskId string) {
	url := fmt.Sprintf("https://restapi.getui.com/v1/%s/stop_task/%s", c.appId, taskId)
	req, err := http.NewRequest("DELETE", url, nil)
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

	var resultData map[string]string

	err = json.Unmarshal(respBody, &resultData)

	result = resultData["result"]
	taskId = resultData["taskid"]
	return
}

type ScheduleTaskResult struct {
	result     string
	taskDetail struct {
		pushContent string // 推送类容（transmission的内容）
		pushTime    string // 推送时间
		creatTime   string // 任务创建时间
		sendResult  string // 任务状态
	}
	taskId string // 任务Id
}

// 定时任务查询接口
//  应用场景: 该接口主要用来在需要查看返回已提交的定时任务的相关信息。
func (c *Client) GetScheduleTask(taskId string) (*ScheduleTaskResult, error) {
	url := fmt.Sprintf("https://restapi.getui.com/v1/%s/get_schedule_task", c.appId)
	req, err := http.NewRequest("POST", url, strings.NewReader(fmt.Sprintf(`{"taskid":"%s"}`, taskId)))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("authtoken", c.authToken)
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	respBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var resultData = &ScheduleTaskResult{}
	err = json.Unmarshal(respBody, &resultData)
	return resultData, nil
}

// 定时任务删除接口
//  应用场景: 用来删除还未下发的任务
func (c *Client) DelScheduleTask(taskId string) (result string, err error) {
	url := fmt.Sprintf("https://restapi.getui.com/v1/%s/del_schedule_task", c.appId)
	req, err := http.NewRequest("POST", url, strings.NewReader(fmt.Sprintf(`{"taskid":"%s"}`, taskId)))
	if err != nil {
		return "", err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("authtoken", c.authToken)
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	respBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	var resultData = map[string]string{}
	err = json.Unmarshal(respBody, &resultData)
	return resultData["result"], nil
}

type Alias struct {
	Cid   string `json:"cid"`
	Alias string `json:"alias"`
}

// 绑定别名
//  一个ClientID只能绑定一个别名，若已绑定过别名的ClientID再次绑定新别名，
//  则认为与前一个别名自动解绑，绑定新别名
//  允许将多个ClientID和一个别名绑定，如用户使用多终端，则可将多终端对应的ClientID绑定为一个别名，
//  目前一个别名最多支持绑定10个ClientID
func (c *Client) BindAlias(aliasList []Alias) (result, desc string, err error) {
	url := fmt.Sprintf("https://restapi.getui.com/v1/%s/bind_alias", c.appId)

	data, _ := json.Marshal(aliasList)
	body := fmt.Sprintf(`{"alias_list":%s}`, data)
	req, err := http.NewRequest("POST", url, strings.NewReader(body))
	if err != nil {
		return "", "", err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("authtoken", c.authToken)
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", "", err
	}

	respBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", "", err
	}

	var resultData = map[string]string{}
	err = json.Unmarshal(respBody, &resultData)
	return resultData["result"], resultData["desc"], nil
}

// 绑定别名
//  一个ClientID只能绑定一个别名，若已绑定过别名的ClientID再次绑定新别名，
//  则认为与前一个别名自动解绑，绑定新别名
//  允许将多个ClientID和一个别名绑定，如用户使用多终端，则可将多终端对应的ClientID绑定为一个别名，
//  目前一个别名最多支持绑定10个ClientID
func (c *Client) BindAlia(alias, cid string) (result, desc string, err error) {
	data := make([]Alias, 1)
	data = append(data, Alias{
		Cid:   cid,
		Alias: alias,
	})

	return c.BindAlias(data)
}

// 单个cid和别名解绑
func (c *Client) UnBindAlias(cid, alias string) (result string, err error) {
	url := fmt.Sprintf("https://restapi.getui.com/v1/%s/unbind_alias", c.appId)

	data, err := json.Marshal(fmt.Sprintf(`{"cid":"%s","alias":"%s"}`, cid, alias))
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", url, strings.NewReader(string(data)))
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

	var resultData = map[string]string{}
	err = json.Unmarshal(respBody, &resultData)
	return resultData["result"], nil
}

// 解绑别名所有cid
func (c *Client) UnBindAliasAll(alias string) (result, desc string, err error) {
	url := fmt.Sprintf("https://restapi.getui.com/v1/%s/unbind_alias_all", c.appId)

	data, err := json.Marshal(fmt.Sprintf(`{"alias":"%s"}`, alias))
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", url, strings.NewReader(string(data)))
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

	var resultData = map[string]string{}
	err = json.Unmarshal(respBody, &resultData)
	return resultData["result"], resultData["desc"], nil
}

// 查询别名cid
//  通过传入的别名查询对应的cid信息
func (c *Client) QueryCid(alias string) (result string, cidList []string, err error) {
	url := fmt.Sprintf("https://restapi.getui.com/v1/%s/query_cid/%s", c.appId, alias)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("authtoken", c.authToken)
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", nil, err
	}

	respBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", nil, err
	}

	var resultData struct {
		Result string   `json:"result"`
		Cid    []string `json:"cid"`
	}
	err = json.Unmarshal(respBody, &resultData)
	return resultData.Result, resultData.Cid, nil
}

// 查询cid别名
//  通过传入的cid查询对应的别名
func (c *Client) QueryAlias(cid string) (result string, alias string, err error) {
	url := fmt.Sprintf("https://restapi.getui.com/v1/%s/query_alias/%s", c.appId, cid)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", "", err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("authtoken", c.authToken)
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", "", err
	}

	respBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", "", err
	}

	var resultData map[string]string
	err = json.Unmarshal(respBody, &resultData)
	return resultData["result"], resultData["alias"], nil
}

// 对指定用户设置tag属性
func (c *Client) SetTags(cid string, tagList []string) (result string, err error) {
	data := struct {
		Cid     string   `json:"cid"`
		TagList []string `json:"tag_list"`
	}{
		Cid:     cid,
		TagList: tagList,
	}

	body, err := json.Marshal(data)
	url := fmt.Sprintf("https://restapi.getui.com/v1/%s/set_tags", c.appKey)
	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
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

	var resultData map[string]string
	err = json.Unmarshal(respBody, &resultData)
	return resultData["result"], nil
}

// 查询指定用户tag属性
func (c *Client) GetTags(cid string) (result, tags string, err error) {
	url := fmt.Sprintf("https://restapi.getui.com/v1/%s/get_tags/%s", c.appKey, cid)
	req, err := http.NewRequest("GET", url, nil)
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

	var resultData map[string]string
	err = json.Unmarshal(respBody, &resultData)
	return resultData["result"], resultData["cid"], nil
}

// 添加黑名单用户
func (c *Client) AddBlackList(cidList []string) (result, desc string, err error) {
	data := fmt.Sprintf(`{"cid":["%s"]}`, strings.Join(cidList, `","`))
	url := fmt.Sprintf("https://restapi.getui.com/v1/%s/user_blk_list", c.appKey)
	req, err := http.NewRequest("POST", url, strings.NewReader(data))
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

	var resultData map[string]string
	err = json.Unmarshal(respBody, &resultData)
	return resultData["result"], resultData["desc"], nil
}

// 移除黑名单用户
func (c *Client) RemoveBlackList(cidList []string) (result, desc string, err error) {
	data := fmt.Sprintf(`{"cid":["%s"]}`, strings.Join(cidList, `","`))
	url := fmt.Sprintf("https://restapi.getui.com/v1/%s/user_blk_list", c.appKey)
	req, err := http.NewRequest("DELETE", url, strings.NewReader(data))
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

	var resultData map[string]string
	err = json.Unmarshal(respBody, &resultData)
	return resultData["result"], resultData["desc"], nil
}

// 查询用户状态
//  调用此接口可获取用户状态，如在线不在线
func (c *Client) UserStatus(cid string) (result, lastLogin string, err error) {
	url := fmt.Sprintf("https://restapi.getui.com/v1/%s/user_status/%s", c.appKey, cid)
	req, err := http.NewRequest("GET", url, nil)
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

	var resultData map[string]string
	err = json.Unmarshal(respBody, &resultData)
	return resultData["result"], resultData["lastlogin"], nil
}

// 查询数据对象
type PushResultDetail struct {
	TaskId     string `json:"taskid"`      // 任务号
	MsgTotal   int    `json:"msg_total"`   // 有效可下发总数
	MsgProcess int    `json:"msg_process"` // 消息回执总数
	ClickNum   int    `json:"click_num"`   // 用户点击数
	PushNum    int    `json:"push_num"`    // im下发总量
	// iOS推送结果数据，详细字段参考GT
	Apn struct {
		Displayed int    `json:"displayed"` //展示数
		Result    string `json:"result"`
		Feedback  int    `json:"feedback"` // 回执数
		Clicked   int    `json:"clicked"`  //点击数
		Sent      int    `json:"sent"`     // 成功下发数
	}
	// 个推推送结果数据
	GT struct {
		Sent      int `json:"sent"`      // 成功下发数
		Feedback  int `json:"feedback"`  // 回执数
		Clicked   int `json:"clicked"`   //点击数
		Displayed int `json:"displayed"` //展示数
	}
}

// 获取推送结果接口
//  调用此接口查询推送数据，可查询消息有效可下发总数，消息回执总数和用户点击数等结果。
func (c *Client) GetPushResult(taskIdList []string) (result string, pushResultList []PushResultDetail, err error) {
	if len(taskIdList) == 0 {
		return
	}
	data := fmt.Sprintf(`{"taskIdList":["%s"]}`, strings.Join(taskIdList, `","`))
	url := fmt.Sprintf("https://restapi.getui.com/v1/%s/push_result", c.appKey)
	req, err := http.NewRequest("POST", url, strings.NewReader(data))
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

	var resultData struct {
		Result string             `json:"result"`
		Data   []PushResultDetail `json:"data"`
	}
	err = json.Unmarshal(respBody, &resultData)
	return resultData.Result, resultData.Data, nil
}

type PushResultByGroup struct {
	Result     string `json:"result"`      // 操作结果 成功返回ok 见详情
	MsgTotal   int    `json:"msg_total"`   // 百日内活跃用户数
	OnlineNum  int    `json:"online_num"`  // 消息实际下发数
	MsgProcess int    `json:"msg_process"` // 消息接收数
	ShowNum    int    `json:"show_num"`    // 消息展示数
	ClickNum   int    `json:"click_num"`   // 消息点击数
	Desc       string `json:"desc"`        // 错误详情
}

// 根据任务组名获取推送结果数据
//  根据任务组名查询推送结果，返回结果包括百日内联网用户数（活跃用户数）、实际下发数、到达数、展示数、点击数。
func (c *Client) GetPushResultByGroup(groupName string) (result PushResultByGroup, err error) {
	url := fmt.Sprintf("https://restapi.getui.com/v1/%s/get_push_result_by_group_name/%s", c.appKey, groupName)
	req, err := http.NewRequest("POST", url, nil)
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

	var resultData PushResultByGroup
	err = json.Unmarshal(respBody, &resultData)
	return resultData, nil
}

type AppUserStat struct {
	AppId              string `json:"app_id"`
	Date               string `json:"date"`               // 查询的日期（格式：yyyy-MM-dd）
	NewRegisterCount   int    `json:"new_regist_count"`   // 新注册用户数
	RegisterTotalCount int    `json:"regist_total_count"` // 累计注册用户数
	ActiveCount        int    `json:"active_count"`       // 活跃用户数
	OnlineCount        int    `json:"online_count"`       // 在线用户数
}

// 获取单日用户数据接口
//  调用此接口查询推送数据，可查询消息有效可下发总数，消息回执总数和用户点击数等结果。
func (c *Client) QueryAppUser(date time.Time) (result string, stat AppUserStat, err error) {
	url := fmt.Sprintf("https://restapi.getui.com/v1/%s/query_app_push/%s", c.appKey, date.Format("20200321"))
	req, err := http.NewRequest("POST", url, nil)
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

	var resultData struct {
		Result string      `json:"result"`
		Data   AppUserStat `json:"data"`
	}
	err = json.Unmarshal(respBody, &resultData)
	return resultData.Result, resultData.Data, nil
}

// 应用角标设置接口(仅iOS)
//  badge	应用icon上显示的数字
//  msgId	请求的msgid
func (c *Client) IosSetBadge(badge int, msgId string, cidList, deviceTokenList []string) (result, desc string, err error) {
	data := struct {
		MsgId           string   `json:"msgid"`
		Badge           int      `json:"badge"`
		CidList         []string `json:"cid_list"`
		DeviceTokenList []string `json:"devicetoken_list"`
	}{
		MsgId:           msgId,
		Badge:           badge,
		CidList:         cidList,
		DeviceTokenList: deviceTokenList,
	}
	body, err := json.Marshal(data)
	if err != nil {
		return
	}
	url := fmt.Sprintf("https://restapi.getui.com/v1/%s/set_badge", c.appKey)
	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
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

	var resultData map[string]string
	err = json.Unmarshal(respBody, &resultData)
	return resultData["result"], resultData["desc"], nil
}

// 按条件查询用户数
//  通过指定查询条件来查询满足条件的用户数量
func (c *Client) QueryUserCount(condition Condition) (result string, userCount int, err error) {
	data := struct {
		Condition Condition `json:"condition"`
	}{Condition: condition}

	body, err := json.Marshal(data)
	if err != nil {
		return
	}
	url := fmt.Sprintf("https://restapi.getui.com/v1/%s/query_user_count", c.appKey)
	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
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

	var resultData map[string]interface{}
	err = json.Unmarshal(respBody, &resultData)
	return resultData["result"].(string), resultData["desc"].(int), nil
}

// 获取可用bi标签
//  查询应用可用的bi标签列表
func (c *Client) QueryBiTags() (result string, tags []string, err error) {
	url := fmt.Sprintf("https://restapi.getui.com/v1/%s/query_bi_tags", c.appKey)
	req, err := http.NewRequest("POST", url, nil)
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

	var resultData map[string]interface{}
	err = json.Unmarshal(respBody, &resultData)
	return resultData["result"].(string), resultData["tags"].([]string), nil
}
