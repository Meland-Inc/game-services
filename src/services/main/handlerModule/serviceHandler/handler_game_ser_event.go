package serviceHandler

import (
	"fmt"
	"game-message-core/grpc/pubsubEventData"
	"game-message-core/proto"
	"time"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/global/contract"
	"github.com/Meland-Inc/game-services/src/global/grpcAPI/grpcInvoke"
	"github.com/Meland-Inc/game-services/src/global/serviceCnf"
	"github.com/Meland-Inc/game-services/src/global/userAgent"
	"github.com/Meland-Inc/game-services/src/services/main/home_model"
	login_model "github.com/Meland-Inc/game-services/src/services/main/loginModel"
	"github.com/Meland-Inc/game-services/src/services/main/playerModel"
)

func GRPCSaveHomeDataEvent(env contract.IModuleEventReq, curMs int64) {
	input := &pubsubEventData.SaveHomeEvent{}
	err := env.UnmarshalToDaprEventData(input)
	if err != nil {
		serviceLog.Error("SaveHomeDataEvent UnmarshalEvent fail err: %v ", err)
		return
	}

	if input.MsgVersion < serviceCnf.GetInstance().StartMs {
		return
	}

	serviceLog.Info("Receive SaveHomeDataEvent: userId [%d]", input.UserId)
	homeModel, _ := home_model.GetHomeModel()
	err = homeModel.UpdateUserHomeData(input.UserId, input.Data)
	if err != nil {
		serviceLog.Error("SaveHomeDataEvent up user home data failed err: %v ", err)
	}
}

func GRPCGranaryStockpileEvent(env contract.IModuleEventReq, curMs int64) {
	input := &pubsubEventData.GranaryStockpileEvent{}
	err := env.UnmarshalToDaprEventData(input)
	if err != nil {
		serviceLog.Error("GranaryStockpileEvent UnmarshalEvent fail err: %v ", err)
		return
	}

	serviceLog.Info("Receive GranaryStockpileEvent: %+v", input)

	homeModel, _ := home_model.GetHomeModel()

	upTm := time.UnixMilli(input.MsgVersion).UTC()
	for _, it := range input.Items {
		if err = homeModel.TryAddGranaryRecord(
			input.HomeOwner, it.Cid, it.Num, it.Quality, upTm, input.OccupantId, input.OccupantName,
		); err != nil {
			serviceLog.Error(err.Error())
		}
	}
}

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
	agentModel := userAgent.GetUserAgentModel()
	agent, exist := agentModel.GetUserAgent(input.UserId)
	if exist {
		agent.InSceneServiceAppId = input.SceneServiceAppId
		agent.SocketId = input.UserSocketId
		agent.AgentAppId = input.AgentAppId
		agent.InMapId = input.MapId
	} else {
		agent, _ = agentModel.AddUserAgentRecord(
			input.UserId,
			input.AgentAppId,
			input.UserSocketId,
			input.SceneServiceAppId,
		)
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

	loginModel, _ := login_model.GetLoginModel()
	loginModel.OnLogOut(input.UserId)
}

func GRPCSavePlayerDataEvent(env contract.IModuleEventReq, curMs int64) {
	input := &pubsubEventData.SavePlayerEventData{}
	err := env.UnmarshalToDaprEventData(input)
	if err != nil {
		serviceLog.Error("SavePlayerDataEvent UnmarshalEvent fail err: %v ", err)
		return
	}

	serviceLog.Info("receive savePlayerEvent: %+v", input)

	playerDataModel, _ := playerModel.GetPlayerDataModel()
	sceneData, err := playerDataModel.GetPlayerSceneData(input.UserId)
	if err != nil {
		serviceLog.Error("SavePlayerEvent scene Data  not found")
		return
	}

	switch input.FormService.SceneSerSubType {
	case proto.SceneServiceSubType_World:
		sceneData.Hp = input.CurHP
		sceneData.MapId = input.FormService.MapId
		sceneData.X = input.PosX
		sceneData.Y = input.PosY
		sceneData.Z = input.PosZ
		sceneData.DirX = input.DirX
		sceneData.DirY = input.DirY
		sceneData.DirZ = input.DirZ
		err = playerDataModel.UpPlayerSceneData(sceneData)
	case proto.SceneServiceSubType_Dungeon, proto.SceneServiceSubType_Home:
		sceneData.Hp = input.CurHP
		err = playerDataModel.UpPlayerSceneData(sceneData)
	default:
		err = fmt.Errorf("invalid service sub type %v", input.FormService.SceneSerSubType)
	}
	if err != nil {
		serviceLog.Error(err.Error())
	}
}

func GRPCKillMonsterEvent(env contract.IModuleEventReq, curMs int64) {
	input := &pubsubEventData.KillMonsterEventData{}
	err := env.UnmarshalToDaprEventData(input)
	if err != nil {
		serviceLog.Error("KillMonsterEvent UnmarshalEvent fail err: %v ", err)
		return
	}

	if input.MsgVersion < serviceCnf.GetInstance().StartMs {
		return
	}

	serviceLog.Info("service receive KillMonsterEvent: %+v", input)
	playerDataModel, _ := playerModel.GetPlayerDataModel()
	err = playerDataModel.AddExp(input.UserId, input.Exp)
	if err != nil {
		serviceLog.Error("KillMonsterEvent add exp failed: %v", err)
	}
	for _, drop := range input.DropList {
		if err := grpcInvoke.Web3MintNFT(
			input.UserId, drop.Cid, drop.Num, drop.Quality, int32(input.PosX), int32(input.PosZ),
		); err != nil {
			serviceLog.Error("mint nft[%d] failed: %v", drop.Cid, err)
		}
	}
}

func GRPCPlayerDeathEvent(env contract.IModuleEventReq, curMs int64) {
	input := &pubsubEventData.PlayerDeathEventData{}
	err := env.UnmarshalToDaprEventData(input)
	if err != nil {
		serviceLog.Error("PlayerDeathEvent UnmarshalEvent fail err: %v ", err)
		return
	}

	if input.MsgVersion < serviceCnf.GetInstance().StartMs {
		return
	}

	serviceLog.Info("service receive PlayerDeathEvent: %+v", input)
	playerDataModel, _ := playerModel.GetPlayerDataModel()
	pos := &proto.Vector3{X: input.PosX, Y: input.PosY, Z: input.PosZ}
	if err = playerDataModel.OnPlayerDeath(
		input.UserId, pos, input.KillerId,
		proto.EntityType(input.KillerType), input.KillerName,
	); err != nil {
		serviceLog.Error("PlayerDeathEventData OnPlayerDeath err: %v", err)
		return
	}
}

func GRPCUserTaskRewardEvent(env contract.IModuleEventReq, curMs int64) {
	input := &pubsubEventData.UserTaskRewardEvent{}
	err := env.UnmarshalToDaprEventData(input)
	if err != nil {
		serviceLog.Error("UserTaskRewardEvent UnmarshalEvent fail err: %v ", err)
		return
	}

	if input.MsgVersion < serviceCnf.GetInstance().StartMs {
		return
	}

	serviceLog.Info("service receive UserTaskRewardEvent: %+v", input)
	playerDataModel, _ := playerModel.GetPlayerDataModel()

	// call mint task reward NFT is in task service, so reward exp add in here
	if err = playerDataModel.AddExp(input.UserId, input.Exp); err != nil {
		serviceLog.Error("UserTaskRewardEvent  addExp err: %v", err)
		return
	}
}

func GRPCUserChangeServiceEvent(env contract.IModuleEventReq, curMs int64) {
	input := &pubsubEventData.UserChangeServiceEvent{}
	err := env.UnmarshalToDaprEventData(input)
	if err != nil {
		serviceLog.Error("UserChangeServiceEvent UnmarshalEvent fail err: %v ", err)
		return
	}

	if input.MsgVersion < serviceCnf.GetInstance().StartMs {
		return
	}

	serviceLog.Info("Receive UserChangeServiceEvent: %+v", input)

	agentModel := userAgent.GetUserAgentModel()
	agent, exist := agentModel.GetUserAgent(input.UserId)
	if exist {
		agent.TryUpdate(agent.UserId, agent.AgentAppId, agent.SocketId, input.ToService.AppId)
	} else {
		agent, _ = agentModel.AddUserAgentRecord(
			input.UserId,
			input.UserAgentAppId,
			input.UserSocketId,
			input.ToService.AppId,
		)
	}
}
