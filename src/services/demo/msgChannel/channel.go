package msgChannel

import (
	"game-message-core/grpc"
	"game-message-core/proto"
	"game-message-core/protoTool"

	"github.com/Meland-Inc/game-services/src/common/daprInvoke"
	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/common/time_helper"
	"github.com/Meland-Inc/game-services/src/global/serviceCnf"
)

var chanInstance *MsgChannel

type MsgChannel struct {
	isClosed      bool
	msgHandler    map[proto.EnvelopeType]HandleFunc
	clientMsgChan chan *proto.PullClientMessageInput
	stopChan      chan chan struct{}
}

func GetInstance() *MsgChannel {
	if chanInstance == nil {
		InitAndRun()
	}
	return chanInstance
}

func InitAndRun() {
	chanInstance = NewMsgChannel()
	chanInstance.run()
}

func NewMsgChannel() *MsgChannel {
	channel := &MsgChannel{
		clientMsgChan: make(chan *proto.PullClientMessageInput, 2048),
		stopChan:      make(chan chan struct{}),
		isClosed:      false,
		msgHandler:    make(map[proto.EnvelopeType]HandleFunc),
	}
	channel.registerHandler()
	return channel
}

func (ch *MsgChannel) CallClientMsg(in *proto.PullClientMessageInput) {
	if ch.isClosed {
		return
	}
	ch.clientMsgChan <- in
}

func (ch *MsgChannel) stop() {
	ch.isClosed = true
	close(ch.stopChan)
	close(ch.clientMsgChan)
}

func (ch *MsgChannel) Stop() {
	stopDone := make(chan struct{}, 1)
	ch.stopChan <- stopDone
	<-stopDone
}

func (ch *MsgChannel) run() {

	go func() {
		defer func() {
			if err := recover(); err != nil {
				serviceLog.Error("msg channel recover panic err: %+v", err)
				ch.isClosed = false
				ch.run()
			}
		}()

		for {
			select {
			case input := <-ch.clientMsgChan:
				ch.onClientMessage(input)
			case stopFinished := <-ch.stopChan:
				ch.stop()
				stopFinished <- struct{}{}
				return
			}
		}
	}()
}

func (ch *MsgChannel) onClientMessage(input *proto.PullClientMessageInput) {
	serviceLog.Info("client msg: %+v", input)
	if handler, exist := ch.msgHandler[input.Msg.Type]; exist {
		handler(input)
	}
}

func (ch *MsgChannel) SendToPlayer(agentAppId, socketId string, userId int64, msg *proto.Envelope) error {
	input := &proto.BroadCastToClientInput{
		MsgVersion:   time_helper.NowUTCMill(),
		ServiceAppId: serviceCnf.GetInstance().ServerName,
		UserId:       userId,
		SocketId:     socketId,
		MsgId:        int32(msg.Type),
		Msg:          msg,
	}

	inputBytes, err := protoTool.MarshalProto(input)
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
