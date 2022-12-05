package playerModel

import (
	"fmt"
	"game-message-core/grpc/methodData"
	"game-message-core/grpc/pubsubEventData"
	"game-message-core/proto"

	base_data "game-message-core/grpc/baseData"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/global/component"
	"github.com/Meland-Inc/game-services/src/global/grpcAPI/grpcInvoke"
	"github.com/Meland-Inc/game-services/src/global/grpcAPI/grpcNetTool"
	"github.com/Meland-Inc/game-services/src/global/serviceCnf"
	"github.com/Meland-Inc/game-services/src/global/userAgent"
	login_model "github.com/Meland-Inc/game-services/src/services/main/loginModel"
	"github.com/dapr/go-sdk/service/common"
)

func (p *PlayerDataModel) GRPCGetUserDataHandler(env *component.ModelEventReq, curMs int64) {
	inputBs, ok := env.Msg.([]byte)
	serviceLog.Debug("received service getUserData : %s, [%v]", inputBs, ok)
	if !ok {
		serviceLog.Error("service getUserData to string failed: %s", inputBs)
		return
	}

	output := &methodData.GetUserDataOutput{}
	result := &component.ModelEventResult{}
	defer func() {
		result.SetResult(output)
		env.WriteResult(result)
	}()

	input := &methodData.GetUserDataInput{}
	err := grpcNetTool.UnmarshalGrpcData(inputBs, input)
	if err != nil {
		result.Err = err
		return
	}

	baseData, sceneData, avatars, profile, err := p.PlayerAllData(input.UserId)
	if err != nil {
		result.Err = err
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

func (p *PlayerDataModel) GRPCTakeUserNftHandler(env *component.ModelEventReq, curMs int64) {
	inputBs, ok := env.Msg.([]byte)
	serviceLog.Debug("received service TakeUserNft : %s, [%v]", inputBs, ok)
	if !ok {
		serviceLog.Error("service TakeUserNft to string failed: %s", inputBs)
		return
	}

	output := &methodData.MainServiceActionTakeNftOutput{}
	result := &component.ModelEventResult{}
	defer func() {
		result.SetResult(output)
		env.WriteResult(result)
	}()

	input := &methodData.MainServiceActionTakeNftInput{}
	err := grpcNetTool.UnmarshalGrpcData(inputBs, input)
	if err != nil {
		result.Err = err
		return
	}
	if input.UserId < 1 {
		result.Err = fmt.Errorf("invalid user id: %d", input.UserId)
		return
	}

	playerItem, err := p.GetPlayerItems(input.UserId)
	if err != nil {
		result.Err = err
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
			result.Err = fmt.Errorf("not fund NFT %+v", tn)
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

func (p *PlayerDataModel) GRPCUserEnterGameEvent(env *component.ModelEventReq, curMs int64) {
	msg, ok := env.Msg.(*common.TopicEvent)
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

func (p *PlayerDataModel) GRPCUserLeaveGameEvent(env *component.ModelEventReq, curMs int64) {
	msg, ok := env.Msg.(*common.TopicEvent)
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

func (p *PlayerDataModel) GRPCSavePlayerDataEvent(env *component.ModelEventReq, curMs int64) {
	msg, ok := env.Msg.(*common.TopicEvent)
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
	if err = p.UpPlayerSceneData(
		input.UserId, input.CurHP, sceneData.Level, sceneData.Exp,
		input.MapId, input.PosX, input.PosY, input.PosZ,
		input.DirX, input.DirY, input.DirZ,
	); err != nil {
		serviceLog.Error(err.Error())
	}
}

func (p *PlayerDataModel) GRPCKillMonsterEvent(env *component.ModelEventReq, curMs int64) {
	msg, ok := env.Msg.(*common.TopicEvent)
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
		if err := grpcInvoke.MintNFT(
			input.UserId, drop.Cid, drop.Num, drop.Quality, int32(input.PosX), int32(input.PosZ),
		); err != nil {
			serviceLog.Error("mint nft[%d] failed: %v", drop.Cid, err)
		}
	}
}

func (p *PlayerDataModel) GRPCPlayerDeathEvent(env *component.ModelEventReq, curMs int64) {
	msg, ok := env.Msg.(*common.TopicEvent)
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

func (p *PlayerDataModel) GRPCUserTaskRewardEvent(env *component.ModelEventReq, curMs int64) {
	msg, ok := env.Msg.(*common.TopicEvent)
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
