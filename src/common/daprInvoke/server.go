package daprInvoke

import (
	"context"
	"fmt"

	"github.com/dapr/go-sdk/service/common"
	"github.com/dapr/go-sdk/service/grpc"
)

var server common.Service

func InitServer(port string) (err error) {
	server, err = grpc.NewService(fmt.Sprintf(":%s", port))
	return err
}

func AddServiceInvocationHandler(
	name string,
	f func(context.Context, *common.InvocationEvent) (*common.Content, error),
) error {
	return server.AddServiceInvocationHandler(name, f)
}

func AddTopicEventHandler(topic string,
	fn func(ctx context.Context, e *common.TopicEvent) (retry bool, err error),
) error {
	sub := &common.Subscription{
		PubsubName: "pubsub",
		Topic:      topic,
		Route:      fmt.Sprintf("/%s", topic),
		Metadata:   map[string]string{},
	}
	return server.AddTopicEventHandler(sub, fn)
}

func Start() error {
	return server.Start()
}

