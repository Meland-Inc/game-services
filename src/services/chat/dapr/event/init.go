package daprEvent

import (
	"game-message-core/grpc"

	"github.com/Meland-Inc/game-services/src/common/daprInvoke"
)

func InitDaprPubsubEvent() (err error) {
	if err := initServiceGrpcPubsubEventHandle(); err != nil {
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
		string(grpc.SubscriptionEventUserLeaveGame),
		UserLeaveGameHandler,
	); err != nil {
		return err
	}

	if err := daprInvoke.AddTopicEventHandler(
		string(grpc.SubscriptionEventTickOutPlayer),
		TickOutUserHandler,
	); err != nil {
		return err
	}

	if err := daprInvoke.AddTopicEventHandler(
		string(grpc.SubscriptionEventSavePlayerData),
		SavePlayerDataEventHandler,
	); err != nil {
		return err
	}

	return nil
}
