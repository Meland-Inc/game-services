package msgChannel

import (
	"game-message-core/proto"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
)

var chanInstance *MsgChannel

type MsgChannel struct {
	isClosed         bool
	clientMsgHandler map[proto.EnvelopeType]HandleFunc
	clientMsgChan    chan *proto.PullClientMessageInput
	serviceMsgChan   chan *ServiceMsgData
	stopChan         chan chan struct{}
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
		stopChan:         make(chan chan struct{}),
		isClosed:         false,
		clientMsgChan:    make(chan *proto.PullClientMessageInput, 2048),
		clientMsgHandler: make(map[proto.EnvelopeType]HandleFunc),
		serviceMsgChan:   make(chan *ServiceMsgData, 2048),
	}

	channel.registerClientMsgHandler()
	return channel
}

func (ch *MsgChannel) CallServiceMsg(in *ServiceMsgData) {
	if ch.isClosed {
		return
	}
	ch.serviceMsgChan <- in
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
	close(ch.serviceMsgChan)
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
			case serviceMsg := <-ch.serviceMsgChan:
				ch.onServiceMessage(serviceMsg)
			case stopFinished := <-ch.stopChan:
				ch.stop()
				stopFinished <- struct{}{}
				return
			}
		}
	}()
}

func (ch *MsgChannel) onClientMessage(input *proto.PullClientMessageInput) {
	serviceLog.Info("received player[%v] message: %v", input.UserId, input.Msg.Type)

	if handler, exist := ch.clientMsgHandler[input.Msg.Type]; exist {
		handler(input)
	}
}
