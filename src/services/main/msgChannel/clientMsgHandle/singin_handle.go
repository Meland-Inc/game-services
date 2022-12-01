package clientMsgHandle

import (
	"encoding/json"
	"fmt"
	"game-message-core/grpc"
	"game-message-core/grpc/methodData"
	"game-message-core/proto"
	"game-message-core/protoTool"

	"github.com/Meland-Inc/game-services/src/common/daprInvoke"
	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/common/time_helper"
	"github.com/Meland-Inc/game-services/src/global/auth"
	"github.com/Meland-Inc/game-services/src/global/grpcAPI/grpcInvoke"
	"github.com/Meland-Inc/game-services/src/global/serviceCnf"
	login_model "github.com/Meland-Inc/game-services/src/services/main/loginModel"
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
		responseSingInMessage(input.AgentAppId, input.SocketId, respMsg)
	}()

	req := msg.GetSigninPlayerRequest()
	if req == nil {
		serviceLog.Error("main service singIn player request is nil")
		return
	}

	userId, err := checkUserToken(req.Token)
	if err != nil {
		respMsg.ErrorMessage = err.Error()
		return
	}

	playerData, err := getPlayerData(userId)
	if err != nil {
		respMsg.ErrorMessage = err.Error()
		return
	}

	sceneAppId, err := getSceneAppId(req.SceneServiceAppId, playerData.MapId, userId)
	if err != nil {
		respMsg.ErrorMessage = err.Error()
		return
	}

	err = checkRepeatedSingIn(userId, input.AgentAppId, input.SocketId, sceneAppId)
	if err != nil {
		respMsg.ErrorMessage = err.Error()
		return
	}

	input.UserId = userId
	agent := GetOrStoreUserAgent(input)
	agent.InMapId = playerData.MapId

	res.SceneServiceAppId = sceneAppId
	res.ClientTime = req.ClientTime
	res.ServerTime = time_helper.NowUTCMill()
	res.Player = playerData
}

func checkUserToken(token string) (userId int64, err error) {
	userIdStr, err := auth.CheckDefaultAuth(token)
	if err != nil {
		return 0, err
	}
	return cast.ToInt64(userIdStr), nil
}

func responseSingInMessage(agentAppId, UserSocketId string, msg *proto.Envelope) error {
	msgBody, err := protoTool.MarshalProto(msg)
	if err != nil {
		return err
	}
	input := methodData.BroadCastToClientInput{
		MsgVersion:   time_helper.NowUTCMill(),
		ServiceAppId: serviceCnf.GetInstance().AppId,
		SocketId:     UserSocketId,
		MsgId:        int32(msg.Type),
		MsgBody:      msgBody,
	}

	inputBytes, err := json.Marshal(input)
	if err != nil {
		serviceLog.Error("SendToPlayer Marshal BroadCastInput failed err: %+v", err)
		return err
	}
	_, err = daprInvoke.InvokeMethod(
		agentAppId,
		string(grpc.ProtoMessageActionBroadCastToClient),
		inputBytes,
	)
	if err != nil {
		serviceLog.Error("UserAgentData SendToPlayer InvokeMethod  failed err:%+v", err)
	}
	return err
}

func getSceneAppId(clientPushSceneAppId string, mapId int32, userId int64) (string, error) {
	if serviceCnf.GetInstance().IsDevelop && clientPushSceneAppId != "" {
		return clientPushSceneAppId, nil
	}

	serviceOut, err := grpcInvoke.RPCSelectService(
		proto.ServiceType_ServiceTypeScene,
		proto.SceneServiceSubType_World,
		0,
		mapId,
	)
	serviceLog.Debug("getSceneAppId output = %+v", serviceOut)
	if err != nil {
		return "", err
	}
	if serviceOut.ErrorCode > 0 {
		return "", fmt.Errorf(serviceOut.ErrorMessage)
	}
	return serviceOut.Service.AppId, nil
}

func getPlayerData(userId int64) (*proto.Player, error) {
	dataModel, err := playerModel.GetPlayerDataModel()
	if err != nil {
		return nil, err
	}

	baseData, sceneData, avatars, profile, err := dataModel.PlayerAllData(userId)
	if err != nil {
		return nil, err
	}

	pos := &proto.Vector3{X: sceneData.X, Y: sceneData.Y, Z: sceneData.Z}
	dir := &proto.Vector3{X: sceneData.DirX, Y: sceneData.DirY, Z: sceneData.DirZ}
	player := &proto.Player{
		BaseData: baseData.ToNetPlayerBaseData(),
		Profile:  profile,
		Active:   sceneData.Hp > 0,
		MapId:    sceneData.MapId,
		Position: pos,
		Dir:      dir,
	}
	for _, avatar := range avatars {
		player.Avatars = append(player.Avatars, avatar.ToNetPlayerAvatar())
	}

	return player, nil
}

func checkRepeatedSingIn(userId int64, agentAppId, socketId, sceneAppId string) error {
	loginModel, err := login_model.GetLoginModel()
	if err != nil {
		return err
	}

	return loginModel.OnLogin(userId, agentAppId, socketId, sceneAppId)
}
