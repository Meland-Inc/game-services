package daprCalls

import (
	"github.com/Meland-Inc/game-services/src/common/daprInvoke"
)

func InitDaprCallHandle() (err error) {
	daprInvoke.AddServiceInvocationHandler("DemoServiceTestCallsHandler", DemoServiceTestCallsHandler)
	if err != nil {
		return err
	}

	return nil
}
