package userAgent

import (
	"encoding/json"
	"errors"
	"game-message-core/grpc"
	"game-message-core/grpc/methodData"
	"game-message-core/proto"
	"game-message-core/protoTool"

	"github.com/Meland-Inc/game-services/src/common/daprInvoke"
	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/common/time_helper"
)

func BroadCastToClient(
	agentAppId, serviceAppId string, userId int64,
	userSocketId string, msg *proto.Envelope,
) error {
	if msg == nil {
		return errors.New("msg is nil")
	}

	msgBody, err := protoTool.MarshalProto(msg)
	if err != nil {
		return err
	}

	input := methodData.BroadCastToClientInput{
		MsgVersion:   time_helper.NowUTCMill(),
		ServiceAppId: serviceAppId,
		UserId:       userId,
		SocketId:     userSocketId,
		MsgId:        int32(msg.Type),
		MsgBody:      msgBody,
	}

	inputBytes, err := json.Marshal(input)
	if err != nil {
		serviceLog.Error("SendToPlayer Marshal BroadCastInput failed err: %+v", err)
		return err
	}

	_, err = daprInvoke.InvokeMethod(
		agentAppId,
		string(grpc.ProtoMessageActionBroadCastToClient),
		inputBytes,
	)
	return err

}

func MultipleBroadCastToClient(agentAppId, serviceAppId string, userIds []int64, msg *proto.Envelope) error {
	if msg == nil {
		return errors.New("msg is nil")
	}

	msgBody, err := protoTool.MarshalProto(msg)
	if err != nil {
		return err
	}

	input := methodData.MultipleBroadCastToClientInput{
		MsgVersion:   time_helper.NowUTCMill(),
		ServiceAppId: serviceAppId,
		UserList:     userIds,
		MsgId:        int32(msg.Type),
		MsgBody:      msgBody,
	}

	inputBytes, err := json.Marshal(input)
	if err != nil {
		serviceLog.Error("SendToPlayer Marshal MultipleBroadCastInput failed err: %+v", err)
		return err
	}

	_, err = daprInvoke.InvokeMethod(
		agentAppId,
		string(grpc.ProtoMessageActionMultipleBroadCastToClient),
		inputBytes,
	)
	return err
}
