package GeTuiGo

import (
	"testing"
)

func getClient() (*Client, error) {
	return NewClient("8pBAMeizL7AToQifGbUqn1", "aj3YmXBs5l7Vj9x4UvFyiA", "kHUVG5uojo9rVJ4XrZ0yx2")
}

func TestNewClient(t *testing.T) {
	_, err := getClient()

	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_SinglePush(t *testing.T) {
	client, err := getClient()
	if err != nil {
		t.Fatal(err)
	}

	push := &Push{
		Message: NewMessage(TypeNotification),
		Notification: &TmplNotification{
			TransmissionType:    false,
			TransmissionContent: "",
			Style:               NewStyleSystem(),
		},
		Cid: "xxxxx",
	}

	result, err := client.SinglePush(push)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(result)
}

func TestClient_SaveListBody(t *testing.T) {
	client, err := getClient()
	if err != nil {
		t.Fatal(err)
	}

	push := &Push{
		Message: NewMessage(TypeNotification),
		Notification: &TmplNotification{
			TransmissionType:    false,
			TransmissionContent: "",
			Style:               NewStyleSystem(),
		},
	}

	result, taskId, desc, err := client.SaveListBody(push)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(result, taskId, desc)
}
