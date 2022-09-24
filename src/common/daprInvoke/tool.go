package daprInvoke

import (
	"encoding/json"
	"game-message-core/protoTool"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/dapr/go-sdk/service/common"
	googleProto "google.golang.org/protobuf/proto"
)

func MakeJsonOutputContent(in *common.InvocationEvent, resp interface{}) (*common.Content, error) {
	bytes, err := json.Marshal(resp)
	if err != nil {
		serviceLog.Error("make output content fail marshal err : %+v", err)
		return nil, err
	}
	out := &common.Content{
		Data:        bytes,
		ContentType: in.ContentType,
		DataTypeURL: in.DataTypeURL,
	}
	return out, nil
}

func MakeProtoOutputContent(in *common.InvocationEvent, resp googleProto.Message) (*common.Content, error) {
	bytes, err := protoTool.MarshalProto(resp)
	if err != nil {
		serviceLog.Error("make output content fail marshal err : %+v", err)
		return nil, err
	}
	out := &common.Content{
		Data:        bytes,
		ContentType: in.ContentType,
		DataTypeURL: in.DataTypeURL,
	}
	return out, nil
}
