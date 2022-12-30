package grpcPubsubEvent

import (
	"encoding/json"

	"github.com/Meland-Inc/game-services/src/common/daprInvoke"
	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	message "github.com/Meland-Inc/game-services/src/global/web3Message"
)

func Web3RPCEventCloseDynamicSceneService(serAppId string) error {
	env := message.CloseServer{
		ServerAppId: serAppId,
	}

	inputBytes, err := json.Marshal(env)
	if err != nil {
		serviceLog.Error("Web3 CloseServer Event Marshal failed err: %+v", err)
		return err
	}

	serviceLog.Info("Web3 CloseServer Event : %+v", env)

	return daprInvoke.PubSubEventCall(string(message.SubscriptionEventCloseServer), string(inputBytes))
}
