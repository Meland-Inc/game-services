package userChannel

import (
	"game-message-core/grpc"
	"game-message-core/proto"
	"game-message-core/protoTool"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
)

func (uc *UserChannel) UnMarshalProtoMessage(data []byte) (*proto.Envelope, error) {
	serviceLog.Info("userChannel unMarshal proto data  [%v]", data)
	msg := &proto.Envelope{}
	err := protoTool.UnmarshalProto(data, msg)
	serviceLog.Info("userChannel unMarshal proto data err:%+v,  msg:%+v", err, msg)
	if err != nil {
		serviceLog.Error("userChannel Unmarshal Proto msg failed err:+v", err)
		return nil, err
	}
	return msg, nil
}

func (uc *UserChannel) MarshalProtoMessage(msg *proto.Envelope) ([]byte, error) {
	bs, err := protoTool.MarshalProto(msg)
	if err != nil {
		serviceLog.Error("userChannel marshal Proto msg failed err:+v", err)
		return nil, err
	}
	return bs, err
}

func (uc *UserChannel) getServiceAppId(serviceType proto.ServiceType) (appId string) {
	switch proto.ServiceType(serviceType) {
	case proto.ServiceType_ServiceTypeMain:
		appId = string(grpc.AppIdMelandServiceMain)
	case proto.ServiceType_ServiceTypeAccount:
		appId = string(grpc.AppIdMelandServiceAccount)
	case proto.ServiceType_ServiceTypeScene:
		appId = uc.sceneServiceAppId
	case proto.ServiceType_ServiceTypeTask:
		appId = string(grpc.AppIdMelandServiceTask)
	case proto.ServiceType_ServiceTypeChat:
		appId = string(grpc.AppIdMelandServiceChat)
	default:
	}
	return
}
