package userChannel

import (
	"encoding/json"
	"fmt"
	"game-message-core/grpc"
	"game-message-core/grpc/methodData"
	"game-message-core/proto"
	"game-message-core/protoTool"

	"github.com/Meland-Inc/game-services/src/common/daprInvoke"
	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/common/time_helper"
	"github.com/Meland-Inc/game-services/src/global/grpcAPI/grpcPubsubEvent"
	"github.com/Meland-Inc/game-services/src/global/serviceCnf"
)

func (uc *UserChannel) clientMsgIsLegal(msgType proto.EnvelopeType) (bool, error) {
	if uc.enterSceneService {
		return true, nil
	}
	switch msgType {
	case proto.EnvelopeType_SigninPlayer,
		proto.EnvelopeType_EnterMap,
		proto.EnvelopeType_QueryPlayer,
		proto.EnvelopeType_CreatePlayer:
		return true, nil
	default:
		return false, fmt.Errorf("msg is not Legal: %v", msgType)
	}
}

func (uc *UserChannel) makePullClientMessageInputBytes(data []byte) ([]byte, error) {
	msg, _ := protoTool.UnMarshalToEnvelope(data)
	input := &methodData.PullClientMessageInput{
		MsgVersion: time_helper.NowUTCMill(),
		AgentAppId: serviceCnf.GetInstance().AppId,
		UserId:     uc.owner,
		SocketId:   uc.id,
		MsgId:      int32(msg.Type),
		MsgBody:    data,
	}
	serviceLog.Info("pull client msg[%+v] to other service", msg.Type)

	inputBytes, err := json.Marshal(input)
	if err != nil {
		serviceLog.Error("Marshal client msg input failed err:+v", err)
	}
	return inputBytes, err
}

func (uc *UserChannel) parsePullClientMessageOutput(data []byte) (*methodData.PullClientMessageOutput, error) {
	output := &methodData.PullClientMessageOutput{}
	err := json.Unmarshal(data, output)
	if err != nil {
		serviceLog.Error("Unmarshal client msg output failed err:+v", err)
	}
	return output, err
}

func (uc *UserChannel) getServiceAppId(serviceType proto.ServiceType) (appId string) {
	switch proto.ServiceType(serviceType) {
	case proto.ServiceType_ServiceTypeMain:
		appId = string(grpc.GAME_SERVICE_APPID_MAIN)
	case proto.ServiceType_ServiceTypeAccount:
		appId = string(grpc.GAME_SERVICE_APPID_ACCOUNT)
	case proto.ServiceType_ServiceTypeScene:
		appId = uc.sceneServiceAppId
	case proto.ServiceType_ServiceTypeTask:
		appId = string(grpc.GAME_SERVICE_APPID_TASK)
	case proto.ServiceType_ServiceTypeChat:
		appId = string(grpc.GAME_SERVICE_APPID_CHAT)
	default:
	}
	return
}

func (uc *UserChannel) callOtherServiceClientMsg(data []byte, msg *proto.Envelope) {
	errResponseF := func(errorCode int32, errMsg string) {
		resMsg := &proto.Envelope{
			Type:         msg.Type,
			SeqId:        msg.SeqId,
			ErrorCode:    errorCode,
			ErrorMessage: errMsg,
		}
		if byes, err := protoTool.MarshalEnvelope(resMsg); err == nil {
			uc.tcpSession.Write(byes)
		}
	}

	if _, err := uc.clientMsgIsLegal(msg.Type); err != nil {
		errResponseF(70000, err.Error())
		return
	}

	serviceType := protoTool.EnvelopeTypeToServiceType(msg.Type)
	appId := uc.getServiceAppId(serviceType)
	if appId == "" {
		serviceLog.Error("other service msg APPID is nil")
		errResponseF(70001, "other service msg APPID is nil")
		return
	}

	inputBytes, err := uc.makePullClientMessageInputBytes(data)
	if err != nil {
		serviceLog.Error("make client proto msg bytes failed err:+v", err)
		errResponseF(70001, "make client proto msg bytes failed")
		return
	}

	serviceLog.Info("UserChannel call [%v]", serviceType)

	resp, err := daprInvoke.InvokeMethod(
		appId,
		string(grpc.ProtoMessageActionPullClientMessage),
		inputBytes,
	)
	if err != nil {
		serviceLog.Error("send client msg to [%s] failed err: %+v", appId, err)
		errResponseF(70001, err.Error())
		return
	}

	output, err := uc.parsePullClientMessageOutput(resp)
	serviceLog.Info("UserChannel call [%v] resp msg: %+v", serviceType, output)
	if err != nil {
		errResponseF(70001, err.Error())
		return
	}
	if !output.Success {
		errResponseF(70001, output.ErrMsg)
		return
	}
}

func (uc *UserChannel) callPlayerLeaveGame() {
	err := grpcPubsubEvent.RPCPubsubEventLeaveGame(uc.owner)
	if err != nil {
		serviceLog.Error("call [%s] UserActionLeaveGame failed err: %+v", serviceCnf.GetInstance().AppId, err)
	}
}
