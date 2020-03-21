package GeTuiGo

import (
	"log"
	"testing"
)

func getClient(t *testing.T) *Client {
	client, err := NewClient("8pBAMeizL7AToQifGbUqn1", "aj3YmXBs5l7Vj9x4UvFyiA", "kHUVG5uojo9rVJ4XrZ0yx2")
	if err != nil {
		t.Fatal(err)
	}

	return client
}

func TestNewClient(t *testing.T) {
	getClient(t)
}

func TestClient_SinglePush(t *testing.T) {
	client := getClient(t)

	push := &Push{
		Message: NewMessage(TypeNotification),
		Notification: &TmplNotification{
			TransmissionType:    false,
			TransmissionContent: "",
			Style:               NewStyleSystem(),
		},
		Cid: "44b4da5e84150d87ea1509442d41e175",
	}

	result, err := client.SinglePush(push)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(result)
}

func TestClient_SinglePushBatch(t *testing.T) {
	client := getClient(t)
	pushList := make([]*Push, 0)
	message := NewMessage(TypeNotification)
	style := NewStyleSystem()
	style.Title = "测试title"
	style.Text = "测试test"

	push := &Push{
		Message: message,
		Notification: &TmplNotification{
			TransmissionType:    false,
			TransmissionContent: "",
			Style:               style,
		},
		Cid: "44b4da5e84150d87ea1509442d41e175",
	}

	push2 := &Push{
		Message: NewMessage(TypeNotification),
		Notification: &TmplNotification{
			TransmissionType:    false,
			TransmissionContent: "",
			Style:               NewStyleSystem(),
		},
		Cid: "xx",
	}
	pushList = append(pushList, push, push2)

	result, err := client.SinglePushBatch(pushList, true)
	t.Log(result, err)
}

func TestClient_SaveListBody(t *testing.T) {
	client := getClient(t)

	style := NewStyleSystem()
	style.Title = "测试title"
	style.Text = "测试test"

	push := &Push{
		Message: NewMessage(TypeNotification),
		Notification: &TmplNotification{
			TransmissionType:    false,
			TransmissionContent: "",
			Style:               style,
		},
	}

	result, taskId, desc, err := client.SaveListBody(push)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(result, taskId, desc)

	pushList := &PushList{
		Cid:        []string{"44b4da5e84150d87ea1509442d41e175"},
		TaskId:     taskId,
		Alias:      nil,
		NeedDetail: true,
	}

	result2, err2 := client.PushList(pushList)
	if err2 != nil {
		t.Fatal(err2)
	}

	t.Log(result2)
}

func TestClient_PushToApp(t *testing.T) {
	client := getClient(t)

	style := NewStyleSystem()
	style.Title = "PushToApp 测试title"
	style.Text = "PushToApp 测试test"

	push := &Push{
		Message: NewMessage(TypeNotification),
		Notification: &TmplNotification{
			TransmissionType:    false,
			TransmissionContent: "",
			Style:               style,
		},
	}

	result, taskId, desc, err := client.PushToApp(push)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(result, taskId, desc)
}

func TestClient_StopTask(t *testing.T) {
	client := getClient(t)
	result, respTaskId, err := client.StopTask("asd")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(result, respTaskId)
}

func TestClient_GetScheduleTask(t *testing.T) {
	client := getClient(t)

	result, err := client.GetScheduleTask("task123")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(result)
}

func TestClient_DelScheduleTask(t *testing.T) {
	client := getClient(t)

	result, err := client.DelScheduleTask("test123")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(result)
}

func TestClient_BindAlias(t *testing.T) {
	client := getClient(t)

	alias := make([]Alias, 1)
	alias[0] = Alias{
		Cid:   "44b4da5e84150d87ea1509442d41e175",
		Alias: "lee",
	}
	result, desc, err := client.BindAlias(alias)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(result, desc)

	result1, alias1, err := client.QueryAlias("44b4da5e84150d87ea1509442d41e175")
	if err != nil {
		t.Fatal(err)
	}

	log.Println(result1, alias1)
	if alias1 != "lee" {
		t.Fatal("别名设置失败")
	}
}

func TestClient_SetTags(t *testing.T) {
	client := getClient(t)

	result, err := client.SetTags("44b4da5e84150d87ea1509442d41e175", []string{"tag1"})
	if err != nil {
		t.Fatal(err)
	}

	t.Log(result)
	result1, tags, err := client.GetTags("44b4da5e84150d87ea1509442d41e175")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(result1, tags)
}

func TestClient_UserStatus(t *testing.T) {
	client := getClient(t)

	result, lastLogin, err := client.UserStatus("44b4da5e84150d87ea1509442d41e175")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(result, lastLogin)
}

func TestClient_GetPushResult(t *testing.T) {
	client := getClient(t)

	push := &Push{
		Message: NewMessage(TypeNotification),
		Notification: &TmplNotification{
			TransmissionType:    false,
			TransmissionContent: "",
			Style:               NewStyleSystem(),
		},
		Cid: "44b4da5e84150d87ea1509442d41e175",
	}

	result, err := client.SinglePush(push)
	if err != nil {
		t.Fatal(err)
	}

	result2, pushResultDetail, err := client.GetPushResult([]string{result.TaskId})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(result2, pushResultDetail)
}
