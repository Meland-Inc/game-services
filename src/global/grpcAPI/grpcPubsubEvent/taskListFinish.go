package grpcPubsubEvent

import (
	"encoding/json"
	"game-message-core/grpc"
	base_data "game-message-core/grpc/baseData"
	"game-message-core/grpc/pubsubEventData"
	"game-message-core/proto"

	"github.com/Meland-Inc/game-services/src/common/daprInvoke"
	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/common/time_helper"
)

func RPCPubsubEventTaskListFinish(
	userId int64,
	taskListType proto.TaskListType,
	rewardItem []*proto.ItemBaseInfo,
) error {
	env := pubsubEventData.TaskListFinishEvent{
		MsgVersion:   time_helper.NowUTCMill(),
		UserId:       userId,
		TaskListType: taskListType,
	}
	for _, item := range rewardItem {
		grpcItem := base_data.GrpcItemBaseInfo{}
		grpcItem.Set(item)
		env.RewardItem = append(env.RewardItem, grpcItem)
	}

	inputBytes, err := json.Marshal(env)
	if err != nil {
		serviceLog.Error("TaskListFinishEvent Marshal Input failed err: %+v", err)
		return err
	}

	serviceLog.Info("pubsubEvent TaskListFinish: %+v", env)

	return daprInvoke.PubSubEventCall(string(grpc.SubscriptionEventTaskListFinish), string(inputBytes))
}
