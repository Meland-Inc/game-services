package userAgent

import (
	"game-message-core/grpc/methodData"
	"game-message-core/proto"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/common/time_helper"
	"github.com/Meland-Inc/game-services/src/global/serviceCnf"
)

func ResponseClientMessage(
	agent *UserAgentData,
	input *methodData.PullClientMessageInput,
	respMsg *proto.Envelope,
) {
	if respMsg.ErrorMessage != "" {
		serviceLog.Error(
			"responseClient [%v] Msg err : [%d][%s]",
			respMsg.Type, respMsg.ErrorCode, respMsg.ErrorMessage,
		)
	}
	if agent == nil {
		serviceLog.Error("player[%d] agent not found", input.UserId)
		return
	}
	agent.SendToPlayer(serviceCnf.GetInstance().AppId, respMsg)
}

func MakeResponseMsg(msg *proto.Envelope) *proto.Envelope {
	return &proto.Envelope{
		Type:  msg.Type,
		SeqId: msg.SeqId,
	}
}

func getPlayerAgent(input *methodData.PullClientMessageInput) *UserAgentData {
	agentModel := GetUserAgentModel()
	agent, exist := agentModel.GetUserAgent(input.UserId)
	if !exist {
		agent = &UserAgentData{
			AgentAppId:          input.AgentAppId,
			SocketId:            input.SocketId,
			InSceneServiceAppId: input.SceneServiceId,
			UserId:              input.UserId,
			LoginAt:             time_helper.NowUTCMill(),
		}
		agentModel.AddUserAgentRecord(input.UserId, input.AgentAppId, input.SocketId, input.SceneServiceId)
	} else {
		agent.TryUpdate(input.UserId, input.AgentAppId, input.SocketId, input.SceneServiceId)
	}

	return agent
}

func GetOrStoreUserAgent(input *methodData.PullClientMessageInput) *UserAgentData {
	agentModel := GetUserAgentModel()
	agent, exist := agentModel.GetUserAgent(input.UserId)
	if !exist {
		agent, _ = agentModel.AddUserAgentRecord(
			input.UserId,
			input.AgentAppId,
			input.SocketId,
			input.SceneServiceId,
		)
	} else {
		agent.TryUpdate(input.UserId, input.AgentAppId, input.SocketId, input.SceneServiceId)
	}
	return agent
}
