package clientMsgHandle

import (
	"fmt"
	"game-message-core/grpc/methodData"
	"game-message-core/proto"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/common/time_helper"
	"github.com/Meland-Inc/game-services/src/global/component"
	"github.com/Meland-Inc/game-services/src/global/serviceCnf"
	"github.com/Meland-Inc/game-services/src/global/userAgent"
	"github.com/Meland-Inc/game-services/src/services/main/playerModel"
)

func ResponseClientMessage(
	input *methodData.PullClientMessageInput,
	respMsg *proto.Envelope,
) {
	if respMsg.ErrorMessage != "" {
		serviceLog.Error(
			"responseClient [%v] Msg err : [%d][%s]",
			respMsg.Type, respMsg.ErrorCode, respMsg.ErrorMessage,
		)
	}

	agent := getPlayerAgent(input)
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

func getPlayerDataModel() (*playerModel.PlayerDataModel, error) {
	iPlayerModel, exist := component.GetInstance().GetModel(component.MODEL_NAME_PLAYER_DATA)
	if !exist {
		return nil, fmt.Errorf("player data model not found")
	}
	dataModel, _ := iPlayerModel.(*playerModel.PlayerDataModel)
	return dataModel, nil
}

func getPlayerAgent(input *methodData.PullClientMessageInput) *userAgent.UserAgentData {
	iUserAgentModel, exist := component.GetInstance().GetModel(component.MODEL_NAME_USER_AGENT)
	if !exist {
		return nil
	}
	agentModel := iUserAgentModel.(*userAgent.UserAgentModel)
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
