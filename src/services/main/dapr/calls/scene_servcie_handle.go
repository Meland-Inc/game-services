package daprCalls

import (
	"context"
	base_data "game-message-core/grpc/baseData"
	"game-message-core/grpc/methodData"
	"game-message-core/proto"

	"github.com/Meland-Inc/game-services/src/common/daprInvoke"
	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/global/grpcAPI/grpcNetTool"
	"github.com/Meland-Inc/game-services/src/services/main/playerModel"
	"github.com/dapr/go-sdk/service/common"
)

func GRPCGetUserDataHandler(ctx context.Context, in *common.InvocationEvent) (*common.Content, error) {
	resFunc := func(
		success bool, err error, baseData *proto.PlayerBaseData, profile *proto.EntityProfile,
		mapId int32, pos *proto.Vector3, dir *proto.Vector3, avatars []proto.PlayerAvatar,
	) (*common.Content, error) {
		out := &methodData.GetUserDataOutput{}
		out.Success = success
		if err != nil {
			out.ErrMsg = err.Error()
			serviceLog.Error("get user data err: %v", err)
		} else {
			out.BaseData.Set(baseData)
			out.Profile.Set(profile)
			out.MapId = mapId
			out.Pos.Set(pos)
			out.Dir.Set(dir)
			for _, avatar := range avatars {
				grpcAvatar := base_data.GrpcPlayerAvatar{
					Position: avatar.Position,
					ObjectId: avatar.ObjectId,
				}
				grpcAvatar.Attribute = &base_data.GrpcAvatarAttribute{}
				grpcAvatar.Attribute.Set(avatar.Attribute)
				out.Avatars = append(out.Avatars, grpcAvatar)
			}
		}
		content, _ := daprInvoke.MakeOutputContent(in, out)
		return content, err
	}

	serviceLog.Info("get user data received data: %v", string(in.Data))

	input := &methodData.GetUserDataInput{}
	err := grpcNetTool.UnmarshalGrpcData(in.Data, input)
	if err != nil {
		return nil, err
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
	pbAvatars := []proto.PlayerAvatar{}
	for _, avatar := range avatars {
		pbAvatars = append(pbAvatars, *avatar.ToNetPlayerAvatar())
	}

	return resFunc(
		true, err,
		baseData.ToNetPlayerBaseData(),
		profile,
		sceneData.MapId,
		pos,
		dir,
		pbAvatars,
	)
}
