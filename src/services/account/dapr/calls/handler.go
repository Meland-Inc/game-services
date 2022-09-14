package daprCalls

import (
	"context"
	"fmt"

	"github.com/dapr/go-sdk/service/common"
)

func DemoServiceTestCallsHandler(ctx context.Context, in *common.InvocationEvent) (*common.Content, error) {
	fmt.Println("this is DemoServiceTestCallsHandler")
	out := &common.Content{
		Data:        []byte{},
		ContentType: in.ContentType,
		DataTypeURL: in.DataTypeURL,
	}
	return out, nil
}
