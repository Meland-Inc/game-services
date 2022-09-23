package daprInvoke

import (
	"game-message-core/protoTool"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/dapr/go-sdk/service/common"
	googleProto "google.golang.org/protobuf/proto"
)

func MakeOutputContent(in *common.InvocationEvent, resp googleProto.Message) (*common.Content, error) {
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
