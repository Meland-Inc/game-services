package daprService

import (
	"github.com/Meland-Inc/game-services/src/common/daprInvoke"
	mgrDaprCalls "github.com/Meland-Inc/game-services/src/services/manager/dapr/calls"
	mgrDaprEvent "github.com/Meland-Inc/game-services/src/services/manager/dapr/event"
)

func Init() (err error) {
	if err = daprInvoke.InitClient("5700"); err != nil {
		return err
	}

	if err = daprInvoke.InitServer("5770"); err != nil {
		return err
	}

	if err = mgrDaprEvent.InitDaprPubsubEvent(); err != nil {
		return err
	}

	if err = mgrDaprCalls.InitDaprCallHandle(); err != nil {
		return err
	}

	return nil
}
