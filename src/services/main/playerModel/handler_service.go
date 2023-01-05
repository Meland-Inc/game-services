package playerModel

import (
	"fmt"
	"game-message-core/grpc/methodData"
	"game-message-core/grpc/pubsubEventData"
	"game-message-core/proto"

	base_data "game-message-core/grpc/baseData"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/global/contract"
	"github.com/Meland-Inc/game-services/src/global/grpcAPI/grpcInvoke"
	"github.com/Meland-Inc/game-services/src/global/grpcAPI/grpcNetTool"
	"github.com/Meland-Inc/game-services/src/global/module"
	"github.com/Meland-Inc/game-services/src/global/serviceCnf"
	"github.com/Meland-Inc/game-services/src/global/userAgent"
	login_model "github.com/Meland-Inc/game-services/src/services/main/loginModel"
	"github.com/dapr/go-sdk/service/common"
)

func (p *PlayerDataModel) GRPCGetUserDataHandler(env contract.IModuleEventReq, curMs int64) {
	inputBs, ok := env.GetMsg().([]byte)
	serviceLog.Debug("received service getUserData : %s, [%v]", inputBs, ok)
	if !ok {
		serviceLog.Error("service getUserData to string failed: %s", inputBs)
		return
	}

	output := &methodData.GetUserDataOutput{Success: true}
	result := &module.ModuleEventResult{}
	defer func() {
		if output.ErrMsg != "" {
			output.Success = false
		}
		result.SetResult(output)
		env.WriteResult(result)
	}()

	input := &methodData.GetUserDataInput{}
	err := grpcNetTool.UnmarshalGrpcData(inputBs, input)
	if err != nil {
		output.ErrMsg = err.Error()
		return
	}

	baseData, sceneData, avatars, profile, err := p.PlayerAllData(input.UserId)
	if err != nil {
		output.ErrMsg = err.Error()
		return
	}

	pos := &proto.Vector3{X: sceneData.X, Y: sceneData.Y, Z: sceneData.Z}
	dir := &proto.Vector3{X: sceneData.DirX, Y: sceneData.DirY, Z: sceneData.DirZ}
	pbAvatars := []proto.PlayerAvatar{}
	for _, avatar := range avatars {
		pbAvatars = append(pbAvatars, *avatar.ToNetPlayerAvatar())
	}

	output.BaseData.Set(baseData.ToNetPlayerBaseData())
	output.Profile.Set(profile)
	output.MapId = sceneData.MapId
	output.Pos.Set(pos)
	output.Dir.Set(dir)
	for _, avatar := range pbAvatars {
		grpcAvatar := base_data.GrpcPlayerAvatar{
			Position: avatar.Position,
			ObjectId: avatar.ObjectId,
		}
		grpcAvatar.Attribute = &base_data.GrpcAvatarAttribute{}
		grpcAvatar.Attribute.Set(avatar.Attribute)
		output.Avatars = append(output.Avatars, grpcAvatar)
	}
}

func (p *PlayerDataModel) GRPCMintUserNftHandler(env contract.IModuleEventReq, curMs int64) {
	inputBs, ok := env.GetMsg().([]byte)
	serviceLog.Debug("received service MintUserNft : %s, [%v]", inputBs, ok)
	if !ok {
		serviceLog.Error("service MintUserNft to string failed: %s", inputBs)
		return
	}

	output := &methodData.MainServiceActionMintNftOutput{Success: true}
	result := &module.ModuleEventResult{}
	defer func() {
		if output.FailedMsg != "" {
			output.Success = false
		}
		result.SetResult(output)
		env.WriteResult(result)
	}()

	input := &methodData.MainServiceActionMintNftInput{}
	err := grpcNetTool.UnmarshalGrpcData(inputBs, input)
	if err != nil {
		output.FailedMsg = err.Error()
		return
	}
	if input.UserId < 1 {
		output.FailedMsg = fmt.Sprintf("invalid user id: %d", input.UserId)
		return
	}

	err = grpcInvoke.Web3MintNFT(input.UserId, input.Item.Cid, input.Item.Num, input.Item.Quality, 0, 0)
	if err != nil {
		output.FailedMsg = err.Error()
	}
}

func (p *PlayerDataModel) GRPCTakeUserNftHandler(env contract.IModuleEventReq, curMs int64) {
	inputBs, ok := env.GetMsg().([]byte)
	serviceLog.Debug("received service TakeUserNft : %s, [%v]", inputBs, ok)
	if !ok {
		serviceLog.Error("service TakeUserNft to string failed: %s", inputBs)
		return
	}

	output := &methodData.MainServiceActionTakeNftOutput{Success: true}
	result := &module.ModuleEventResult{}
	defer func() {
		if output.FailedMsg != "" {
			output.Success = false
		}
		result.SetResult(output)
		env.WriteResult(result)
	}()

	input := &methodData.MainServiceActionTakeNftInput{}
	err := grpcNetTool.UnmarshalGrpcData(inputBs, input)
	if err != nil {
		output.FailedMsg = err.Error()
		return
	}
	if input.UserId < 1 {
		output.FailedMsg = fmt.Sprintf("invalid user id: %d", input.UserId)
		return
	}

	playerItem, err := p.GetPlayerItems(input.UserId)
	if err != nil {
		output.FailedMsg = err.Error()
		return
	}

	for _, tn := range input.TakeNfts {
		var giveCount = tn.Num
		for _, item := range playerItem.Items {
			if tn.NftId != "" && tn.NftId != item.Id {
				continue
			}
			if tn.ItemCid != 0 && tn.ItemCid != item.Cid {
				continue
			}
			giveCount -= item.Num
			if giveCount <= 0 {
				break
			}
		}
		if giveCount > 0 {
			output.FailedMsg = fmt.Sprintf("not fund NFT %+v", tn)
			return
		}
	}

	for _, takeNft := range input.TakeNfts {
		if takeNft.NftId != "" {
			err = p.TakeNftById(input.UserId, takeNft.NftId, takeNft.Num)
		} else {
			err = p.TakeNftByItemCid(input.UserId, takeNft.ItemCid, takeNft.Num)
		}
		if err != nil {
			serviceLog.Error(err.Error())
		}
	}
}

// -------------------- pubsub event -----------------------

func (p *PlayerDataModel) GRPCUserChangeServiceEvent(env contract.IModuleEventReq, curMs int64) {
	msg, ok := env.GetMsg().(*common.TopicEvent)
	serviceLog.Info("UserChangeServiceEvent : %s, [%v]", msg, ok)
	if !ok {
		serviceLog.Error("UserChangeServiceEvent to TopicEvent failed: %v", msg)
		return
	}

	input := &pubsubEventData.UserChangeServiceEvent{}
	err := grpcNetTool.UnmarshalGrpcTopicEvent(msg, input)
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

func (p *PlayerDataModel) GRPCUserEnterGameEvent(env contract.IModuleEventReq, curMs int64) {
	msg, ok := env.GetMsg().(*common.TopicEvent)
	serviceLog.Info("UserEnterGameEvent : %s, [%v]", msg, ok)
	if !ok {
		serviceLog.Error("UserEnterGameEvent to TopicEvent failed: %v", msg)
		return
	}

	input := &pubsubEventData.UserEnterGameEvent{}
	err := grpcNetTool.UnmarshalGrpcTopicEvent(msg, input)
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

func (p *PlayerDataModel) GRPCUserLeaveGameEvent(env contract.IModuleEventReq, curMs int64) {
	msg, ok := env.GetMsg().(*common.TopicEvent)
	serviceLog.Info("UserLeaveGameEvent : %s, [%v]", msg, ok)
	if !ok {
		serviceLog.Error("UserLeaveGameEvent to TopicEvent failed: %v", msg)
		return
	}

	input := &pubsubEventData.UserLeaveGameEvent{}
	err := grpcNetTool.UnmarshalGrpcTopicEvent(msg, input)
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

func (p *PlayerDataModel) GRPCSavePlayerDataEvent(env contract.IModuleEventReq, curMs int64) {
	msg, ok := env.GetMsg().(*common.TopicEvent)
	serviceLog.Info("SavePlayerDataEvent : %s, [%v]", msg, ok)
	if !ok {
		serviceLog.Error("SavePlayerDataEvent to TopicEvent failed: %v", msg)
		return
	}

	input := &pubsubEventData.SavePlayerEventData{}
	err := grpcNetTool.UnmarshalGrpcTopicEvent(msg, input)
	if err != nil {
		serviceLog.Error("SavePlayerDataEvent UnmarshalEvent fail err: %v ", err)
		return
	}

	serviceLog.Info("receive savePlayerEvent: %+v", input)

	sceneData, err := p.GetPlayerSceneData(input.UserId)
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
		err = p.UpPlayerSceneData(sceneData)
	case proto.SceneServiceSubType_Dungeon, proto.SceneServiceSubType_Home:
		sceneData.Hp = input.CurHP
		err = p.UpPlayerSceneData(sceneData)
	default:
		err = fmt.Errorf("invalid service sub type %v", input.FormService.SceneSerSubType)
	}
	if err != nil {
		serviceLog.Error(err.Error())
	}
}

func (p *PlayerDataModel) GRPCKillMonsterEvent(env contract.IModuleEventReq, curMs int64) {
	msg, ok := env.GetMsg().(*common.TopicEvent)
	serviceLog.Info("KillMonsterEvent : %s, [%v]", msg, ok)
	if !ok {
		serviceLog.Error("KillMonsterEvent to TopicEvent failed: %v", msg)
		return
	}

	input := &pubsubEventData.KillMonsterEventData{}
	err := grpcNetTool.UnmarshalGrpcTopicEvent(msg, input)
	if err != nil {
		serviceLog.Error("KillMonsterEvent UnmarshalEvent fail err: %v ", err)
		return
	}

	if input.MsgVersion < serviceCnf.GetInstance().StartMs {
		return
	}

	serviceLog.Info("service receive KillMonsterEvent: %+v", input)

	err = p.AddExp(input.UserId, input.Exp)
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

func (p *PlayerDataModel) GRPCPlayerDeathEvent(env contract.IModuleEventReq, curMs int64) {
	msg, ok := env.GetMsg().(*common.TopicEvent)
	serviceLog.Info("PlayerDeathEvent : %s, [%v]", msg, ok)
	if !ok {
		serviceLog.Error("PlayerDeathEvent to TopicEvent failed: %v", msg)
		return
	}

	input := &pubsubEventData.PlayerDeathEventData{}
	err := grpcNetTool.UnmarshalGrpcTopicEvent(msg, input)
	if err != nil {
		serviceLog.Error("PlayerDeathEvent UnmarshalEvent fail err: %v ", err)
		return
	}

	if input.MsgVersion < serviceCnf.GetInstance().StartMs {
		return
	}

	serviceLog.Info("service receive PlayerDeathEvent: %+v", input)

	pos := &proto.Vector3{X: input.PosX, Y: input.PosY, Z: input.PosZ}
	if err = p.OnPlayerDeath(
		input.UserId, pos, input.KillerId,
		proto.EntityType(input.KillerType), input.KillerName,
	); err != nil {
		serviceLog.Error("PlayerDeathEventData OnPlayerDeath err: %v", err)
		return
	}
}

func (p *PlayerDataModel) GRPCUserTaskRewardEvent(env contract.IModuleEventReq, curMs int64) {
	msg, ok := env.GetMsg().(*common.TopicEvent)
	serviceLog.Info("UserTaskRewardEvent : %s, [%v]", msg, ok)
	if !ok {
		serviceLog.Error("UserTaskRewardEvent to TopicEvent failed: %v", msg)
		return
	}

	input := &pubsubEventData.UserTaskRewardEvent{}
	err := grpcNetTool.UnmarshalGrpcTopicEvent(msg, input)
	if err != nil {
		serviceLog.Error("UserTaskRewardEvent UnmarshalEvent fail err: %v ", err)
		return
	}

	if input.MsgVersion < serviceCnf.GetInstance().StartMs {
		return
	}

	serviceLog.Info("service receive UserTaskRewardEvent: %+v", input)

	// call mint task reward NFT is in task service, so reward exp add in here
	if err = p.AddExp(input.UserId, input.Exp); err != nil {
		serviceLog.Error("UserTaskRewardEvent  addExp err: %v", err)
		return
	}
}
