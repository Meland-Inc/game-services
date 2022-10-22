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

func Web3RecyclingEventHandler(ctx context.Context, e *common.TopicEvent) (retry bool, err error) {
	serviceLog.Info("Receive Web3RecyclingEvent data: %v", e.Data)

	input := &message.RecyclingEvent{}
	err = grpcNetTool.UnmarshalGrpcTopicEvent(e, input)
	if err != nil {
		serviceLog.Error("Web3RecyclingEvent UnmarshalEvent fail err: %v ", err)
		return false, err
	}

	if input.Etag < int(serviceCnf.GetInstance().StartMs/1000) {
		return false, nil
	}

	msgChannel.GetInstance().CallServiceMsg(&msgChannel.ServiceMsgData{
		MsgId:   string(message.SubscriptionEventRecyclingEvent),
		MsgBody: input,
	})

	return false, nil
}

func Web3BuildUpdateEventHandler(ctx context.Context, e *common.TopicEvent) (retry bool, err error) {
	serviceLog.Info("Receive Web3BuildUpdateEvent data: %v", e.Data)

	input := &message.BuildUpdateEvent{}
	err = grpcNetTool.UnmarshalGrpcTopicEvent(e, input)
	if err != nil {
		serviceLog.Error("Web3BuildUpdateEvent UnmarshalEvent fail err: %v ", err)
		return false, err
	}

	if input.Etag < int(serviceCnf.GetInstance().StartMs/1000) {
		return false, nil
	}

	msgChannel.GetInstance().CallServiceMsg(&msgChannel.ServiceMsgData{
		MsgId:   string(message.SubscriptionEventBuildUpdateEvent),
		MsgBody: input,
	})

	return false, nil
}
