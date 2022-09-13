package userChannel

import (
	"game-message-core/proto"
	"game-message-core/protoTool"

	"github.com/Meland-Inc/game-services/src/common/net/session"
	"github.com/Meland-Inc/game-services/src/common/serviceLog"
)

type UserChannel struct {
	id                string
	owner             int64
	sceneServiceAppId string
	tcpSession        *session.Session
	channels          []chan []byte
	stopChans         []chan chan struct{}
}

func NewUserChannel(se *session.Session) *UserChannel {
	uc := &UserChannel{}
	uc.tcpSession = se
	uc.id = se.SessionId()

	count := int(proto.ServiceType_ServiceTypeLimit)
	uc.channels = make([]chan []byte, count, count)
	uc.stopChans = make([]chan chan struct{}, count, count)
	for i := 0; i < len(uc.channels); i++ {
		uc.channels[i] = make(chan []byte, 256)
		uc.stopChans[i] = make(chan chan struct{})
	}
	return uc
}

func (uc *UserChannel) GetId() string { return uc.id }
func (uc *UserChannel) SetOwner(owner int64) {
	uc.owner = owner
}
func (uc *UserChannel) GetOwner() int64 {
	return uc.owner
}
func (uc *UserChannel) GetSession() *session.Session {
	return uc.tcpSession
}
func (uc *UserChannel) GetSceneService() string {
	return uc.sceneServiceAppId
}
func (uc *UserChannel) SetSceneService(sceneServiceAppId string) {
	uc.sceneServiceAppId = sceneServiceAppId
}

func (uc *UserChannel) OnSessionReceivedData(data []byte) {
	msg, err := uc.UnMarshalProtoMessage(data)
	if err != nil {
		return
	}

	serviceId := protoTool.EnvelopeTypeToServiceType(msg.Type)
	uc.channels[serviceId] <- data
}

func (uc *UserChannel) OnSessionClose() {
	uc.callPlayerLeaveGame()
	// todo: user channel mgr remove self
	uc.Stop()
}

func (uc *UserChannel) Stop() {
	for _, sh := range uc.stopChans {
		stopDone := make(chan struct{}, 1)
		sh <- stopDone
		<-stopDone
	}
}

func (uc *UserChannel) Run() {
	for idx, ch := range uc.channels {
		uc.runChannel(idx, ch, uc.stopChans[idx])
	}
}

func (uc *UserChannel) runChannel(channelId int, ch chan []byte, stopCh chan chan struct{}) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				serviceLog.Error("user channel [%d] panic err: %v", channelId, err)
				uc.runChannel(channelId, ch, stopCh)
			}
		}()

		for {
			select {
			case data := <-ch:
				uc.onProtoData(data)
			case stopDone := <-stopCh:
				stopDone <- struct{}{}
				return
			}
		}
	}()
}

func (uc *UserChannel) onProtoData(data []byte) {
	msg, err := uc.UnMarshalProtoMessage(data)
	if err != nil {
		return
	}

	serviceType := protoTool.EnvelopeTypeToServiceType(msg.Type)
	if serviceType == proto.ServiceType_ServiceTypeAgent {
		uc.agentClientMsg(msg)
	} else {
		uc.callOtherServiceClientMsg(data, msg)
	}
}
