package clientMsgHandle

import (
	"fmt"
	"game-message-core/grpc/methodData"
	"game-message-core/proto"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/common/time_helper"
	"github.com/Meland-Inc/game-services/src/global/grpcAPI/grpcInvoke"
	"github.com/Meland-Inc/game-services/src/global/serviceCnf"
	"github.com/Meland-Inc/game-services/src/services/main/playerModel"
	"github.com/spf13/cast"
)

func SingInHandler(input *methodData.PullClientMessageInput, msg *proto.Envelope) {
	res := &proto.SigninPlayerResponse{}
	respMsg := makeResponseMsg(msg)
	defer func() {
		if respMsg.ErrorMessage != "" {
			respMsg.ErrorCode = 20001 // TODO: USE PROTO ERROR CODE
			serviceLog.Error("main service SingIn Player err: %s", respMsg.ErrorMessage)
		}
		respMsg.Payload = &proto.Envelope_SigninPlayerResponse{SigninPlayerResponse: res}
		ResponseClientMessage(input, respMsg)
	}()

	req := msg.GetSigninPlayerRequest()
	if req == nil {
		serviceLog.Error("main service singIn player request is nil")
		return
	}
	// check token
	// userIdStr, err := auth.CheckDefaultAuth(req.Token)
	// if err != nil {
	// 	respMsg.ErrorMessage = err.Error()
	// 	return
	// }

	userIdStr := "680"

	input.UserId = cast.ToInt64(userIdStr)
	dataModel, err := playerModel.GetPlayerDataModel()
	if err != nil {
		respMsg.ErrorMessage = err.Error()
		return
	}

	baseData, sceneData, avatars, profile, err := dataModel.PlayerAllData(input.UserId)
	if err != nil {
		respMsg.ErrorMessage = err.Error()
		return
	}

	pos := &proto.Vector3{X: sceneData.X, Y: sceneData.Y, Z: sceneData.Z}
	dir := &proto.Vector3{X: sceneData.DirX, Y: sceneData.DirY, Z: sceneData.DirZ}
	res.ClientTime = req.ClientTime
	res.ServerTime = time_helper.NowUTCMill()
	res.LastLoginTime = sceneData.LastLoginAt.UnixMilli()
	res.Player = &proto.Player{
		BaseData: baseData.ToNetPlayerBaseData(),
		Profile:  profile,
		Active:   sceneData.Hp > 0,
		MapId:    sceneData.MapId,
		Position: pos,
		Dir:      dir,
	}
	for _, avatar := range avatars {
		res.Player.Avatars = append(res.Player.Avatars, avatar.ToNetPlayerAvatar())
	}

	sceneAppId, err := getSceneAppId(req.SceneServiceAppId, sceneData.MapId)
	if err != nil {
		respMsg.ErrorMessage = err.Error()
		return
	}
	res.SceneServiceAppId = sceneAppId
}

func getSceneAppId(clientPushSceneAppId string, mapId int32) (string, error) {
	if serviceCnf.GetInstance().IsDevelop && clientPushSceneAppId != "" {
		return clientPushSceneAppId, nil
	}

	serviceOut, err := grpcInvoke.RPCSelectService(proto.ServiceType_ServiceTypeScene, mapId)
	if err != nil {
		return "", err
	}
	if serviceOut.ErrorCode > 0 {
		return "", fmt.Errorf(serviceOut.ErrorMessage)
	}
	return serviceOut.ServiceAppId, nil
}
