package userChannel

import (
	"game-message-core/proto"
	"game-message-core/protoTool"
)

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
	msgBody, _ := protoTool.MarshalEnvelope(resMsg)
	uc.SendToUser(msg.Type, msgBody)
}
