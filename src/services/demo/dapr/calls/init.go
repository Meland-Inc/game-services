package daprCalls

import (
	"context"
	"net/url"

	"github.com/Meland-Inc/game-services/src/common/daprInvoke"
	"github.com/Meland-Inc/game-services/src/common/serviceLog"
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
	escStr, err := url.QueryUnescape(string(in.Data))
	if err != nil {
		return nil, err
	}

	serviceLog.Info("this is DemoServiceTestCallsHandler, %s", escStr)
	out := &common.Content{
		Data:        []byte{},
		ContentType: in.ContentType,
		DataTypeURL: in.DataTypeURL,
	}
	return out, nil
}
