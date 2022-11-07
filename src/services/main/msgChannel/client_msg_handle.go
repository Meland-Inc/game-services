package msgChannel

import (
	"game-message-core/grpc/methodData"
	"game-message-core/proto"

	"github.com/Meland-Inc/game-services/src/services/main/msgChannel/clientMsgHandle"
)

type HandleFunc func(*methodData.PullClientMessageInput, *proto.Envelope)

func (ch *MsgChannel) registerClientMsgHandler() {
	ch.clientMsgHandler[proto.EnvelopeType_SigninPlayer] = clientMsgHandle.SingInHandler
	ch.clientMsgHandler[proto.EnvelopeType_ItemGet] = clientMsgHandle.ItemGetHandle
	ch.clientMsgHandler[proto.EnvelopeType_ItemUse] = clientMsgHandle.ItemUseHandle
	ch.clientMsgHandler[proto.EnvelopeType_UpdateAvatar] = clientMsgHandle.LoadAvatarHandle
	ch.clientMsgHandler[proto.EnvelopeType_UnloadAvatar] = clientMsgHandle.UnloadAvatarHandle
	ch.clientMsgHandler[proto.EnvelopeType_GetItemSlot] = clientMsgHandle.ItemSlotGetHandle
	ch.clientMsgHandler[proto.EnvelopeType_UpgradeItemSlot] = clientMsgHandle.ItemSlotUpgradeHandle
	ch.clientMsgHandler[proto.EnvelopeType_UpgradePlayerLevel] = clientMsgHandle.UpgradePlayerLevelHandle

	// land and build client msg handles
	ch.clientMsgHandler[proto.EnvelopeType_QueryLands] = clientMsgHandle.QueryLandsHandler
	ch.clientMsgHandler[proto.EnvelopeType_Build] = clientMsgHandle.BuildHandler
	ch.clientMsgHandler[proto.EnvelopeType_Recycling] = clientMsgHandle.RecyclingHandler
	ch.clientMsgHandler[proto.EnvelopeType_MintBattery] = clientMsgHandle.MintBatteryHandler
	ch.clientMsgHandler[proto.EnvelopeType_Charged] = clientMsgHandle.ChargedHandler
	ch.clientMsgHandler[proto.EnvelopeType_Harvest] = clientMsgHandle.HarvestHandler
	ch.clientMsgHandler[proto.EnvelopeType_Collection] = clientMsgHandle.CollectionHandler
	ch.clientMsgHandler[proto.EnvelopeType_SelfNftBuilds] = clientMsgHandle.SelfNftBuildsHandler
}
