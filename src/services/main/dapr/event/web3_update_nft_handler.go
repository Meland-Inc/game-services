package daprEvent

import (
	"context"
	"fmt"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/global/grpcAPI/grpcNetTool"
	"github.com/Meland-Inc/game-services/src/global/serviceCnf"
	message "github.com/Meland-Inc/game-services/src/global/web3Message"
	"github.com/Meland-Inc/game-services/src/services/main/msgChannel"
	"github.com/dapr/go-sdk/service/common"
	"github.com/spf13/cast"
)

func Web3UpdateUserNftHandler(ctx context.Context, e *common.TopicEvent) (retry bool, err error) {
	serviceLog.Info("Receive Web3UpdateUserNft nft: %v, :%s ", e.Data, e.DataContentType)

	input := &message.UpdateUserNFT{}
	err = grpcNetTool.UnmarshalGrpcTopicEvent(e, input)
	if err != nil {
		serviceLog.Error("Web3UpdateUserNft UnmarshalEvent fail err: %v ", err)
		return false, err
	}

	if input.Etag < int(serviceCnf.GetInstance().StartMs/1000) {
		return
	}

	userId := cast.ToInt64(input.UserId)
	if userId < 1 {
		serviceLog.Error("dapr Web3UpdateUserNft invalid nft Data[%v]", input)
		return false, fmt.Errorf("dapr Web3UpdateUserNft invalid userId [%v]", input)
	}

	msgChannel.GetInstance().CallServiceMsg(&msgChannel.ServiceMsgData{
		MsgId:   string(message.SubscriptionEventUpdateUserNFT),
		MsgBody: input,
	})

	return false, nil
}

func Web3MultiUpdateUserNftHandler(ctx context.Context, e *common.TopicEvent) (retry bool, err error) {
	fmt.Printf("Receive Web3 MultiWeb3UpdateUserNft nft: %v, %s ", e.Data, e.DataContentType)

	input := &message.MultiUpdateUserNFT{}
	err = grpcNetTool.UnmarshalGrpcTopicEvent(e, input)
	if err != nil {
		serviceLog.Error("MultiWeb3UpdateUserNft UnmarshalEvent fail err: %v ", err)
		return false, err
	}

	if input.Etag < int(serviceCnf.GetInstance().StartMs/1000) {
		return
	}

	userId := cast.ToInt64(input.UserId)
	if userId < 1 {
		serviceLog.Error("web3 Multi update user nft invalid nft Data[%v]", input)
		return false, fmt.Errorf("web3 Multi update user nft invalid userId [%v]", input)
	}

	msgChannel.GetInstance().CallServiceMsg(&msgChannel.ServiceMsgData{
		MsgId:   string(message.SubscriptionEventMultiUpdateUserNFT),
		MsgBody: input,
	})

	return false, nil
}
