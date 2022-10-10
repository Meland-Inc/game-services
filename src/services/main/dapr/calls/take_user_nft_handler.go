package daprCalls

import (
	"context"
	"fmt"
	"game-message-core/grpc"
	"game-message-core/grpc/methodData"

	"github.com/Meland-Inc/game-services/src/common/daprInvoke"
	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/global/grpcAPI/grpcNetTool"
	"github.com/Meland-Inc/game-services/src/services/main/msgChannel"
	"github.com/Meland-Inc/game-services/src/services/main/playerModel"
	"github.com/dapr/go-sdk/service/common"
)

func TakeUserNftHandler(ctx context.Context, in *common.InvocationEvent) (*common.Content, error) {
	resFunc := func(success bool, err error) (*common.Content, error) {
		out := &methodData.MainServiceActionTakeNftOutput{}
		out.Success = success
		if err != nil {
			out.FailedMsg = err.Error()
			serviceLog.Error("take user nft err: %v", err)
		}
		content, _ := daprInvoke.MakeOutputContent(in, out)
		return content, err
	}

	serviceLog.Info("main service receive take nft, data: %v", string(in.Data))

	input := &methodData.MainServiceActionTakeNftInput{}
	err := grpcNetTool.UnmarshalGrpcData(in.Data, input)
	if err != nil {
		return resFunc(false, err)
	}

	if input.UserId < 1 {
		return resFunc(false, fmt.Errorf("invalid user id: %d", input.UserId))
	}

	if err = checkNfts(input.UserId, input.TakeNfts); err != nil {
		return resFunc(false, err)
	}

	msgChannel.GetInstance().CallServiceMsg(&msgChannel.ServiceMsgData{
		MsgId:   string(grpc.MainServiceActionTakeNFT),
		MsgBody: input,
	})

	return resFunc(true, nil)
}

func checkNfts(userId int64, takeNfts []methodData.TakeNftData) error {
	dataModel, err := playerModel.GetPlayerDataModel()
	if err != nil {
		serviceLog.Error("main service take nft playerDataModel not found")
		return err
	}

	playerItem, err := dataModel.GetPlayerItems(userId)
	if err != nil {
		return err
	}

	for _, tn := range takeNfts {
		var giveCount = tn.Num
		for _, item := range playerItem.Items {
			if tn.NftId != "" && tn.NftId != item.Id {
				continue
			}
			if tn.ItemCid != 0 && tn.ItemCid != item.Cid {
				continue
			}
			giveCount -= item.Num
			if giveCount <= 0 {
				break
			}
		}
		if giveCount > 0 {
			return fmt.Errorf("not fund NFT %+v", tn)
		}
	}
	return nil
}
