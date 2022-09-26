package grpcNetTool

import (
	"encoding/json"
	"net/url"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
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
