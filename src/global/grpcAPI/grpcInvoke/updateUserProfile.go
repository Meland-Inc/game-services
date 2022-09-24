package grpcInvoke

import (
	"encoding/json"
	"fmt"
	"game-message-core/grpc"
	base_data "game-message-core/grpc/baseData"
	"game-message-core/grpc/methodData"
	"game-message-core/proto"

	"github.com/Meland-Inc/game-services/src/common/daprInvoke"
	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/common/time_helper"
	"github.com/Meland-Inc/game-services/src/global/grpcAPI/grpcNetTool"
	"github.com/Meland-Inc/game-services/src/global/userAgent"
)

// to scene service update user profile
func UpdateUsedProfile(userId int64, profile *proto.EntityProfile) error {
	beginMs := time_helper.NowMill()
	defer func() {
		serviceLog.Info("UpdateUsedProfile used time [%04d]Ms", time_helper.NowMill()-beginMs)
	}()

	agentModel := userAgent.GetUserAgentModel()
	userAgent, exist := agentModel.GetUserAgent(userId)
	if !exist || userAgent.InSceneServiceAppId == "" {
		return fmt.Errorf("user in scene service not found")
	}

	input := &methodData.UpdateUserProfileInput{
		MsgVersion: time_helper.NowUTCMill(),
		UserId:     userId,
	}
	input.CurProfile = base_data.GrpcEntityProfile{}
	input.CurProfile.Set(profile)
	serviceLog.Info("UpdateUsedProfile = %+v", input)
	inputBytes, err := json.Marshal(input)
	if err != nil {
		return err
	}

	outBytes, err := daprInvoke.InvokeMethod(
		userAgent.InSceneServiceAppId,
		string(grpc.UserActionUpdateUserProfile),
		inputBytes,
	)

	serviceLog.Info("UpdateUsedProfile outBytes = %+v", string(outBytes))

	output := &methodData.UpdateUserProfileOutput{}
	err = grpcNetTool.UnmarshalGrpcData(outBytes, output)
	if err != nil {
		return err
	}

	if !output.Success {
		return fmt.Errorf("UpdateUsedProfile fail")
	}
	return nil
}
