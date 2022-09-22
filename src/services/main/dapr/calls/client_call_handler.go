package daprCalls

import (
	"context"
	"encoding/json"
	"fmt"
	"game-message-core/grpc/methodData"
	"net/url"

	"github.com/Meland-Inc/game-services/src/common/daprInvoke"
	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/services/main/msgChannel"
	"github.com/dapr/go-sdk/service/common"
)

func ClientMessageHandler(ctx context.Context, in *common.InvocationEvent) (*common.Content, error) {
	resFunc := func(success bool, err error) (*common.Content, error) {
		out := &methodData.PullClientMessageOutput{}
		out.Success = success
		if err != nil {
			out.ErrMsg = err.Error()
		}
		content, _ := daprInvoke.MakeOutputContent(in, out)
		return content, err
	}

	serviceLog.Info("main service received clientPbMsg data: %v", string(in.Data))

	escStr, err := url.QueryUnescape(string(in.Data))
	if err != nil {
		return nil, err
	}

	input := &methodData.PullClientMessageInput{}
	err = json.Unmarshal([]byte(escStr), input)
	if err != nil {
		serviceLog.Error("main service Unmarshal to PullClientMessageInput fail err: %+v", err)
		return resFunc(false, fmt.Errorf("PullClientMessageInput unMarshal fail"))
	}

	msgChannel.GetInstance().CallClientMsg(input)
	return resFunc(true, nil)
}
