package playerModel

import (
	"game-message-core/grpc/methodData"
	"game-message-core/proto"
	"game-message-core/protoTool"

	"github.com/Meland-Inc/game-services/src/common/net/msgParser"
	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/common/time_helper"
	"github.com/Meland-Inc/game-services/src/global/auth"
	"github.com/Meland-Inc/game-services/src/global/contract"
	"github.com/Meland-Inc/game-services/src/global/grpcAPI/grpcNetTool"
	"github.com/Meland-Inc/game-services/src/global/userAgent"
	login_model "github.com/Meland-Inc/game-services/src/services/main/loginModel"
)

func (p *PlayerDataModel) clientMsgHandler(env contract.IModuleEventReq, curMs int64) {
	bs, ok := env.GetMsg().([]byte)
	serviceLog.Info("client msg: %s, [%v]", bs, ok)
	if !ok {
		serviceLog.Error("client msg to string failed: %v", bs)
		return
	}

	serviceLog.Info("main service received clientPbMsg data: %v", string(bs))

	input := &methodData.PullClientMessageInput{}
	err := grpcNetTool.UnmarshalGrpcData(bs, input)
	if err != nil {
		serviceLog.Error("client msg input Unmarshal error: %v", err)
		return
	}

	agent := userAgent.GetOrStoreUserAgent(input)
	msg, err := protoTool.UnMarshalToEnvelope(input.MsgBody)
	if err != nil {
		serviceLog.Error("Unmarshal Envelope fail err: %+v", err)
		return
	}

	switch proto.EnvelopeType(input.MsgId) {
	case proto.EnvelopeType_SigninPlayer:
		p.SingInHandler(agent, input, msg)
	case proto.EnvelopeType_ItemGet:
		p.ItemGetHandler(agent, input, msg)
	case proto.EnvelopeType_ItemUse:
		p.ItemUseHandler(agent, input, msg)
	case proto.EnvelopeType_UpdateAvatar:
		p.LoadAvatarHandler(agent, input, msg)
	case proto.EnvelopeType_UnloadAvatar:
		p.UnloadAvatarHandler(agent, input, msg)
	case proto.EnvelopeType_GetItemSlot:
		p.ItemSlotGetHandler(agent, input, msg)
	case proto.EnvelopeType_UpgradeItemSlot:
		p.ItemSlotUpgradeHandler(agent, input, msg)
	case proto.EnvelopeType_UpgradePlayerLevel:
		p.PlayerLevelUpgradHandler(agent, input, msg)
	}
}

func (p *PlayerDataModel) SingInHandler(
	agent *userAgent.UserAgentData, input *methodData.PullClientMessageInput, msg *proto.Envelope,
) {
	res := &proto.SigninPlayerResponse{}
	respMsg := userAgent.MakeResponseMsg(msg)
	defer func() {
		if respMsg.ErrorMessage != "" {
			respMsg.ErrorCode = 20001 // TODO: USE PROTO ERROR CODE
			serviceLog.Error("main service SingIn Player err: %s", respMsg.ErrorMessage)
		}
		respMsg.Payload = &proto.Envelope_SigninPlayerResponse{SigninPlayerResponse: res}
		userAgent.ResponseClientMessage(agent, input, respMsg)
	}()

	req := msg.GetSigninPlayerRequest()
	if req == nil {
		respMsg.ErrorMessage = "singIn player request is nil"
		serviceLog.Error(respMsg.ErrorMessage)
		return
	}

	userId, err := auth.GetUserIdByToken(req.Token)
	if err != nil {
		respMsg.ErrorMessage = err.Error()
		return
	}

	loginModel, _ := login_model.GetLoginModel()
	sceneAppId, err := loginModel.GetUserLoginData(userId, input.AgentAppId, input.SocketId)
	if err != nil {
		respMsg.ErrorMessage = err.Error()
		return
	}
	if req.SceneServiceAppId != "" {
		sceneAppId = req.SceneServiceAppId
	}

	playerData, err := p.PlayerProtoData(userId)
	if err != nil {
		respMsg.ErrorMessage = err.Error()
		return
	}

	// 登录时 input.UserId == 0 所以此处需要重新init userAgent
	input.UserId = userId
	agent = userAgent.GetOrStoreUserAgent(input)
	agent.InMapId = playerData.MapId

	res.SceneServiceAppId = sceneAppId
	res.ClientTime = req.ClientTime
	res.ServerTime = time_helper.NowUTCMill()
	res.Player = playerData
}

func (p *PlayerDataModel) ItemGetHandler(
	agent *userAgent.UserAgentData, input *methodData.PullClientMessageInput, msg *proto.Envelope,
) {
	respMsg := userAgent.MakeResponseMsg(msg)
	defer func() {
		if respMsg.ErrorMessage != "" {
			respMsg.ErrorCode = 20002 // TODO: USE PROTO ERROR CODE
			serviceLog.Error(respMsg.ErrorMessage)
		}
		respMsg.Payload = &proto.Envelope_ItemGetResponse{
			ItemGetResponse: &proto.ItemGetResponse{},
		}
		userAgent.ResponseClientMessage(agent, input, respMsg)
	}()

	playerItems, err := p.GetPlayerItems(input.UserId)
	if err != nil {
		respMsg.ErrorMessage = err.Error()
		return
	}

	serviceLog.Info("main service userId[%v] itemLength[%v]", input.UserId, len(playerItems.Items))

	initRes := &proto.BroadCastInitItemResponse{Items: []*proto.Item{}}
	dbMsg := &proto.Envelope{
		Type: proto.EnvelopeType_BroadCastInitItem,
		Payload: &proto.Envelope_BroadCastInitItemResponse{
			BroadCastInitItemResponse: initRes,
		},
	}

	maxIdx := len(playerItems.Items) - 1
	for idx, it := range playerItems.Items {
		initRes.Items = append(initRes.Items, it.ToNetItem())
		var msgDataLength int
		if len(initRes.Items) >= 10 {
			msgBody, _ := protoTool.MarshalProto(dbMsg)
			msgDataLength = len(msgBody)
		}
		if idx >= maxIdx || msgDataLength >= msgParser.MSG_LIMIT-1000 {
			userAgent.ResponseClientMessage(agent, input, dbMsg)
			initRes.Items = []*proto.Item{}
		}
	}
}

func (p *PlayerDataModel) ItemUseHandler(
	agent *userAgent.UserAgentData, input *methodData.PullClientMessageInput, msg *proto.Envelope,
) {
	res := &proto.ItemUseResponse{}
	respMsg := userAgent.MakeResponseMsg(msg)
	defer func() {
		if respMsg.ErrorMessage != "" {
			respMsg.ErrorCode = 20004 // TODO: USE PROTO ERROR CODE
			serviceLog.Error(respMsg.ErrorMessage)
		}
		respMsg.Payload = &proto.Envelope_ItemUseResponse{ItemUseResponse: res}
		userAgent.ResponseClientMessage(agent, input, respMsg)
	}()

	if input.UserId < 1 {
		respMsg.ErrorMessage = "item use Invalid User ID"
		return
	}
	req := msg.GetItemUseRequest()
	if req == nil {
		respMsg.ErrorMessage = "main service use item request is nil"
		return
	}

	err := p.UseItem(input.UserId, req.ItemId, req.Args)
	if err != nil {
		respMsg.ErrorMessage = err.Error()
		return
	}
}

func (p *PlayerDataModel) LoadAvatarHandler(
	agent *userAgent.UserAgentData, input *methodData.PullClientMessageInput, msg *proto.Envelope,
) {
	res := &proto.UpdateAvatarResponse{}
	respMsg := userAgent.MakeResponseMsg(msg)
	defer func() {
		if respMsg.ErrorMessage != "" {
			respMsg.ErrorCode = 20004 // TODO: USE PROTO ERROR CODE
			serviceLog.Error(respMsg.ErrorMessage)
		}
		respMsg.Payload = &proto.Envelope_UpdateAvatarResponse{UpdateAvatarResponse: res}
		userAgent.ResponseClientMessage(agent, input, respMsg)
	}()

	if input.UserId < 1 {
		respMsg.ErrorMessage = "load avatar Invalid User ID"
		return
	}
	req := msg.GetUpdateAvatarRequest()
	if req == nil {
		serviceLog.Error("main service load avatar request is nil")
		return
	}
	err := p.LoadAvatar(input.UserId, req.ItemId, req.IsAppearance)
	if err != nil {
		respMsg.ErrorMessage = err.Error()
		return
	}
}

func (p *PlayerDataModel) UnloadAvatarHandler(
	agent *userAgent.UserAgentData, input *methodData.PullClientMessageInput, msg *proto.Envelope,
) {
	res := &proto.UnloadAvatarResponse{}
	respMsg := userAgent.MakeResponseMsg(msg)
	defer func() {
		if respMsg.ErrorMessage != "" {
			respMsg.ErrorCode = 20005 // TODO: USE PROTO ERROR CODE
		}
		respMsg.Payload = &proto.Envelope_UnloadAvatarResponse{UnloadAvatarResponse: res}
		userAgent.ResponseClientMessage(agent, input, respMsg)
	}()

	if input.UserId < 1 {
		respMsg.ErrorMessage = "load avatar Invalid User ID"
		return
	}
	req := msg.GetUpdateAvatarRequest()
	if req == nil {
		serviceLog.Error("main service load avatar request is nil")
		return
	}
	err := p.UnloadAvatar(input.UserId, req.ItemId, true, false)
	if err != nil {
		respMsg.ErrorMessage = err.Error()
		return
	}
}

func (p *PlayerDataModel) ItemSlotGetHandler(
	agent *userAgent.UserAgentData, input *methodData.PullClientMessageInput, msg *proto.Envelope,
) {
	res := &proto.GetItemSlotResponse{}
	respMsg := userAgent.MakeResponseMsg(msg)
	defer func() {
		if respMsg.ErrorMessage != "" {
			respMsg.ErrorCode = 20006 // TODO: USE PROTO ERROR CODE
		}
		respMsg.Payload = &proto.Envelope_GetItemSlotResponse{GetItemSlotResponse: res}
		userAgent.ResponseClientMessage(agent, input, respMsg)
	}()

	if input.UserId < 1 {
		respMsg.ErrorMessage = "item slot get Invalid User ID"
		return
	}
	playerSlot, err := p.GetPlayerItemSlots(input.UserId)
	if err != nil {
		respMsg.ErrorMessage = err.Error()
		return
	}
	for _, s := range playerSlot.GetSlotList().SlotList {
		res.Slots = append(res.Slots, &proto.ItemSlot{
			Level:    int32(s.Level),
			Position: proto.AvatarPosition(s.Position),
		})
	}
}

func (p *PlayerDataModel) PlayerLevelUpgradHandler(
	agent *userAgent.UserAgentData, input *methodData.PullClientMessageInput, msg *proto.Envelope,
) {
	res := &proto.UpgradeItemSlotResponse{}
	respMsg := userAgent.MakeResponseMsg(msg)
	defer func() {
		if respMsg.ErrorMessage != "" {
			respMsg.ErrorCode = 20007 // TODO: USE PROTO ERROR CODE
		}
		respMsg.Payload = &proto.Envelope_UpgradeItemSlotResponse{UpgradeItemSlotResponse: res}
		userAgent.ResponseClientMessage(agent, input, respMsg)
	}()

	if input.UserId < 1 {
		respMsg.ErrorMessage = "item slot upgrade Invalid User ID"
		return
	}
	req := msg.GetUpgradeItemSlotRequest()
	if req == nil {
		serviceLog.Error("main service upgrade slot request is nil")
		return
	}
	slotData, err := p.UpgradeItemSlots(input.UserId, req.Position, false)
	if err != nil {
		respMsg.ErrorMessage = err.Error()
		return
	}
	for _, s := range slotData.GetSlotList().SlotList {
		res.Slots = append(res.Slots, &proto.ItemSlot{
			Level:    int32(s.Level),
			Position: proto.AvatarPosition(s.Position),
		})
	}
}

func (p *PlayerDataModel) ItemSlotUpgradeHandler(
	agent *userAgent.UserAgentData, input *methodData.PullClientMessageInput, msg *proto.Envelope,
) {
	res := &proto.UpgradePlayerLevelResponse{}
	respMsg := userAgent.MakeResponseMsg(msg)
	defer func() {
		if respMsg.ErrorMessage != "" {
			respMsg.ErrorCode = 20006 // TODO: USE PROTO ERROR CODE
		}
		respMsg.Payload = &proto.Envelope_UpgradePlayerLevelResponse{UpgradePlayerLevelResponse: res}
		userAgent.ResponseClientMessage(agent, input, respMsg)
	}()

	if input.UserId < 1 {
		respMsg.ErrorMessage = "upgrade player level Invalid User ID"
		return
	}

	lv, exp, err := p.UpgradePlayerLevel(input.UserId)
	if err != nil {
		respMsg.ErrorMessage = err.Error()
		return
	}
	res.CurExp = int64(exp)
	res.CurLevel = lv
}
