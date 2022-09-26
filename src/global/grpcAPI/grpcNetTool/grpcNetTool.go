package grpcNetTool

import (
	"encoding/json"
	"net/url"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/dapr/go-sdk/service/common"
)

func UnmarshalGrpcData(data []byte, v interface{}) error {
	err := json.Unmarshal(data, v)
	if err != nil {
		escStr, err := url.QueryUnescape(string(data))
		serviceLog.Info("QueryUnescape data: %v, err: %v", escStr, err)
		if err != nil {
			serviceLog.Info("QueryUnescape data failed err: %v", err)
			return err
		}
		err = json.Unmarshal([]byte(escStr), v)
		if err != nil {
			serviceLog.Error("Unmarshal QueryUnescape failed err:%+v", err)
			return err
		}
	}
	return err
}

func UnmarshalGrpcTopicEvent(e *common.TopicEvent, v interface{}) error {
	inputJson, ok := e.Data.(string)
	if ok {
		escStr, err := url.QueryUnescape(inputJson)
		if err != nil {
			serviceLog.Info("QueryUnescape TopicEvent failed err: %v", err)
			return err
		}
		err = json.Unmarshal([]byte(escStr), v)
		if err == nil {
			return nil
		}
	}

	inputBytes, err := json.Marshal(e.Data)
	if err != nil {
		serviceLog.Info("TopicEvent Marshal(e.Data) fail err: %v ", err)
		return err
	}

	return UnmarshalGrpcData(inputBytes, v)
}
