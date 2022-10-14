package userChannel

import (
	"game-message-core/grpc/pubsubEventData"
	"game-message-core/proto"
	"game-message-core/protoTool"

	"github.com/Meland-Inc/game-services/src/common/net/session"
	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/common/time_helper"
	"github.com/Meland-Inc/game-services/src/global/grpcAPI/grpcPubsubEvent"
	"github.com/Meland-Inc/game-services/src/global/serviceCnf"
)

// type HandleFunc func(proto.EnvelopeType, msgBody []byte)

type ResMsgData struct {
	MsgType proto.EnvelopeType
	MsgBody []byte
}

type UserChannel struct {
	id                string
	owner             int64
	sceneServiceAppId string
	enterSceneService bool
	tcpSession        *session.Session
	reqMsgChannels    []chan []byte
	stopChans         []chan chan struct{}
	resMsgChannel     chan *ResMsgData
	resMsgStopChannel chan chan struct{}
	isClosed          bool
}

func NewUserChannel(se *session.Session) *UserChannel {
	uc := &UserChannel{}
	uc.tcpSession = se
	uc.id = se.SessionId()
	uc.isClosed = false
	uc.resMsgChannel = make(chan *ResMsgData, 256)
	uc.resMsgStopChannel = make(chan chan struct{})

	count := int(proto.ServiceType_ServiceTypeLimit)
	uc.reqMsgChannels = make([]chan []byte, count, count)
	uc.stopChans = make([]chan chan struct{}, count, count)
	for i := 0; i < len(uc.reqMsgChannels); i++ {
		uc.reqMsgChannels[i] = make(chan []byte, 256)
		uc.stopChans[i] = make(chan chan struct{})
	}
	return uc
}

func (uc *UserChannel) GetId() string                { return uc.id }
func (uc *UserChannel) SetOwner(owner int64)         { uc.owner = owner }
func (uc *UserChannel) GetOwner() int64              { return uc.owner }
func (uc *UserChannel) GetSession() *session.Session { return uc.tcpSession }
func (uc *UserChannel) GetSceneService() string      { return uc.sceneServiceAppId }
func (uc *UserChannel) SetSceneService(sceneServiceAppId string) {
	uc.sceneServiceAppId = sceneServiceAppId
}

func (uc *UserChannel) OnSessionReceivedData(s *session.Session, data []byte) {
	msg, err := protoTool.UnMarshalToEnvelope(data)
	if err != nil {
		return
	}
	
	if msg.Type != proto.EnvelopeType_Ping {
		serviceLog.Debug("user[%d] channel收到客户端消息 [%v]", uc.owner, msg.Type)
	}

	serviceId := protoTool.EnvelopeTypeToServiceType(msg.Type)
	uc.reqMsgChannels[serviceId] <- data
}

func (uc *UserChannel) OnSessionClose(s *session.Session) {
	serviceLog.Info("channel Id:[%s] user:[%d] closed", uc.id, uc.owner)
	uc.callPlayerLeaveGame()
	GetInstance().RemoveUserChannel(uc)
	uc.Stop()
}

func (uc *UserChannel) Stop() {
	for _, sh := range uc.stopChans {
		stopDone := make(chan struct{}, 1)
		sh <- stopDone
		<-stopDone
	}

	stopDone := make(chan struct{}, 1)
	uc.resMsgStopChannel <- stopDone
	<-stopDone

	uc.tcpSession.Stop()
	uc.tcpSession = nil
	uc.isClosed = true
	uc.sceneServiceAppId = ""
	uc.enterSceneService = false
}

func (uc *UserChannel) Run() {
	for idx, ch := range uc.reqMsgChannels {
		uc.runReqMsgChannel(idx, ch, uc.stopChans[idx])
	}
	uc.runResMsgChannel()
}

func (uc *UserChannel) runResMsgChannel() {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				serviceLog.StackError("user resp channel panic err: %v", err)
				uc.runResMsgChannel()
			}
		}()

		for {
			select {
			case data := <-uc.resMsgChannel:
				uc.writeResMsg(data)
			case stopDone := <-uc.resMsgStopChannel:
				stopDone <- struct{}{}
				close(uc.resMsgChannel)
				close(uc.resMsgStopChannel)
				return
			}
		}
	}()
}

func (uc *UserChannel) runReqMsgChannel(channelId int, ch chan []byte, stopCh chan chan struct{}) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				serviceLog.StackError("user req msg channel [%d] panic err: %v", channelId, err)
				uc.runReqMsgChannel(channelId, ch, stopCh)
			}
		}()

		for {
			select {
			case data := <-ch:
				uc.onProtoMessage(data)
			case stopDone := <-stopCh:
				stopDone <- struct{}{}
				close(ch)
				close(stopCh)
				return
			}
		}
	}()
}

func (uc *UserChannel) onProtoMessage(data []byte) {
	msg, err := protoTool.UnMarshalToEnvelope(data)
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

func (uc *UserChannel) writeResMsg(data *ResMsgData) {
	if uc.isClosed || uc.tcpSession == nil || data == nil {
		return
	}

	err := uc.tcpSession.Write(data.MsgBody)
	if err != nil {
		serviceLog.Error("tcp send to user err : %+v", err)
		return
	}

	// update channel owner and sceneServiceAppId by SingInMsg
	switch data.MsgType {
	case proto.EnvelopeType_SigninPlayer:
		uc.onUserSingInGame(data.MsgType, data.MsgBody)
	case proto.EnvelopeType_EnterMap:
		uc.onUserEnterMap(data.MsgBody)
	}

}

func (uc *UserChannel) SendToUser(msgType proto.EnvelopeType, msgBody []byte) {
	if uc.isClosed || uc.tcpSession == nil {
		return
	}

	uc.resMsgChannel <- &ResMsgData{MsgType: msgType, MsgBody: msgBody}
}

func (uc *UserChannel) onUserSingInGame(msgType proto.EnvelopeType, msgBody []byte) {
	respMsg, err := protoTool.UnMarshalToEnvelope(msgBody)
	if err != nil {
		serviceLog.Error("SigninPlayer response message UnMarshal failed")
		return
	}

	if respMsg.ErrorCode > 0 {
		serviceLog.Error("SigninPlayer fail err: %s", respMsg.ErrorMessage)
		uc.Stop()
		return
	}

	payload := respMsg.GetSigninPlayerResponse()
	uc.SetSceneService(payload.GetSceneServiceAppId())
	uc.SetOwner(payload.Player.BaseData.UserId)
	GetInstance().AddUserChannelByOwner(uc)
}

func (uc *UserChannel) onUserEnterMap(msgBody []byte) {
	respMsg, err := protoTool.UnMarshalToEnvelope(msgBody)
	if err != nil {
		serviceLog.Error("enterMap response message UnMarshal failed")
		return
	}
	if respMsg.ErrorCode > 0 {
		serviceLog.Error("enterMap fail err: %s", respMsg.ErrorMessage)
		uc.Stop()
		return
	}

	payload := respMsg.GetEnterMapResponse()
	if payload == nil || payload.Me == nil || payload.Location == nil ||
		payload.Me.Position == nil || payload.Me.Dir == nil {
		serviceLog.Error("enterMapResponse msg payload data: %+v, is invalid", payload)
		return
	}

	env := &pubsubEventData.UserEnterGameEvent{
		MsgVersion:        time_helper.NowUTCMill(),
		SceneServiceAppId: uc.GetSceneService(),
		AgentAppId:        serviceCnf.GetInstance().ServerName,
		UserSocketId:      uc.GetSession().SessionId(),
		UserId:            payload.Me.BaseData.UserId,
		Name:              payload.Me.BaseData.Name,
		RoleId:            payload.Me.BaseData.RoleId,
		Gender:            payload.Me.BaseData.Gender,
		RoleIcon:          payload.Me.BaseData.RoleIcon,
		Guide:             payload.Me.BaseData.Guide,
		Eyebrow:           payload.Me.BaseData.Feature.Eyebrow,
		Mouth:             payload.Me.BaseData.Feature.Mouth,
		Eye:               payload.Me.BaseData.Feature.Eye,
		Face:              payload.Me.BaseData.Feature.Face,
		Hair:              payload.Me.BaseData.Feature.Hair,
		Glove:             payload.Me.BaseData.Feature.Glove,
		Clothes:           payload.Me.BaseData.Feature.Clothes,
		Pants:             payload.Me.BaseData.Feature.Pants,
		MapId:             payload.Location.MapId,
		X:                 payload.Me.Position.X,
		Y:                 payload.Me.Position.Y,
		Z:                 payload.Me.Position.Z,
		DirX:              payload.Me.Dir.X,
		DirY:              payload.Me.Dir.Y,
		DirZ:              payload.Me.Dir.Z,
	}

	uc.enterSceneService = true
	grpcPubsubEvent.RPCPubsubEventEnterGame(env)
}
