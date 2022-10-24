package daprEvent

import (
	"game-message-core/grpc"

	message "github.com/Meland-Inc/game-services/src/global/web3Message"

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
	if err := daprInvoke.AddTopicEventHandler(
		string(message.SubscriptionEventUpdateUserNFT),
		Web3UpdateUserNftHandler,
	); err != nil {
		return err
	}

	if err := daprInvoke.AddTopicEventHandler(
		string(message.SubscriptionEventMultiUpdateUserNFT),
		Web3MultiUpdateUserNftHandler,
	); err != nil {
		return err
	}

	if err := daprInvoke.AddTopicEventHandler(
		string(message.SubscriptionEventMultiLandDataUpdateEvent),
		Web3MultiLandDataUpdateEventHandler,
	); err != nil {
		return err
	}

	if err := daprInvoke.AddTopicEventHandler(
		string(message.SubscriptionEventRecyclingEvent),
		Web3RecyclingEventHandler,
	); err != nil {
		return err
	}

	if err := daprInvoke.AddTopicEventHandler(
		string(message.SubscriptionEventBuildUpdateEvent),
		Web3BuildUpdateEventHandler,
	); err != nil {
		return err
	}

	return nil
}

func initServiceGrpcPubsubEventHandle() error {
	if err := daprInvoke.AddTopicEventHandler(
		string(grpc.SubscriptionEventUserEnterGame),
		UserEnterGameEventHandler,
	); err != nil {
		return err
	}

	if err := daprInvoke.AddTopicEventHandler(
		string(grpc.SubscriptionEventSavePlayerData),
		SavePlayerDataEventHandle,
	); err != nil {
		return err
	}

	if err := daprInvoke.AddTopicEventHandler(
		string(grpc.SubscriptionEventKillMonster),
		PlayerKillMonsterEventHandle,
	); err != nil {
		return err
	}

	if err := daprInvoke.AddTopicEventHandler(
		string(grpc.SubscriptionEventPlayerDeath),
		PlayerDeathEventHandle,
	); err != nil {
		return err
	}

	if err := daprInvoke.AddTopicEventHandler(
		string(grpc.SubscriptionEventUserTaskReward),
		TaskRewardEventHandler,
	); err != nil {
		return err
	}

	return nil
}
