package userAgent

import (
	"encoding/json"
	"game-message-core/grpc"
	"game-message-core/grpc/methodData"
	"game-message-core/proto"
	"game-message-core/protoTool"

	"github.com/Meland-Inc/game-services/src/common/daprInvoke"
	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/common/time_helper"
)

type UserAgentData struct {
	AgentAppId string `json:"agentAppId"`
	SocketId   string `json:"socketId"`
	UserId     int64  `json:"userId"`
	LoginAt    int64  `json:"loginAt"`
}

func (p *UserAgentData) SendToPlayer(serviceAppId string, msg *proto.Envelope) error {
	msgBody, err := protoTool.MarshalProto(msg)
	if err != nil {
		return err
	}

	input := methodData.BroadCastToClientInput{
		MsgVersion:   time_helper.NowUTCMill(),
		ServiceAppId: serviceAppId,
		UserId:       p.UserId,
		SocketId:     p.SocketId,
		MsgId:        int32(msg.Type),
		MsgBody:      msgBody,
	}

	inputBytes, err := json.Marshal(input)
	if err != nil {
		serviceLog.Error("SendToPlayer Marshal BroadCastInput failed err: %+v", err)
		return err
	}

	_, err = daprInvoke.InvokeMethod(
		p.AgentAppId,
		string(grpc.ProtoMessageActionBroadCastToClient),
		inputBytes,
	)
	if err != nil {
		serviceLog.Error("UserAgentData SendToPlayer InvokeMethod  failed err:%+v", err)
	}
	return err
}