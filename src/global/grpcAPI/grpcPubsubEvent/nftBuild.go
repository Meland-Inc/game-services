package grpcPubsubEvent

import (
	"encoding/json"
	"game-message-core/grpc"
	base_data "game-message-core/grpc/baseData"
	"game-message-core/grpc/pubsubEventData"

	"github.com/Meland-Inc/game-services/src/common/daprInvoke"
	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/common/time_helper"
)

func RPCPubsubEventNftBuildUpdate(build base_data.GrpcNftBuild) error {
	env := &pubsubEventData.NftBuildUpdateEvent{
		MsgVersion: time_helper.NowUTCMill(),
		Build:      build,
	}
	inputBytes, err := json.Marshal(env)
	if err != nil {
		serviceLog.Error("RPCPubsubEventNftBuildUpdate Marshal Input failed err: %+v", err)
		return err
	}
	return daprInvoke.PubSubEventCall(string(grpc.SubscriptionEventNftBuildUpdate), string(inputBytes))
}

func RPCPubsubEventNftBuildAdd(build base_data.GrpcNftBuild) error {
	env := &pubsubEventData.NftBuildUpdateEvent{
		MsgVersion: time_helper.NowUTCMill(),
		Build:      build,
	}
	inputBytes, err := json.Marshal(env)
	if err != nil {
		serviceLog.Error("RPCPubsubEventNftBuildAdd Marshal Input failed err: %+v", err)
		return err
	}
	return daprInvoke.PubSubEventCall(string(grpc.SubscriptionEventNftBuildAdd), string(inputBytes))
}

func RPCPubsubEventNftBuildRemove(build base_data.GrpcNftBuild) error {
	env := &pubsubEventData.NftBuildRemoveEvent{
		MsgVersion: time_helper.NowUTCMill(),
		EntityId:   build.Id,
		FromNft:    build.FromNft,
		Owner:      build.Owner,
	}
	inputBytes, err := json.Marshal(env)
	if err != nil {
		serviceLog.Error("RPCPubsubEventNftBuildRemove Marshal Input failed err: %+v", err)
		return err
	}
	return daprInvoke.PubSubEventCall(string(grpc.SubscriptionEventNftBuildRemove), string(inputBytes))
}
