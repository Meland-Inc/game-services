package msgChannel

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

var chanInstance *ClientMsgChannel

func GetInstance() *ClientMsgChannel {
	if chanInstance == nil {
		NewClientMsgChannel()
	}
	return chanInstance
}

func NewClientMsgChannel() *ClientMsgChannel {
	chanInstance = &ClientMsgChannel{
		msgChan:    make(chan *methodData.PullClientMessageInput, 2048),
		stopChan:   make(chan chan struct{}),
		isClosed:   false,
		msgHandler: make(map[proto.EnvelopeType]HandleFunc),
	}

	chanInstance.registerHandler()
	chanInstance.run()
	return chanInstance
}

type ClientMsgChannel struct {
	isClosed   bool
	msgHandler map[proto.EnvelopeType]HandleFunc
	msgChan    chan *methodData.PullClientMessageInput
	stopChan   chan chan struct{}
}

func (ch *ClientMsgChannel) CallClientMsg(in *methodData.PullClientMessageInput) {
	if ch.isClosed {
		return
	}
	ch.msgChan <- in
}

func (ch *ClientMsgChannel) stop() {
	ch.isClosed = true
	close(ch.stopChan)
	close(ch.msgChan)
}

func (ch *ClientMsgChannel) Stop() {
	stopDone := make(chan struct{}, 1)
	ch.stopChan <- stopDone
	<-stopDone
}

func (ch *ClientMsgChannel) run() {

	go func() {
		if err := recover(); err != nil {
			ch.isClosed = false
			ch.run()
		}

		for {
			select {
			case input := <-ch.msgChan:
				ch.onClientMessage(input)
			case stopFinished := <-ch.stopChan:
				ch.stop()
				stopFinished <- struct{}{}
				return
			}
		}
	}()
}

func (ch *ClientMsgChannel) onClientMessage(input *methodData.PullClientMessageInput) {
	msg, err := protoTool.UnMarshalToEnvelope(input.MsgBody)
	if err != nil {
		serviceLog.Error("account Unmarshal Envelope fail err: %+v", err)
		return
	}

	if handler, exist := ch.msgHandler[msg.Type]; exist {
		handler(input, msg)
	}
}

func (ch *ClientMsgChannel) SendToPlayer(
	agentAppId, socketId string,
	userId int64,
	msg *proto.Envelope,
) error {
	msgBody, err := protoTool.MarshalProto(msg)
	if err != nil {
		return err
	}

	input := methodData.BroadCastToClientInput{
		MsgVersion:   time_helper.NowUTCMill(),
		ServiceAppId: serviceCnf.GetInstance().ServerName,
		UserId:       userId,
		SocketId:     socketId,
		MsgId:        int32(msg.Type),
		MsgBody:      msgBody,
	}

	inputBytes, err := json.Marshal(input)
	if err != nil {
		serviceLog.Error("SendToPlayer Marshal BroadCastInput failed err: %+v", err)
		return err
	}

	_, err = daprInvoke.InvokeMethod(
		agentAppId,
		string(grpc.ProtoMessageActionBroadCastToClient),
		inputBytes,
	)
	if err != nil {
		serviceLog.Error("SendToPlayer InvokeMethod  failed err:%+v", err)
	}

	return err
}
