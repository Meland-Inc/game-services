package daprEvent

import (
	"context"
	"fmt"

	"github.com/Meland-Inc/game-services/src/common/daprInvoke"
	"github.com/dapr/go-sdk/service/common"
)

func InitDaprPubsubEvent() (err error) {
	daprInvoke.AddTopicEventHandler("DemoServiceTestEventHandler", DemoServiceTestEventHandler)
	if err != nil {
		return err
	}


	
	return nil
}

func DemoServiceTestEventHandler(ctx context.Context, e *common.TopicEvent) (retry bool, err error) {
	fmt.Println("this is DemoServiceTestEvent")
	return false, nil
}
