package daprEvent

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/Meland-Inc/game-services/src/common/daprInvoke"
	"github.com/Meland-Inc/game-services/src/common/serviceLog"
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
	bs, err := json.Marshal(e.Data)
	if err != nil {
		serviceLog.Error("DemoServiceTestEventHandler  marshal e.Data  fail err: %+v", err)
		return false, fmt.Errorf("DemoServiceTestEventHandler  marshal e.Data  fail err: %+v", err)
	}

	escStr, err := url.QueryUnescape(string(bs))
	serviceLog.Info("Receive DemoServiceTestEventHandler data: %v, err: %v", escStr, err)
	return false, nil
}
