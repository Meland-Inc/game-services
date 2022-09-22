package daprCalls

import (
	"context"
	"encoding/json"
	"fmt"
	"game-message-core/grpc/methodData"
	"net/url"

	"github.com/Meland-Inc/game-services/src/common/daprInvoke"
	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/services/account/msgChannel"
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

	escStr, err := url.QueryUnescape(string(in.Data))
	if err != nil {
		return nil, err
	}

	serviceLog.Info("account received clientPbMsg data: %s", escStr)

	input := &methodData.PullClientMessageInput{}
	err = json.Unmarshal([]byte(escStr), input)
	if err != nil {
		serviceLog.Error("account Unmarshal to PullClientMessageInput fail err: %+v", err)
		return resFunc(false, fmt.Errorf("PullClientMessageInput unMarshal fail"))
	}

	msgChannel.GetInstance().CallClientMsg(input)
	return resFunc(true, nil)
}
