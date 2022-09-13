package userChannel

import "game-message-core/proto"

func (uc *UserChannel) agentClientMsg(msg *proto.Envelope) {
	switch msg.Type {
	case proto.EnvelopeType_Ping:
		uc.handlePing(msg)
	}
}

func (uc *UserChannel) handlePing(msg *proto.Envelope) {
	resMsg := &proto.Envelope{
		Type:  proto.EnvelopeType_Ping,
		SeqId: msg.SeqId,
		Payload: &proto.Envelope_PingResponse{
			PingResponse: &proto.PingResponse{},
		},
	}
	if byes, err := uc.MarshalProtoMessage(resMsg); err == nil {
		uc.tcpSession.Write(byes)
	}
}
