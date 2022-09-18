package clientMsgHandle

import (
	"game-message-core/grpc/methodData"
	"game-message-core/proto"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/common/time_helper"
	"github.com/Meland-Inc/game-services/src/global/auth"
	"github.com/Meland-Inc/game-services/src/global/component"
	"github.com/Meland-Inc/game-services/src/global/serviceCnf"
	"github.com/Meland-Inc/game-services/src/global/userAgent"
	"github.com/Meland-Inc/game-services/src/services/main/playerModel"
	"github.com/spf13/cast"
)

func ResponseClientMessage(userId int64, respMsg *proto.Envelope) {
	iUserAgentModel, exist := component.GetInstance().GetModel(component.MODEL_NAME_USER_AGENT)
	if !exist {
		return
	}
	agentModel := iUserAgentModel.(*userAgent.UserAgentModel)
	agent, exist := agentModel.GetUserAgent(userId)
	if !exist {
		serviceLog.Warning("user [%d] agent data not found", userId)
		return
	}
	agent.SendToPlayer(serviceCnf.GetInstance().ServerName, respMsg)

	if respMsg.ErrorMessage != "" {
		serviceLog.Error(
			"responseClient [%v] Msg err : [%d][%s]",
			respMsg.Type, respMsg.ErrorCode, respMsg.ErrorMessage,
		)
	}
}

func makeResponseMsg(msg *proto.Envelope) *proto.Envelope {
	return &proto.Envelope{
		Type:  msg.Type,
		SeqId: msg.SeqId,
	}
}

func SingInHandle(input *methodData.PullClientMessageInput, msg *proto.Envelope) {
	res := &proto.SigninPlayerResponse{}
	respMsg := makeResponseMsg(msg)
	defer func() {
		if respMsg.ErrorMessage != "" {
			respMsg.ErrorCode = 20001 // TODO: USE PROTO ERROR CODE
		}
		respMsg.Payload = &proto.Envelope_SigninPlayerResponse{SigninPlayerResponse: res}
		ResponseClientMessage(input.UserId, respMsg)
	}()

	req := msg.GetSigninPlayerRequest()
	if req == nil {
		serviceLog.Error("main service singIn player request is nil")
		return
	}
	userIdStr, err := auth.CheckDefaultAuth(req.Token)
	if err != nil {
		respMsg.ErrorMessage = err.Error()
		return
	}

	input.UserId = cast.ToInt64(userIdStr)
	iPlayerModel, exist := component.GetInstance().GetModel(component.MODEL_NAME_PLAYER_DATA)
	if !exist {
		respMsg.ErrorMessage = "player data select failed"
		return
	}

	dataModel, _ := iPlayerModel.(*playerModel.PlayerDataModel)
	baseData, sceneData, avatars, profile, err := dataModel.PlayerAllData(input.UserId)
	if err != nil {
		respMsg.ErrorMessage = err.Error()
		return
	}

	pbAvatars := []*proto.PlayerAvatar{}
	for _, avatar := range avatars {
		pbAvatars = append(pbAvatars, &proto.PlayerAvatar{
			Position:  proto.AvatarPosition(avatar.AvatarPos),
			ObjectId:  avatar.Cid,
			Attribute: avatar.Attribute,
		})
	}

	res.ClientTime = req.ClientTime
	res.ServerTime = time_helper.NowUTCMill()
	res.LastLoginTime = sceneData.LastLoginAt.UnixMilli()
	res.Player = &proto.Player{
		BaseData: baseData.ToNetPlayerBaseData(),
		Avatars:  pbAvatars,
		Profile:  profile,
		Active:   sceneData.Hp > 0,
		MapId:    sceneData.MapId,
		Position: &proto.Vector3{X: float32(sceneData.X), Y: float32(sceneData.Y), Z: float32(sceneData.Z)},
		Dir:      &proto.Vector3{X: float32(sceneData.DirX), Y: float32(sceneData.DirY), Z: float32(sceneData.DirZ)},
	}

	if serviceCnf.GetInstance().IsDevelop && req.SceneServiceAppId != "" {
		res.SceneServiceAppId = req.SceneServiceAppId
	} else {
		// TODO GET SCENE SERVICE appId
	}
}
