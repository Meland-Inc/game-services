package daprInvoke

import (
	"context"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	daprC "github.com/dapr/go-sdk/client"
)

var client daprC.Client

func InitClient(port string) (err error) {
	client, err = daprC.NewClientWithPort(port)
	return err
}

func InvokeMethod(serviceAppId, methodName string, data []byte) ([]byte, error) {
	content := &daprC.DataContent{
		ContentType: "application/json",
		Data:        data,
	}
	return client.InvokeMethodWithContent(
		context.Background(),
		serviceAppId,
		methodName,
		"post",
		content,
	)
}

func PubSubEventCall(topic string, jsonString string) error {
	serviceLog.Debug("pubsubEventCall  [%v], [%v]", topic, jsonString)
	return client.PublishEvent(context.Background(), "pubsub", topic, jsonString)
}

func Stop() {
	client.Close()
}
