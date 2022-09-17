package itemModel

import (
	"fmt"
	"game-message-core/proto"
)

// 穿装备
func (p *ItemModel) LoadAvatar(userId int64, userLv int32, itemId string, pos proto.AvatarPosition) error {
	if pos < proto.AvatarPosition_AvatarPositionHead || pos > proto.AvatarPosition_AvatarPositionWeapon {
		return fmt.Errorf("invalid avatar position [%v]", pos)
	}

	item, err := p.ItemById(userId, itemId)
	if err != nil {
		return err
	}
	if item.Attribute == nil {
		return fmt.Errorf("item [%s] attribute not found", itemId)
	}
	if item.Attribute.Durability < 1 {
		return fmt.Errorf("item [%s] Durability is zero", itemId)
	}

	itemSlotLv := 1 // TODO: DATA FROM ITEM SLOT MODEL
	if itemSlotLv < int(item.NFTData.UseLevel()) {
		return fmt.Errorf("item socket level < item need level")
	}

	// 检查目标avatar POS 是否有装备正在使用, 有就先卸下
	usingAvatars, _ := p.UsingAvatars(userId)
	for _, it := range usingAvatars {
		if it.Used && it.AvatarPos == int32(pos) {
			p.UnloadAvatar(userId, it.Id, false)
			break
		}
	}
	// 使用装备
	item.Used = true
	item.AvatarPos = int32(pos)
	err = p.addUsingNftRecord(item)
	if err != nil {
		return err
	}

	// TODO CLIENT PLAYER PROFILE AND call scene service up player profile
	p.NoticePlayer(userId, proto.EnvelopeType_BroadCastItemUpdate, []*Item{item})
	return nil
}

// 卸装备
func (p *ItemModel) UnloadAvatar(userId int64, itemId string, callProfileUp bool) error {
	item, err := p.ItemById(userId, itemId)
	if err != nil {
		return err
	}

	item.Used = false
	item.AvatarPos = int32(proto.AvatarPosition_AvatarPositionNone)
	if err := p.removeUsingNftRecord(userId, item.Id); err != nil {
		return err
	}

	if callProfileUp {
		// TODO CLIENT PLAYER PROFILE AND call scene service up player profile
	}
	p.NoticePlayer(userId, proto.EnvelopeType_BroadCastItemUpdate, []*Item{item})
	return nil
}
