package daprEvent

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/global/serviceCnf"
	message "github.com/Meland-Inc/game-services/src/global/web3Message"
	"github.com/Meland-Inc/game-services/src/services/main/msgChannel"
	"github.com/dapr/go-sdk/service/common"
	"github.com/spf13/cast"
)

func Web3UpdateUserNftHandler(ctx context.Context, e *common.TopicEvent) (retry bool, err error) {
	serviceLog.Info("Receive Web3 update user nft: %v %s \n", e.Data, e.DataContentType)

	bs, err := json.Marshal(e.Data)
	input := message.UpdateUserNFT{}
	err = input.UnmarshalJSON(bs)
	if err != nil {
		serviceLog.Error("not math to dapr msg UpdateUserNFT data : %+v", e.Data)
		return false, fmt.Errorf("not math to dapr msg UpdateUserNFT")
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
	fmt.Printf("Receive Web3 Multi update user nft: %v %s \n", e.Data, e.DataContentType)

	bs, err := json.Marshal(e.Data)
	input := message.MultiUpdateUserNFT{}
	err = input.UnmarshalJSON(bs)
	if err != nil {
		serviceLog.Error("not math to dapr msg MultiUpdateUserNFT data : %+v", e.Data)
		return false, fmt.Errorf("not math to dapr msg MultiUpdateUserNFT")
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
