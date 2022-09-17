package itemModel

import (
	"encoding/json"
	"game-message-core/proto"

	message "github.com/Meland-Inc/game-services/src/global/web3Message"
	"github.com/spf13/cast"
)

// db table name player_items

type Item struct {
	Id        string                      `db:"id"`
	Owner     int64                       `db:"owner"`
	Cid       int32                       `db:"cid"`
	Num       int32                       `db:"num"`
	AvatarPos int32                       `db:"avatar_pos"`
	Used      bool                        `db:"-"`
	ItemType  proto.ItemType              `db:"-"`
	Attribute *proto.AvatarAttribute      `db:"-"`
	NFTType   proto.NFTType               `db:"-"`
	NFTData   message.NFT                 `db:"-"`
	TimeOut   message.NFTPlaceableTimeout `db:"-"`
}

func (it *Item) ToNetItem() *proto.Item {
	pbIt := &proto.Item{
		ItemType:      it.ItemType,
		Id:            it.Id,
		ObjectCid:     it.Cid,
		Num:           it.Num,
		UserId:        it.Owner,
		Attribute:     it.Attribute,
		AvatarPos:     proto.AvatarPosition(it.AvatarPos),
		NftUsing:      it.Used,
		NftTimeOutSec: int32(it.TimeOut.TimeoutSec),
	}
	if bs, err := json.Marshal(it.NFTData); err == nil {
		pbIt.NftJsonData = string(bs)
	}
	return pbIt
}

type PlayerItems struct {
	UserId int64
	Items  []*Item
}

func NFTToItem(userId int64, nft message.NFT) *Item {
	item := &Item{
		Id:       nft.Id,
		Owner:    userId,
		Cid:      cast.ToInt32(nft.ItemId),
		Num:      cast.ToInt32(nft.Amount),
		ItemType: proto.ItemType_ItemTypeNFT,
		NFTType:  message.NFTPbType(nft),
		NFTData:  nft,
	}

	switch item.NFTType {
	case proto.NFTType_NFTTypeEquipment:
		if _, _, attr := nft.GetEquipmentData(); attr != nil {
			item.Attribute = attr
		}
	case proto.NFTType_NFTTypeWearable:
		if _, _, attr := nft.GetWearablePbData(); attr != nil {
			item.Attribute = attr
		}
	default:
	}
	return item
}
