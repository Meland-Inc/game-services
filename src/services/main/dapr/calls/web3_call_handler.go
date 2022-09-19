package daprCalls

import (
	"context"
	"fmt"

	"github.com/Meland-Inc/game-services/src/common/daprInvoke"
	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/global/component"
	message "github.com/Meland-Inc/game-services/src/global/web3Message"
	"github.com/Meland-Inc/game-services/src/services/main/msgChannel"
	"github.com/Meland-Inc/game-services/src/services/main/playerModel"
	"github.com/dapr/go-sdk/service/common"
	"github.com/spf13/cast"
)

func Web3DeductUserExpHandler(ctx context.Context, in *common.InvocationEvent) (*common.Content, error) {
	resFunc := func(success bool, err error) (*common.Content, error) {
		out := &message.DeductUserExpOutput{}
		out.DeductSuccess = success
		if err != nil {
			out.FailedReason = err.Error()
			serviceLog.Error("web3 deduct user exp err: %v", err)
		}
		content, _ := daprInvoke.MakeOutputContent(in, out)
		return content, err
	}

	serviceLog.Info("web3 deduct user exp received data: %v", string(in.Data))

	input := &message.DeductUserExpInput{}
	err := input.UnmarshalJSON(in.Data)
	if err != nil {
		return resFunc(false, fmt.Errorf("not math to dapr msg DeductUserExpInput"))
	}

	deductExp := cast.ToInt64(input.DeductExp)
	userId, err := cast.ToInt64E(input.UserId)
	if err != nil || userId < 1 || deductExp < 1 {
		return resFunc(false,
			fmt.Errorf("web3 deduct user exp invalid userId [%s] or expr[%d]", input.UserId, input.DeductExp),
		)
	}

	iPlayerModel, exist := component.GetInstance().GetModel(component.MODEL_NAME_PLAYER_DATA)
	if !exist {
		return resFunc(false, fmt.Errorf("player data model not found"))
	}
	dataModel, _ := iPlayerModel.(*playerModel.PlayerDataModel)
	sceneData, err := dataModel.GetPlayerSceneData(userId)
	if sceneData.Exp < int32(deductExp) {
		out, err := resFunc(false, fmt.Errorf("insufficient experience"))
		return out, err
	}

	msgChannel.GetInstance().CallServiceMsg(&msgChannel.ServiceMsgData{
		MsgId:   string(message.GameServiceActionDeductUserExp),
		MsgBody: input,
	})

	return resFunc(true, nil)
}
