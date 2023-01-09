package serviceHandler

import (
	"game-message-core/grpc/pubsubEventData"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/global/contract"
	"github.com/Meland-Inc/game-services/src/global/serviceCnf"
	"github.com/Meland-Inc/game-services/src/global/userAgent"
	"github.com/Meland-Inc/game-services/src/services/chat/chatModel"
)

func GRPCUserEnterGameEvent(env contract.IModuleEventReq, curMs int64) {
	input := &pubsubEventData.UserEnterGameEvent{}
	err := env.UnmarshalToDaprEventData(input)
	if err != nil {
		serviceLog.Error("UserEnterGameEvent UnmarshalEvent fail err: %v ", err)
		return
	}

	if input.MsgVersion < serviceCnf.GetInstance().StartMs {
		return
	}

	serviceLog.Info("Receive UserEnterGameEvent: %+v", input)
	if input.UserId <= 0 {
		serviceLog.Error("invalid userId [%d] in enter game event")
		return
	}

	agentModel := userAgent.GetUserAgentModel()
	if agent, err := agentModel.CheckAndAddUserAgentRecord(
		input.UserId, input.AgentAppId, input.UserSocketId, input.SceneServiceAppId,
	); err != nil {
		serviceLog.Error(err.Error())
	} else {
		agent.InMapId = input.MapId
	}

	model, _ := chatModel.GetChatModel()
	if err := model.OnPlayerEnterGame(input); err != nil {
		serviceLog.Error(err.Error())
	}
}

func GRPCUserLeaveGameEvent(env contract.IModuleEventReq, curMs int64) {
	input := &pubsubEventData.UserLeaveGameEvent{}
	err := env.UnmarshalToDaprEventData(input)
	if err != nil {
		serviceLog.Error("UserLeaveGameEvent UnmarshalEvent fail err: %v ", err)
		return
	}

	if input.MsgVersion < serviceCnf.GetInstance().StartMs {
		return
	}

	serviceLog.Info("service receive LeaveGame: %+v", input)

	agentModel := userAgent.GetUserAgentModel()
	agentModel.RemoveUserAgentRecord(input.UserId)

	model, _ := chatModel.GetChatModel()
	model.OnPlayerLeaveGame(input.UserId)
}

func GRPCTickOutUserEvent(env contract.IModuleEventReq, curMs int64) {
	input := &pubsubEventData.TickOutPlayerEvent{}
	err := env.UnmarshalToDaprEventData(input)
	if err != nil {
		serviceLog.Error("TickOutPlayerEvent UnmarshalEvent fail err: %v ", err)
		return
	}

	if input.MsgVersion < serviceCnf.GetInstance().StartMs {
		return
	}

	serviceLog.Info("service receive TickOutPlayerEvent: %+v", input)

	agentModel := userAgent.GetUserAgentModel()
	agentModel.RemoveUserAgentRecord(input.UserId)

	model, _ := chatModel.GetChatModel()
	model.OnPlayerLeaveGame(input.UserId)
}

func GRPCSavePlayerDataEvent(env contract.IModuleEventReq, curMs int64) {
	input := &pubsubEventData.SavePlayerEventData{}
	err := env.UnmarshalToDaprEventData(input)
	if err != nil {
		serviceLog.Error("SavePlayerDataEvent UnmarshalEvent fail err: %v ", err)
		return
	}

	serviceLog.Info("receive savePlayerEvent: %+v", input)

	model, _ := chatModel.GetChatModel()
	model.OnUpdatePlayerData(input)
}
