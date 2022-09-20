package playerModel

import (
	"game-message-core/grpc/pubsubEventData"
	"game-message-core/proto"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/common/time_helper"
	"github.com/Meland-Inc/game-services/src/global/grpcAPI/grpcInvoke"
	"github.com/Meland-Inc/game-services/src/global/grpcAPI/grpcPubsubEvent"
)

func (p *PlayerDataModel) RPCEventUsedConsumable(userId int64, item *Item) error {
	_, conData := item.NFTData.GetConsumableData()

	env := pubsubEventData.UserUseNFTEvent{
		MsgVersion:     time_helper.NowUTCMill(),
		UserId:         userId,
		NftId:          item.Id,
		Cid:            item.Cid,
		NftType:        item.NFTType,
		Num:            1,
		ConsumableData: conData,
		X:              0,
		Z:              0,
	}
	return grpcPubsubEvent.RPCPubsubEventUseNft(env)
}

func (p *PlayerDataModel) RPCCallUpdateUserUsingAvatar(userId int64) {
	profile, err := p.GetPlayerProfile(userId)
	if err != nil {
		serviceLog.Error("call UpdateUsedAvatar get profile failed err: %v", err)
		return
	}
	items, err := p.UsingAvatars(userId)
	if err != nil {
		serviceLog.Error("call UpdateUsedAvatar get avatar failed err: %v", err)
		return
	}

	avatars := []proto.PlayerAvatar{}
	for _, it := range items {
		avatar := it.ToNetPlayerAvatar()
		avatars = append(avatars, *avatar)
	}
	err = grpcInvoke.UpdateUsedAvatar(userId, avatars, *profile)
	if err != nil {
		serviceLog.Error("call UpdateUsedAvatar failed err: %v", err)
		return
	}
}

func (p *PlayerDataModel) RPCCallUpdateUserProfile(userId int64) {
	profile, err := p.GetPlayerProfile(userId)
	if err != nil {
		serviceLog.Error("call UpdateUsedProfile get profile failed err: %v", err)
		return
	}
	err = grpcInvoke.UpdateUsedProfile(userId, *profile)
	if err != nil {
		serviceLog.Error("call UpdateUsedAvatar failed err: %v", err)
		return
	}
}
