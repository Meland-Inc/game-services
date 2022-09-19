package msgChannel

import (
	"game-message-core/grpc/methodData"
	"game-message-core/proto"

	"github.com/Meland-Inc/game-services/src/services/main/msgChannel/clientMsgHandle"
)

type HandleFunc func(*methodData.PullClientMessageInput, *proto.Envelope)

func (ch *MsgChannel) registerHandler() {
	ch.msgHandler[proto.EnvelopeType_SigninPlayer] = clientMsgHandle.SingInHandle
	ch.msgHandler[proto.EnvelopeType_ItemGet] = clientMsgHandle.ItemGetHandle
	ch.msgHandler[proto.EnvelopeType_ItemUse] = clientMsgHandle.ItemUseHandle
	ch.msgHandler[proto.EnvelopeType_UpdateAvatar] = clientMsgHandle.LoadAvatarHandle
	ch.msgHandler[proto.EnvelopeType_UnloadAvatar] = clientMsgHandle.UnloadAvatarHandle
}
