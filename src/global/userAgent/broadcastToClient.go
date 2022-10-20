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

func MultipleBroadCastToClient(fromServiceAppId string, userIds []int64, msg *proto.Envelope) error {
	agentModel := GetUserAgentModel()
	if agentModel == nil {
		return errors.New("UserAgentModel is nil")
	}

	agentList := make(map[string][]int64)
	for _, userId := range userIds {
		if agent, _ := agentModel.GetUserAgent(userId); agent != nil {
			if _, exist := agentList[agent.AgentAppId]; exist {
				agentList[agent.AgentAppId] = append(agentList[agent.AgentAppId], userId)
			} else {
				agentList[agent.AgentAppId] = []int64{userId}
			}
		}
	}

	msgBody, err := protoTool.MarshalProto(msg)
	if err != nil {
		return err
	}
	for agentId, userIds := range agentList {
		input := methodData.MultipleBroadCastToClientInput{
			MsgVersion:   time_helper.NowUTCMill(),
			ServiceAppId: fromServiceAppId,
			UserList:     userIds,
			MsgId:        int32(msg.Type),
			MsgBody:      msgBody,
		}
		inputBytes, _ := json.Marshal(input)
		_, err := daprInvoke.InvokeMethod(
			agentId,
			string(grpc.ProtoMessageActionMultipleBroadCastToClient),
			inputBytes,
		)
		if err != nil {
			serviceLog.Error(err.Error())
		}
	}
	return nil
}
