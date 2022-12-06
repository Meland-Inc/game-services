package daprEvent

import (
	"context"

	"game-message-core/grpc"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/global/component"
	message "github.com/Meland-Inc/game-services/src/global/web3Message"
	"github.com/dapr/go-sdk/service/common"

	"github.com/Meland-Inc/game-services/src/common/daprInvoke"
)

func InitDaprPubsubEvent() (err error) {
	if err := initWeb3ServicesPubsubEventHandler(); err != nil {
		return err
	}

	if err := initServiceGrpcPubsubEventHandle(); err != nil {
		return err
	}

	return nil
}

func initWeb3ServicesPubsubEventHandler() error {
	if err := daprInvoke.AddTopicEventHandler(makePubsubEventHandler(
		string(message.SubscriptionEventUpdateUserNFT), component.MODEL_NAME_PLAYER_DATA,
	)); err != nil {
		return err
	}

	if err := daprInvoke.AddTopicEventHandler(makePubsubEventHandler(
		string(message.SubscriptionEventMultiUpdateUserNFT), component.MODEL_NAME_PLAYER_DATA,
	)); err != nil {
		return err
	}

	if err := daprInvoke.AddTopicEventHandler(makePubsubEventHandler(
		string(message.SubscriptionEventMultiLandDataUpdateEvent), component.MODEL_NAME_LAND,
	)); err != nil {
		return err
	}

	if err := daprInvoke.AddTopicEventHandler(makePubsubEventHandler(
		string(message.SubscriptionEventMultiRecyclingEvent), component.MODEL_NAME_LAND,
	)); err != nil {
		return err
	}

	if err := daprInvoke.AddTopicEventHandler(makePubsubEventHandler(
		string(message.SubscriptionEventMultiBuildUpdateEvent), component.MODEL_NAME_LAND,
	)); err != nil {
		return err
	}

	return nil
}

func initServiceGrpcPubsubEventHandle() error {
	if err := daprInvoke.AddTopicEventHandler(makePubsubEventHandler(
		string(grpc.SubscriptionEventUserEnterGame), component.MODEL_NAME_PLAYER_DATA,
	)); err != nil {
		return err
	}

	if err := daprInvoke.AddTopicEventHandler(makePubsubEventHandler(
		string(grpc.SubscriptionEventUserLeaveGame), component.MODEL_NAME_PLAYER_DATA,
	)); err != nil {
		return err
	}

	if err := daprInvoke.AddTopicEventHandler(makePubsubEventHandler(
		string(grpc.SubscriptionEventSavePlayerData), component.MODEL_NAME_PLAYER_DATA,
	)); err != nil {
		return err
	}

	if err := daprInvoke.AddTopicEventHandler(makePubsubEventHandler(
		string(grpc.SubscriptionEventKillMonster), component.MODEL_NAME_PLAYER_DATA,
	)); err != nil {
		return err
	}

	if err := daprInvoke.AddTopicEventHandler(makePubsubEventHandler(
		string(grpc.SubscriptionEventPlayerDeath), component.MODEL_NAME_PLAYER_DATA,
	)); err != nil {
		return err
	}

	if err := daprInvoke.AddTopicEventHandler(makePubsubEventHandler(
		string(grpc.SubscriptionEventUserTaskReward), component.MODEL_NAME_PLAYER_DATA,
	)); err != nil {
		return err
	}

	if err := daprInvoke.AddTopicEventHandler(makePubsubEventHandler(
		string(grpc.SubscriptionEventSaveHomeData), component.MODEL_NAME_HOME,
	)); err != nil {
		return err
	}

	return nil
}

func makePubsubEventHandler(name string, modelName string) (
	string, func(ctx context.Context, e *common.TopicEvent) (retry bool, err error),
) {
	return name, func(ctx context.Context, e *common.TopicEvent) (retry bool, err error) {
		serviceLog.Info("pubsub event[%s] [%s] data:%v", name, modelName, e.Data)
		model, exist := component.GetInstance().GetModel(modelName)
		if !exist {
			serviceLog.Error("model [%s] not found", modelName)
			return false, nil
		}

		model.EventCallNoReturn(&component.ModelEventReq{
			EventType: name,
			Msg:       e,
		})
		return false, nil
	}
}
