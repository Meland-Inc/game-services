package userChannel

import (
	"encoding/json"
	"game-message-core/grpc"
	"game-message-core/grpc/methodData"
	"game-message-core/proto"
	"game-message-core/protoTool"

	"github.com/Meland-Inc/game-services/src/common/daprInvoke"
	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/common/time_helper"
	"github.com/Meland-Inc/game-services/src/global/serviceCnf"
)

func (uc *UserChannel) makeClientProtoMessageInputBytes(data []byte) ([]byte, error) {
	input := methodData.ClientProtoMessageInput{
		MsgVersion:    time_helper.NowUTCMill(),
		AgentAppId:    serviceCnf.GetInstance().ServerName,
		UserId:        uc.owner,
		ProtoBytesReq: data,
	}

	inputBytes, err := json.Marshal(input)
	if err != nil {
		serviceLog.Error("Marshal client msg input failed err:+v", err)
	}
	return inputBytes, err
}

func (uc *UserChannel) parseClientProtoMessageOutput(data []byte) (*methodData.ClientProtoMessageOutput, error) {
	output := &methodData.ClientProtoMessageOutput{}
	err := json.Unmarshal(data, output)
	if err != nil {
		serviceLog.Error("Unmarshal client msg output failed err:+v", err)
	}
	return output, err
}
func (uc *UserChannel) callOtherServiceClientMsg(data []byte, msg *proto.Envelope) {
	errResponseF := func(errorCode int32, errMsg string) {
		resMsg := &proto.Envelope{
			Type:         msg.Type,
			SeqId:        msg.SeqId,
			ErrorCode:    errorCode,
			ErrorMessage: errMsg,
		}
		if byes, err := uc.MarshalProtoMessage(resMsg); err == nil {
			uc.tcpSession.Write(byes)
		}
	}

	serviceType := protoTool.EnvelopeTypeToServiceType(msg.Type)
	appId := uc.getServiceAppId(serviceType)
	if appId == "" {
		serviceLog.Error("other service msg APPID is nil")
		errResponseF(70001, "other service msg APPID is nil")
		return
	}

	inputBytes, err := uc.makeClientProtoMessageInputBytes(data)
	if err != nil {
		serviceLog.Error("make client proto msg bytes failed err:+v", err)
		errResponseF(70001, "make client proto msg bytes failed")
		return
	}

	serviceLog.Info("UserChannel call [%v]", serviceType)

	resp, err := daprInvoke.InvokeMethod(appId, string(grpc.GameServiceActionClientProtoMessage), inputBytes)
	if err != nil {
		serviceLog.Error("send client mst to main service failed err:+v", err)
		errResponseF(70001, err.Error())
		return
	}
	output, err := uc.parseClientProtoMessageOutput(resp)
	if err != nil {
		errResponseF(70001, err.Error())
		return
	}

	serviceLog.Info("UserChannel call [%v] resp msg: %+v", serviceType, output)
	uc.tcpSession.Write(output.ProtoBytesResp)

	// TODO: set channel linked scene service info from SingIn message
}

func (uc *UserChannel) callPlayerLeaveGame() {
	input := methodData.PlayerLeaveGameInput{
		MsgVersion: time_helper.NowUTCMill(),
		AgentAppId: serviceCnf.GetInstance().ServerName,
		UserId:     uc.owner,
	}

	inputBytes, err := json.Marshal(input)
	if err != nil {
		serviceLog.Error("Marshal player leave input failed err:+v", err)
		return
	}

	serviceLog.Info("call player leave game: %+v", input)

	if _, err = daprInvoke.InvokeMethod(
		string(grpc.AppIdMelandServiceChat),
		string(grpc.GameServiceActionPlayerLeaveGame),
		inputBytes,
	); err != nil {
		serviceLog.Error("call chat service player leave game failed err:+v", err)
	}

	if _, err = daprInvoke.InvokeMethod(
		uc.sceneServiceAppId,
		string(grpc.GameServiceActionPlayerLeaveGame),
		inputBytes,
	); err != nil {
		serviceLog.Error("call chat service player leave game failed err:+v", err)
	}
}
