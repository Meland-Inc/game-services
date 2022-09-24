package daprCalls

import (
	"context"

	"game-message-core/proto"
	"game-message-core/protoTool"
	"net/url"

	"github.com/Meland-Inc/game-services/src/common/daprInvoke"
	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/services/main/playerModel"
	"github.com/dapr/go-sdk/service/common"
)

func GRPCGetUserDataHandler(ctx context.Context, in *common.InvocationEvent) (*common.Content, error) {
	resFunc := func(
		success bool, err error, baseData *proto.PlayerBaseData,
		profile *proto.EntityProfile, mapId int32, pos *proto.Vector3,
		dir *proto.Vector3, avatars []*proto.PlayerAvatar,
	) (*common.Content, error) {
		out := &proto.GetUserDataOutput{}
		out.Success = success
		if err != nil {
			out.ErrMsg = err.Error()
			serviceLog.Error("get user data err: %v", err)
		} else {
			out.BaseData = baseData
			out.Profile = profile
			out.MapId = mapId
			out.Position = pos
			out.Dir = dir
			out.Avatars = avatars
		}
		content, _ := daprInvoke.MakeProtoOutputContent(in, out)
		return content, err
	}

	input := &proto.GetUserDataInput{}
	err := protoTool.UnmarshalProto(in.Data, input)
	if err != nil {
		escStr, err := url.QueryUnescape(string(in.Data))
		serviceLog.Info("GetUserDataInput data: %v, err: %+v", escStr, err)
		if err != nil {
			return resFunc(false, err, nil, nil, 0, nil, nil, nil)
		}
		err = protoTool.UnmarshalProto([]byte(escStr), input)
		if err != nil {
			serviceLog.Error("Unmarshal to GetUserDataInput data : %+v, err: $+v", string(in.Data), err)
			return resFunc(false, err, nil, nil, 0, nil, nil, nil)
		}
	}

	dataModel, err := playerModel.GetPlayerDataModel()
	if err != nil {
		return resFunc(false, err, nil, nil, 0, nil, nil, nil)
	}

	baseData, sceneData, avatars, profile, err := dataModel.PlayerAllData(input.UserId)
	if err != nil {
		return resFunc(false, err, nil, nil, 0, nil, nil, nil)
	}

	pos := &proto.Vector3{X: sceneData.X, Y: sceneData.Y, Z: sceneData.Z}
	dir := &proto.Vector3{X: sceneData.DirX, Y: sceneData.DirY, Z: sceneData.DirZ}
	pbAvatars := []*proto.PlayerAvatar{}
	for _, avatar := range avatars {
		pbAvatars = append(pbAvatars, avatar.ToNetPlayerAvatar())
	}

	return resFunc(
		true, nil, baseData.ToNetPlayerBaseData(), profile,
		sceneData.MapId, pos, dir, pbAvatars,
	)
}
