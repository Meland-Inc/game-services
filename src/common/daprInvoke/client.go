package daprInvoke

import (
	"context"

	daprC "github.com/dapr/go-sdk/client"
)

type DaprClient struct {
	serviceAppId string
	c            daprC.Client
}

func NewDaprClient(port string) (*DaprClient, error) {
	client, err := daprC.NewClientWithPort(port)
	if err != nil {
		return nil, err
	}

	return &DaprClient{c: client}, nil

}

func (this *DaprClient) SendMsg(serviceAppId, methodName string, data []byte) ([]byte, error) {
	content := &daprC.DataContent{
		ContentType: "application/json",
		Data:        data,
	}
	return this.c.InvokeMethodWithContent(context.Background(), serviceAppId, methodName, "post", content)
}

func (this *DaprClient) PubSubEventCall(topic string, data interface{}) error {
	return this.c.PublishEvent(context.Background(), "pubsub", topic, data)
}

func (this *DaprClient) Stop() {
	this.c.Close()
}
