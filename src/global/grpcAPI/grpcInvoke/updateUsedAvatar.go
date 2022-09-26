package grpcInvoke

import (
	"encoding/json"
	"fmt"
	"game-message-core/grpc"
	"game-message-core/grpc/methodData"
	"game-message-core/proto"

	base_data "game-message-core/grpc/baseData"

	"github.com/Meland-Inc/game-services/src/common/daprInvoke"
	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/common/time_helper"
	"github.com/Meland-Inc/game-services/src/global/grpcAPI/grpcNetTool"
	"github.com/Meland-Inc/game-services/src/global/userAgent"
)

// to scene service update user used avatars and profile
func UpdateUsedAvatar(
	userId int64,
	avatars []*proto.PlayerAvatar,
	profile *proto.EntityProfile,
) error {
	beginMs := time_helper.NowMill()
	defer func() {
		serviceLog.Info("UpdateUsedAvatar used time [%04d]Ms", time_helper.NowMill()-beginMs)
	}()

	agentModel := userAgent.GetUserAgentModel()
	userAgent, exist := agentModel.GetUserAgent(userId)
	if !exist || userAgent.InSceneServiceAppId == "" {
		return fmt.Errorf("user in scene service not found")
	}

	input := &methodData.UpdateUsedAvatarInput{
		MsgVersion: time_helper.NowUTCMill(),
		UserId:     userId,
	}
	input.CurProfile = base_data.GrpcEntityProfile{}
	input.CurProfile.Set(profile)
	for _, avatar := range avatars {
		grpcAvatar := base_data.GrpcPlayerAvatar{}
		grpcAvatar.Set(avatar)
		input.UsingAvatars = append(input.UsingAvatars, grpcAvatar)
	}

	serviceLog.Info("UpdateUsedAvatar: %+v", input)
	inputBytes, err := json.Marshal(input)
	if err != nil {
		return err
	}

	outBytes, err := daprInvoke.InvokeMethod(
		userAgent.InSceneServiceAppId,
		string(grpc.UserActionUpdateUsedAvatar),
		inputBytes,
	)

	serviceLog.Info("UpdateUsedAvatar outBytes = %+v", string(outBytes))

	output := &methodData.UpdateUsedAvatarOutput{}
	err = grpcNetTool.UnmarshalGrpcData(outBytes, output)
	if err != nil {
		return err
	}
	if !output.Success {
		return fmt.Errorf("UpdateUsedAvatar fail err: %s", output.Success)
	}
	return nil
}
