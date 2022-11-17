package daprEvent

import (
	"context"

	"game-message-core/grpc/pubsubEventData"
	"game-message-core/proto"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/global/grpcAPI/grpcNetTool"
	"github.com/Meland-Inc/game-services/src/global/serviceCnf"
	"github.com/Meland-Inc/game-services/src/services/agent/userChannel"
	"github.com/dapr/go-sdk/service/common"
)

func ServiceUnRegisterHandler(ctx context.Context, e *common.TopicEvent) (retry bool, err error) {
	input := &pubsubEventData.ServiceUnregisterEvent{}
	err = grpcNetTool.UnmarshalGrpcTopicEvent(e, input)
	if err != nil {
		serviceLog.Error("ServiceUnregisterEvent Unmarshal fail err: %v", err)
		return false, nil
	}

	// 抛弃过期事件
	if input.MsgVersion < serviceCnf.GetInstance().StartMs {
		return false, nil
	}

	serviceLog.Info("service UnRegister: %v", input)

	switch input.Service.ServiceType {
	case proto.ServiceType_ServiceTypeScene:
		onSceneServiceUnregister(input)
	}
	return false, nil
}

func onSceneServiceUnregister(input *pubsubEventData.ServiceUnregisterEvent) {
	inSceneUserChArr := []*userChannel.UserChannel{}
	userChannel.GetInstance().Range(
		func(userCh *userChannel.UserChannel) bool {
			if userCh.GetSceneService() == input.Service.AppId {
				inSceneUserChArr = append(inSceneUserChArr, userCh)
			}
			return true
		},
	)
	for _, userCh := range inSceneUserChArr {
		broadcastTickOutPlayer(userCh, proto.TickOutType_ServiceClose)
		userCh.GetSession().Stop()
	}
}
