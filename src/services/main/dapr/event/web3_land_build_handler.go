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
	input := &message.MultiLandDataUpdateEvent{}
	err = grpcNetTool.UnmarshalGrpcTopicEvent(e, input)
	if err != nil {
		serviceLog.Error("Web3MultiLandDataUpdateEvent Unmarshal fail err: %v ", err)
		return false, nil
	}

	if input.Etag < int(serviceCnf.GetInstance().StartMs/1000) {
		return false, nil
	}

	serviceLog.Info("Receive Web3MultiLandDataUpdateEvent: %+v", input)

	msgChannel.GetInstance().CallServiceMsg(&msgChannel.ServiceMsgData{
		MsgId:   string(message.SubscriptionEventMultiLandDataUpdateEvent),
		MsgBody: input,
	})

	return false, nil
}

func Web3MultiRecyclingEventHandler(ctx context.Context, e *common.TopicEvent) (retry bool, err error) {
	input := &message.MultiRecyclingEvent{}
	err = grpcNetTool.UnmarshalGrpcTopicEvent(e, input)
	if err != nil {
		serviceLog.Error("Web3RecyclingEvent Unmarshal fail err: %v ", err)
		return false, nil
	}

	if input.Etag < int(serviceCnf.GetInstance().StartMs/1000) {
		return false, nil
	}

	serviceLog.Info("Receive Web3RecyclingEvent: %+v", input)

	msgChannel.GetInstance().CallServiceMsg(&msgChannel.ServiceMsgData{
		MsgId:   string(message.SubscriptionEventMultiRecyclingEvent),
		MsgBody: input,
	})

	return false, nil
}

func Web3MultiBuildUpdateEventHandler(ctx context.Context, e *common.TopicEvent) (retry bool, err error) {
	input := &message.MultiBuildUpdateEvent{}
	err = grpcNetTool.UnmarshalGrpcTopicEvent(e, input)
	if err != nil {
		serviceLog.Error("Web3BuildUpdateEvent UnmarshalEvent fail err: %v ", err)
		return false, err
	}

	if input.Etag < int(serviceCnf.GetInstance().StartMs/1000) {
		return false, nil
	}

	serviceLog.Info("Receive Web3BuildUpdateEvent: %v", input)

	msgChannel.GetInstance().CallServiceMsg(&msgChannel.ServiceMsgData{
		MsgId:   string(message.SubscriptionEventMultiBuildUpdateEvent),
		MsgBody: input,
	})

	return false, nil
}
