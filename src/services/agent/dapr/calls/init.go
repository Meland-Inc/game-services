package daprCalls

import (
	"context"
	"fmt"

	"github.com/Meland-Inc/game-services/src/common/daprInvoke"
	"github.com/dapr/go-sdk/service/common"
)

func InitDaprCallHandle() (err error) {
	daprInvoke.AddServiceInvocationHandler("DemoServiceTestCallsHandler", DemoServiceTestCallsHandler)
	if err != nil {
		return err
	}

	return nil
}

func DemoServiceTestCallsHandler(ctx context.Context, in *common.InvocationEvent) (*common.Content, error) {
	fmt.Println("this is DemoServiceTestCallsHandler")
	out := &common.Content{
		Data:        []byte{},
		ContentType: in.ContentType,
		DataTypeURL: in.DataTypeURL,
	}
	return out, nil
}
