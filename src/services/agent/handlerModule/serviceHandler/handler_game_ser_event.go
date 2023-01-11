package serviceHandler

import (
	"game-message-core/grpc/pubsubEventData"
	"game-message-core/proto"
	"game-message-core/protoTool"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/global/contract"
	"github.com/Meland-Inc/game-services/src/global/serviceCnf"
	"github.com/Meland-Inc/game-services/src/services/agent/userChannel"
)

func broadcastTickOutPlayer(userCh *userChannel.UserChannel, tickCode proto.TickOutType) {
	if userCh == nil {
		return
	}
	msg := &proto.Envelope{
		Type: proto.EnvelopeType_BroadCastTickOut,
		Payload: &proto.Envelope_BroadCastTickOutResponse{
			BroadCastTickOutResponse: &proto.BroadCastTickOutResponse{
				Kind: tickCode,
			},
		},
	}
	msgBody, err := protoTool.MarshalEnvelope(msg)
	if err != nil {
		serviceLog.Error(err.Error())
		return
	}
	serviceLog.Debug("broad cast tick out user[%v][%v][%v]", userCh.GetOwner(), tickCode, userCh.GetSession().SessionId())
	userCh.SendToUser(msg.Type, msgBody)
}
func GRPCTickOutUserEvent(env contract.IModuleEventReq, curMs int64) {
	input := &pubsubEventData.TickOutPlayerEvent{}
	err := env.UnmarshalToDaprEventData(input)
	if err != nil {
		serviceLog.Error("TickOutPlayer UnmarshalEvent fail err: %v ", err)
		return
	}

	if input.MsgVersion < serviceCnf.GetInstance().StartMs {
		return
	}

	serviceLog.Info("agent service receive TickOutPlayer: %+v", input)

	var userCh *userChannel.UserChannel
	if input.SocketId != "" {
		userCh = userChannel.GetInstance().UserChannelById(input.SocketId)
	} else if input.UserId > 0 {
		userCh = userChannel.GetInstance().UserChannelByOwner(input.UserId)
	}
	if userCh == nil {
		return
	}

	broadcastTickOutPlayer(userCh, input.TickOutCode)
	userCh.Stop()
	userCh.Stop()
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
func GRPCServiceUnRegisterEvent(env contract.IModuleEventReq, curMs int64) {
	input := &pubsubEventData.ServiceUnregisterEvent{}
	err := env.UnmarshalToDaprEventData(input)
	if err != nil {
		serviceLog.Error("ServiceUnregisterEvent UnmarshalEvent fail err: %v ", err)
		return
	}

	if input.MsgVersion < serviceCnf.GetInstance().StartMs {
		return
	}

	serviceLog.Info("service UnRegister: %v", input)

	switch input.Service.ServiceType {
	case proto.ServiceType_ServiceTypeScene:
		onSceneServiceUnregister(input)
	}
}

func GRPCChangeServiceEvent(env contract.IModuleEventReq, curMs int64) {
	input := &pubsubEventData.UserChangeServiceEvent{}
	err := env.UnmarshalToDaprEventData(input)
	if err != nil {
		serviceLog.Error("UserChangeServiceEvent UnmarshalEvent fail err: %v ", err)
		return
	}

	if input.MsgVersion < serviceCnf.GetInstance().StartMs {
		return
	}

	// 抛弃过期事件
	if input.MsgVersion < serviceCnf.GetInstance().StartMs {
		return
	}

	serviceLog.Info("user change service: %v", input)

	userAgent := userChannel.GetInstance().UserChannelByOwner(input.UserId)
	if userAgent != nil {
		userAgent.SetSceneService(input.ToService.AppId)
		serviceLog.Debug("user change scene service new data = %v", userAgent.GetSceneService())
	}

}
