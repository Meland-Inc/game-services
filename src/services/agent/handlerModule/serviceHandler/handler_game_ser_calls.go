package serviceHandler

import (
	"fmt"
	"game-message-core/grpc/methodData"
	"game-message-core/proto"
	"game-message-core/protoTool"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/global/contract"
	"github.com/Meland-Inc/game-services/src/global/grpcAPI/grpcPubsubEvent"
	"github.com/Meland-Inc/game-services/src/global/module"
	"github.com/Meland-Inc/game-services/src/services/agent/clientMsgLogCtrl"
	"github.com/Meland-Inc/game-services/src/services/agent/userChannel"
)

func getUserChannel(userId int64, socketId string) *userChannel.UserChannel {
	var userCh *userChannel.UserChannel
	if socketId != "" {
		userCh = userChannel.GetInstance().UserChannelById(socketId)
	} else if userId > 0 {
		userCh = userChannel.GetInstance().UserChannelByOwner(userId)
	}
	return userCh
}

func BroadCastToClientHandler(env contract.IModuleEventReq, curMs int64) {
	output := &methodData.BroadCastToClientOutput{Success: true}
	result := &module.ModuleEventResult{}
	defer func() {
		if output.ErrMsg != "" {
			output.Success = false
			serviceLog.Warning("BroadCastToClient fail err: %v", output.ErrMsg)
		}
		result.SetResult(output)
		env.WriteResult(result)
	}()

	input := &methodData.BroadCastToClientInput{}
	err := env.UnmarshalToDaprCallData(input)
	if err != nil {
		output.ErrMsg = err.Error()
		return
	}

	resMsg, err := protoTool.UnMarshalToEnvelope(input.MsgBody)
	if clientMsgLogCtrl.PrintCliMsgLog(resMsg.Type) {
		serviceLog.Info("BroadCastToClient user[%d] msg[%v], err:%+v", input.UserId, resMsg.Type, err)
	}

	userCh := getUserChannel(input.UserId, input.SocketId)
	if userCh == nil {
		output.ErrMsg = fmt.Sprintf("BroadCastToClient userCh not found  userId[%d], socketId[%v]", input.UserId, input.SocketId)
		serviceLog.Error(output.ErrMsg)
		if input.UserId > 0 {
			grpcPubsubEvent.RPCPubsubEventLeaveGame(input.UserId)
		}
		return
	}

	userCh.SendToUser(proto.EnvelopeType(input.MsgId), input.MsgBody)
}

func MultipleBroadCastToClientHandler(env contract.IModuleEventReq, curMs int64) {
	output := &methodData.MultipleBroadCastToClientOutput{Success: true}
	result := &module.ModuleEventResult{}
	defer func() {
		if output.ErrMsg != "" {
			output.Success = false
			serviceLog.Warning("MultipleBroadCastToClient fail err: %v", output.ErrMsg)
		}
		result.SetResult(output)
		env.WriteResult(result)
	}()

	input := &methodData.MultipleBroadCastToClientInput{}
	err := env.UnmarshalToDaprCallData(input)
	if err != nil {
		output.ErrMsg = err.Error()
		return
	}

	resMsg, err := protoTool.UnMarshalToEnvelope(input.MsgBody)
	if clientMsgLogCtrl.PrintCliMsgLog(resMsg.Type) {
		serviceLog.Info("MultipleBroadCastToClient Users:%v, msg[%+v], err:%+v", input.UserList, resMsg.Type, err)
	}

	for _, userId := range input.UserList {
		if userId <= 0 {
			continue
		}
		userCh := getUserChannel(userId, "")
		if userCh != nil {
			userCh.SendToUser(proto.EnvelopeType(input.MsgId), input.MsgBody)
		} else {
			serviceLog.Warning("UserChannel [%d] not found", userId)
			grpcPubsubEvent.RPCPubsubEventLeaveGame(userId)
		}
	}
}
