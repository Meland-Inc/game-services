package msgChannel

import (
	"game-message-core/grpc/methodData"
	"game-message-core/proto"

	"github.com/Meland-Inc/game-services/src/services/chat/msgChannel/clientMsgHandle"
)

type HandleFunc func(*methodData.PullClientMessageInput, *proto.Envelope)

func (ch *MsgChannel) registerClientMsgHandler() {
	ch.clientMsgHandler[proto.EnvelopeType_SendChatMessage] = clientMsgHandle.ChatMsgHandle
}
