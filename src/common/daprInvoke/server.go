package daprInvoke

import (
	"context"
	"fmt"
	"net/http"

	"github.com/dapr/go-sdk/service/common"
	daprd "github.com/dapr/go-sdk/service/grpc"
)

type DaprServer struct {
	service common.Service
}

func NewDaprServer(port string) (*DaprServer, error) {
	s, err := daprd.NewService(fmt.Sprintf(":%s", port))
	if err != nil {
		return nil, err
	}
	return &DaprServer{service: s}, nil
}

func (this *DaprServer) RegistMothod(
	name string,
	f func(context.Context, *common.InvocationEvent) (*common.Content, error),
) error {
	return this.service.AddServiceInvocationHandler(name, f)
}

// RegistEventHandle 每一个topic只能一个地方订阅
// 如果需要实现多次订阅.需要自己实现消息派发模块统一派发回调
func (this *DaprServer) RegistEventHandle(topic string, fn func(ctx context.Context, e *common.TopicEvent) (retry bool, err error)) error {
	sub := &common.Subscription{
		PubsubName: "pubsub",
		Topic:      topic,
		Route:      fmt.Sprintf("/%s", topic),
		Metadata:   map[string]string{},
	}
	return this.service.AddTopicEventHandler(sub, fn)
}

func (this *DaprServer) Run() error {
	var err error
	go func() {
		defer func() {
			if err := recover(); err != nil {
				err = fmt.Errorf("save mapResource panic err:%+v", err)
			}
		}()

		err := this.service.Start()
		if err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()
	return err
}
