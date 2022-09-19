package daprEvent

import (
	"game-message-core/grpc"

	message "github.com/Meland-Inc/game-services/src/global/web3Message"

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

	return nil
}
