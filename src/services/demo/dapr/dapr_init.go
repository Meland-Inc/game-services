package daprService

import (
	"github.com/Meland-Inc/game-services/src/common/daprInvoke"
	demoDaprCalls "github.com/Meland-Inc/game-services/src/services/demo/dapr/calls"
	demoDaprEvent "github.com/Meland-Inc/game-services/src/services/demo/dapr/event"
)

func Init() (err error) {
	if err = daprInvoke.InitClient("5700"); err != nil {
		return err
	}

	if err = daprInvoke.InitServer("5770"); err != nil {
		return err
	}

	if err = demoDaprEvent.InitDaprPubsubEvent(); err != nil {
		return err
	}

	if err = demoDaprCalls.InitDaprCallHandle(); err != nil {
		return err
	}

	return nil
}
