package clientMsgHandle

import (
	"game-message-core/grpc/methodData"
	"game-message-core/proto"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/common/time_helper"
	"github.com/Meland-Inc/game-services/src/global/serviceCnf"
	"github.com/Meland-Inc/game-services/src/global/userAgent"
)

func ResponseClientMessage(
	agent *userAgent.UserAgentData,
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
	serviceLog.Info("response player:[%d], msg:[%s]", input.UserId, respMsg.Type)
	agent.SendToPlayer(serviceCnf.GetInstance().ServerName, respMsg)
}

func makeResponseMsg(msg *proto.Envelope) *proto.Envelope {
	return &proto.Envelope{
		Type:  msg.Type,
		SeqId: msg.SeqId,
	}
}

func getPlayerAgent(input *methodData.PullClientMessageInput) *userAgent.UserAgentData {
	agentModel := userAgent.GetUserAgentModel()
	agent, exist := agentModel.GetUserAgent(input.UserId)
	if !exist {
		agent = &userAgent.UserAgentData{
			AgentAppId: input.AgentAppId,
			SocketId:   input.SocketId,
			UserId:     input.UserId,
			LoginAt:    time_helper.NowUTCMill(),
		}
		agentModel.AddUserAgentRecord(input.UserId, input.AgentAppId, input.SocketId)
	} else {
		agent.TryUpdate(input.UserId, input.AgentAppId, input.SocketId)
	}

	return agent
}

func GetOrStoreUserAgent(input *methodData.PullClientMessageInput) *userAgent.UserAgentData {
	agentModel := userAgent.GetUserAgentModel()
	agent, exist := agentModel.GetUserAgent(input.UserId)
	if !exist {
		agent, _ = agentModel.AddUserAgentRecord(input.UserId, input.AgentAppId, input.SocketId)
	} else {
		agent.TryUpdate(input.UserId, input.AgentAppId, input.SocketId)
	}
	return agent
}
