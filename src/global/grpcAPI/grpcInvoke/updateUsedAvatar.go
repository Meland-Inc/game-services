package grpcInvoke

import (
	"encoding/json"
	"fmt"
	"game-message-core/grpc"
	"game-message-core/grpc/methodData"
	"game-message-core/proto"

	"github.com/Meland-Inc/game-services/src/common/daprInvoke"
	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/common/time_helper"
	"github.com/Meland-Inc/game-services/src/global/userAgent"
)

// to scene service update user used avatars and profile
func UpdateUsedAvatar(
	userId int64,
	avatars []proto.PlayerAvatar,
	profile proto.EntityProfile,
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

	input := methodData.UpdateUsedAvatarInput{
		MsgVersion:   time_helper.NowUTCMill(),
		UserId:       userId,
		UsingAvatars: avatars,
		CurProfile:   profile,
	}
	serviceLog.Info("UpdateUsedAvatar = %+v", input)
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
	err = json.Unmarshal(outBytes, output)
	if err != nil {
		serviceLog.Error("UpdateUsedAvatar Unmarshal : err : %+v", err)
		return err
	}
	if !output.Success {
		return fmt.Errorf("UpdateUsedAvatar fail err: %s", output.Success)
	}
	return nil
}
