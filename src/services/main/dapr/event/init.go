package daprEvent

import (
	"context"

	"game-message-core/grpc"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/global/module"
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
		string(message.SubscriptionEventUpdateUserNFT), module.MODULE_NAME_PLAYER_DATA,
	)); err != nil {
		return err
	}

	if err := daprInvoke.AddTopicEventHandler(makePubsubEventHandler(
		string(message.SubscriptionEventMultiUpdateUserNFT), module.MODULE_NAME_PLAYER_DATA,
	)); err != nil {
		return err
	}

	if err := daprInvoke.AddTopicEventHandler(makePubsubEventHandler(
		string(message.SubscriptionEventMultiLandDataUpdateEvent), module.MODULE_NAME_LAND,
	)); err != nil {
		return err
	}

	if err := daprInvoke.AddTopicEventHandler(makePubsubEventHandler(
		string(message.SubscriptionEventMultiRecyclingEvent), module.MODULE_NAME_LAND,
	)); err != nil {
		return err
	}

	if err := daprInvoke.AddTopicEventHandler(makePubsubEventHandler(
		string(message.SubscriptionEventMultiBuildUpdateEvent), module.MODULE_NAME_LAND,
	)); err != nil {
		return err
	}

	return nil
}

func initServiceGrpcPubsubEventHandle() error {
	if err := daprInvoke.AddTopicEventHandler(makePubsubEventHandler(
		string(grpc.SubscriptionEventUserEnterGame), module.MODULE_NAME_PLAYER_DATA,
	)); err != nil {
		return err
	}

	if err := daprInvoke.AddTopicEventHandler(makePubsubEventHandler(
		string(grpc.SubscriptionEventUserLeaveGame), module.MODULE_NAME_PLAYER_DATA,
	)); err != nil {
		return err
	}

	if err := daprInvoke.AddTopicEventHandler(makePubsubEventHandler(
		string(grpc.SubscriptionEventSavePlayerData), module.MODULE_NAME_PLAYER_DATA,
	)); err != nil {
		return err
	}

	if err := daprInvoke.AddTopicEventHandler(makePubsubEventHandler(
		string(grpc.SubscriptionEventKillMonster), module.MODULE_NAME_PLAYER_DATA,
	)); err != nil {
		return err
	}

	if err := daprInvoke.AddTopicEventHandler(makePubsubEventHandler(
		string(grpc.SubscriptionEventPlayerDeath), module.MODULE_NAME_PLAYER_DATA,
	)); err != nil {
		return err
	}

	if err := daprInvoke.AddTopicEventHandler(makePubsubEventHandler(
		string(grpc.SubscriptionEventUserTaskReward), module.MODULE_NAME_PLAYER_DATA,
	)); err != nil {
		return err
	}

	if err := daprInvoke.AddTopicEventHandler(makePubsubEventHandler(
		string(grpc.SubscriptionEventSaveHomeData), module.MODULE_NAME_HOME,
	)); err != nil {
		return err
	}

	if err := daprInvoke.AddTopicEventHandler(makePubsubEventHandler(
		string(grpc.SubscriptionEventGranaryStockpile), module.MODULE_NAME_HOME,
	)); err != nil {
		return err
	}

	if err := daprInvoke.AddTopicEventHandler(makePubsubEventHandler(
		string(grpc.SubscriptionEventUserChangeService), module.MODULE_NAME_PLAYER_DATA,
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
		model, exist := module.GetModel(modelName)
		if !exist {
			serviceLog.Error("model [%s] not found", modelName)
			return false, nil
		}

		model.EventCallNoReturn(module.NewModuleEventReq(name, e, false, nil))
		return false, nil
	}
}
