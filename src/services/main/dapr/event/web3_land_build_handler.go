package daprEvent

import (
	"context"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/global/grpcAPI/grpcNetTool"
	"github.com/Meland-Inc/game-services/src/global/serviceCnf"
	message "github.com/Meland-Inc/game-services/src/global/web3Message"
	"github.com/Meland-Inc/game-services/src/services/main/msgChannel"
	"github.com/dapr/go-sdk/service/common"
)

func Web3MultiLandDataUpdateEventHandler(ctx context.Context, e *common.TopicEvent) (retry bool, err error) {
	serviceLog.Info("Receive Web3MultiLandDataUpdateEvent data: %v", e.Data)

	input := &message.MultiLandDataUpdateEvent{}
	err = grpcNetTool.UnmarshalGrpcTopicEvent(e, input)
	if err != nil {
		serviceLog.Error("Web3MultiLandDataUpdateEvent UnmarshalEvent fail err: %v ", err)
		return false, err
	}

	if input.Etag < int(serviceCnf.GetInstance().StartMs/1000) {
		return false, nil
	}

	msgChannel.GetInstance().CallServiceMsg(&msgChannel.ServiceMsgData{
		MsgId:   string(message.SubscriptionEventMultiLandDataUpdateEvent),
		MsgBody: input,
	})

	return false, nil
}

func Web3MultiRecyclingEventHandler(ctx context.Context, e *common.TopicEvent) (retry bool, err error) {
	serviceLog.Info("Receive Web3RecyclingEvent data: %v", e.Data)

	input := &message.MultiRecyclingEvent{}
	err = grpcNetTool.UnmarshalGrpcTopicEvent(e, input)
	if err != nil {
		serviceLog.Error("Web3RecyclingEvent UnmarshalEvent fail err: %v ", err)
		return false, err
	}

	if input.Etag < int(serviceCnf.GetInstance().StartMs/1000) {
		return false, nil
	}

	msgChannel.GetInstance().CallServiceMsg(&msgChannel.ServiceMsgData{
		MsgId:   string(message.SubscriptionEventMultiRecyclingEvent),
		MsgBody: input,
	})

	return false, nil
}

func Web3MultiBuildUpdateEventHandler(ctx context.Context, e *common.TopicEvent) (retry bool, err error) {
	serviceLog.Info("Receive Web3BuildUpdateEvent data: %v", e.Data)

	input := &message.MultiBuildUpdateEvent{}
	err = grpcNetTool.UnmarshalGrpcTopicEvent(e, input)
	if err != nil {
		serviceLog.Error("Web3BuildUpdateEvent UnmarshalEvent fail err: %v ", err)
		return false, err
	}

	if input.Etag < int(serviceCnf.GetInstance().StartMs/1000) {
		return false, nil
	}

	msgChannel.GetInstance().CallServiceMsg(&msgChannel.ServiceMsgData{
		MsgId:   string(message.SubscriptionEventMultiBuildUpdateEvent),
		MsgBody: input,
	})

	return false, nil
}
