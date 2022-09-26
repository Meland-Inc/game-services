package msgChannel

import (
	"game-message-core/grpc/methodData"
	"game-message-core/proto"

	"github.com/Meland-Inc/game-services/src/services/task/msgChannel/clientMsgHandle"
)

type HandleFunc func(*methodData.PullClientMessageInput, *proto.Envelope)

func (ch *MsgChannel) registerClientMsgHandler() {
	ch.clientMsgHandler[proto.EnvelopeType_SelfTasks] = clientMsgHandle.SelfTasksHandler
}
