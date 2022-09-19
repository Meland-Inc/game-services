package msgChannel

import (
	"game-message-core/grpc/methodData"
	"game-message-core/proto"

	"github.com/Meland-Inc/game-services/src/services/main/msgChannel/clientMsgHandle"
)

type HandleFunc func(*methodData.PullClientMessageInput, *proto.Envelope)

func (ch *MsgChannel) registerClientMsgHandler() {
	ch.clientMsgHandler[proto.EnvelopeType_SigninPlayer] = clientMsgHandle.SingInHandle
	ch.clientMsgHandler[proto.EnvelopeType_ItemGet] = clientMsgHandle.ItemGetHandle
	ch.clientMsgHandler[proto.EnvelopeType_ItemUse] = clientMsgHandle.ItemUseHandle
	ch.clientMsgHandler[proto.EnvelopeType_UpdateAvatar] = clientMsgHandle.LoadAvatarHandle
	ch.clientMsgHandler[proto.EnvelopeType_UnloadAvatar] = clientMsgHandle.UnloadAvatarHandle
	ch.clientMsgHandler[proto.EnvelopeType_GetItemSlot] = clientMsgHandle.ItemSlotGetHandle
	ch.clientMsgHandler[proto.EnvelopeType_UpgradeItemSlot] = clientMsgHandle.ItemSlotUpgradeHandle
	ch.clientMsgHandler[proto.EnvelopeType_UpgradePlayerLevel] = clientMsgHandle.UpgradePlayerLevelHandle
}
