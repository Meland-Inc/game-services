package daprEvent

import (
	"context"
	"encoding/json"
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
	serviceLog.Info("Receive Web3 update user nft: %v, :%s ", e.Data, e.DataContentType)

	inputBytes, err := json.Marshal(e.Data)
	if err != nil {
		serviceLog.Error("Web3UpdateUserNftHandler  marshal e.Data  fail err: %+v", err)
		return false, fmt.Errorf("Web3UpdateUserNftHandler  marshal e.Data  fail err: %+v", err)
	}

	input := &message.UpdateUserNFT{}
	err = grpcNetTool.UnmarshalGrpcData(inputBytes, input)
	if err != nil {
		return false, err
	}

	if input.Etag < int(serviceCnf.GetInstance().StartMs) {
		return
	}

	userId := cast.ToInt64(input.UserId)
	if userId < 1 {
		serviceLog.Error("dapr update user nft invalid nft Data[%v]", input)
		return false, fmt.Errorf("dapr update user nft invalid userId [%v]", input)
	}

	msgChannel.GetInstance().CallServiceMsg(&msgChannel.ServiceMsgData{
		MsgId:   string(message.SubscriptionEventUpdateUserNFT),
		MsgBody: input,
	})

	return false, nil
}

func Web3MultiUpdateUserNftHandler(ctx context.Context, e *common.TopicEvent) (retry bool, err error) {
	fmt.Printf("Receive Web3 Multi update user nft: %v, %s ", e.Data, e.DataContentType)

	inputBytes, err := json.Marshal(e.Data)
	if err != nil {
		serviceLog.Error("Web3MultiUpdateUserNftHandler  marshal e.Data  fail err: %+v", err)
		return false, fmt.Errorf("Web3MultiUpdateUserNftHandler  marshal e.Data  fail err: %+v", err)
	}

	input := &message.MultiUpdateUserNFT{}
	err = grpcNetTool.UnmarshalGrpcData(inputBytes, input)
	if err != nil {
		return false, err
	}

	if input.Etag < int(serviceCnf.GetInstance().StartMs) {
		return
	}

	userId := cast.ToInt64(input.UserId)
	if userId < 1 {
		serviceLog.Error("web3 Multi update user nft invalid nft Data[%v]", input)
		return false, fmt.Errorf("web3 Multi update user nft invalid userId [%v]", input)
	}

	msgChannel.GetInstance().CallServiceMsg(&msgChannel.ServiceMsgData{
		MsgId:   string(message.SubscriptionEventMultiUpdateUserNFT),
		MsgBody: input.Nfts,
	})

	return false, nil
}
