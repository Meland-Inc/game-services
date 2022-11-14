package playerModel

import (
	"game-message-core/proto"

	message "github.com/Meland-Inc/game-services/src/global/web3Message"
	"github.com/spf13/cast"
)

// db table name player_items

type Item struct {
	Id           string                      `db:"id"`
	Owner        int64                       `db:"owner"`
	Cid          int32                       `db:"cid"`
	Num          int32                       `db:"num"`
	AvatarPos    int32                       `db:"avatar_pos"`
	Used         bool                        `db:"-"`
	ItemType     proto.ItemType              `db:"-"`
	Attribute    *proto.AvatarAttribute      `db:"-"`
	NFTType      proto.NFTType               `db:"-"`
	NFTData      message.NFT                 `db:"-"`
	TimeOut      message.NFTPlaceableTimeout `db:"-"`
	ProtoNftData *proto.NftData              `db:"-"`
}

func (it *Item) ToNetItem() *proto.Item {
	pbIt := &proto.Item{
		ItemType:      it.ItemType,
		Id:            it.Id,
		ObjectCid:     it.Cid,
		Num:           it.Num,
		UserId:        it.Owner,
		AvatarPos:     proto.AvatarPosition(it.AvatarPos),
		NftUsing:      it.Used,
		NftTimeOutSec: int32(it.TimeOut.TimeoutSec),
		NftData:       it.ProtoNftData,
	}
	return pbIt
}

func (it *Item) ToNetPlayerAvatar() *proto.PlayerAvatar {
	return &proto.PlayerAvatar{
		Position:  proto.AvatarPosition(it.AvatarPos),
		ObjectId:  it.Cid,
		Attribute: it.Attribute,
	}
}

type PlayerItems struct {
	UserId int64
	Items  []*Item
}

func (p *PlayerItems) AddItem(item *Item) {
	p.Items = append(p.Items, item)
}

func (p *PlayerItems) DelItem(itemId string) {
	for idx, item := range p.Items {
		if item.Id == itemId {
			p.Items = append(p.Items[:idx], p.Items[idx+1:]...)
			break
		}
	}
}

func NFTToItem(userId int64, nft message.NFT) *Item {
	item := &Item{
		Id:           nft.Id,
		Owner:        userId,
		Cid:          cast.ToInt32(nft.ItemId),
		Num:          cast.ToInt32(nft.Amount),
		ItemType:     proto.ItemType_ItemTypeNFT,
		NFTType:      message.NFTPbType(nft),
		NFTData:      nft,
		ProtoNftData: message.ToProtoNftData(nft),
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
