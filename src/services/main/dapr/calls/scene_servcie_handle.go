package daprCalls

import (
	"context"
	"encoding/json"
	"game-message-core/grpc/methodData"
	"game-message-core/proto"
	"net/url"

	"github.com/Meland-Inc/game-services/src/common/daprInvoke"
	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/services/main/playerModel"
	"github.com/dapr/go-sdk/service/common"
)

func GRPCGetUserDataHandler(ctx context.Context, in *common.InvocationEvent) (*common.Content, error) {
	resFunc := func(
		success bool, err error, baseData proto.PlayerBaseData, profile proto.EntityProfile,
		mapId int32, pos proto.Vector3, dir proto.Vector3, avatars []proto.PlayerAvatar,
	) (*common.Content, error) {
		out := &methodData.GetUserDataOutput{}
		out.Success = success
		if err != nil {
			out.ErrMsg = err.Error()
			serviceLog.Error("get user data err: %v", err)
		} else {
			out.BaseData = baseData
			out.Profile = profile
			out.MapId = mapId
			out.Pos = pos
			out.Dir = dir
			out.Avatars = avatars
		}

		content, _ := daprInvoke.MakeOutputContent(in, out)
		return content, err
	}

	serviceLog.Info("get user data received data: %v", string(in.Data))
	escStr, err := url.QueryUnescape(string(in.Data))
	if err != nil {
		return nil, err
	}

	input := &methodData.GetUserDataInput{}
	err = json.Unmarshal([]byte(escStr), input)
	if err != nil {
		return resFunc(
			false, err, proto.PlayerBaseData{}, proto.EntityProfile{},
			0, proto.Vector3{}, proto.Vector3{}, []proto.PlayerAvatar{},
		)
	}

	dataModel, err := playerModel.GetPlayerDataModel()
	if err != nil {
		return resFunc(
			false, err, proto.PlayerBaseData{}, proto.EntityProfile{},
			0, proto.Vector3{}, proto.Vector3{}, []proto.PlayerAvatar{},
		)
	}

	baseData, sceneData, avatars, profile, err := dataModel.PlayerAllData(input.UserId)
	if err != nil {
		return resFunc(
			false, err, proto.PlayerBaseData{}, proto.EntityProfile{},
			0, proto.Vector3{}, proto.Vector3{}, []proto.PlayerAvatar{},
		)
	}

	pos := proto.Vector3{X: float32(sceneData.X), Y: float32(sceneData.Y), Z: float32(sceneData.Z)}
	dir := proto.Vector3{X: float32(sceneData.DirX), Y: float32(sceneData.DirY), Z: float32(sceneData.DirZ)}
	pbAvatars := []proto.PlayerAvatar{}
	for _, avatar := range avatars {
		pbAvatars = append(pbAvatars, *avatar.ToNetPlayerAvatar())
	}
	return resFunc(
		false, err,
		*baseData.ToNetPlayerBaseData(),
		*profile,
		sceneData.MapId,
		pos,
		dir,
		pbAvatars,
	)
}
