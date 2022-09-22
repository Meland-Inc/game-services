package playerModel

import (
	"fmt"
	"game-message-core/proto"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/common/time_helper"
	"github.com/Meland-Inc/game-services/src/global/gameDB"
	dbData "github.com/Meland-Inc/game-services/src/global/gameDB/data"
	"github.com/Meland-Inc/game-services/src/global/grpcAPI/grpcInvoke"
	message "github.com/Meland-Inc/game-services/src/global/web3Message"
)

func (p *PlayerDataModel) GetPlayerItems(userId int64) (*PlayerItems, error) {
	cacheKey := p.playerItemsCacheKey(userId)
	iData, err := p.cache.GetOrStore(
		cacheKey,
		func() (interface{}, error) {
			items, err := p.LoadItems(userId)
			if err != nil {
				return nil, err
			}
			playerItems := &PlayerItems{
				UserId: userId,
				Items:  items,
			}
			return playerItems, err
		},
		p.cacheTTL,
	)
	if err != nil {
		return nil, err
	}

	p.cache.Touch(cacheKey, p.cacheTTL)
	return iData.(*PlayerItems), nil
}

func (p *PlayerDataModel) LoadItems(userId int64) (items []*Item, err error) {
	return []*Item{}, nil

	checkUsed := func(item *Item, usingNfts []dbData.UsingNft) (used bool, avatarPos int32) {
		for _, info := range usingNfts {
			if info.NftId == item.Id {
				return true, int32(info.AvatarPos)
			}
		}
		return false, 0
	}

	nfts, err := p.loadNFTS(userId)
	if err != nil {
		return nil, err
	}

	usingNfts, err := p.loadUsingNfts(userId)
	if err != nil {
		return nil, err
	}

	for _, nft := range nfts {
		used, avatarPos := checkUsed(nft, usingNfts)
		nft.Used = used
		nft.AvatarPos = avatarPos
		items = append(items, nft)
	}

	return items, err
}

func (p *PlayerDataModel) loadNFTS(userId int64) ([]*Item, error) {
	beginMs := time_helper.NowUTCMill()
	defer func() {
		serviceLog.Info("web3 load player[%d] NFT use time MS[%v]", userId, time_helper.NowMill()-beginMs)
	}()

	userNfts, err := grpcInvoke.RPCLoadUserNFTS(userId)
	if err != nil {
		serviceLog.Error("loadItemsByDapr User[%v] NFTS err : %+v", userId, err)
		return nil, err
	}

	items, err := p.parseUserNft(userId, userNfts)
	serviceLog.Info("user NFT list = %+v", userNfts)
	serviceLog.Info("user NFT list len(items)=%+v, err: %+v", len(items), err)
	return items, err
}

func (p *PlayerDataModel) parseUserNft(userId int64, userNfts *message.GetUserNFTsOutput) ([]*Item, error) {
	var items []*Item
	for _, nft := range userNfts.Nfts {
		item := NFTToItem(userId, nft)
		for _, out := range userNfts.PlaceableTimeouts {
			if out.NftId == item.Id {
				item.TimeOut = out
			}
		}
		items = append(items, item)
	}
	serviceLog.Info("user [%v] NFT item Length = %+v", userId, len(items))
	return items, nil
}

func (p *PlayerDataModel) loadUsingNfts(userId int64) ([]dbData.UsingNft, error) {
	if userId < 1 {
		return nil, fmt.Errorf("all using nft invalid user id [%d]", userId)
	}

	var usingNfts []dbData.UsingNft
	err := gameDB.GetGameDB().Where("user_id = ?", userId).Find(&usingNfts).Error
	return usingNfts, err
}

func (p *PlayerDataModel) ItemById(userId int64, nftId string) (*Item, error) {
	userItems, err := p.GetPlayerItems(userId)
	if err != nil {
		return nil, err
	}

	for _, it := range userItems.Items {
		if it.Id == nftId {
			return it, nil
		}
	}

	return nil, fmt.Errorf("Item not found")
}

func (p *PlayerDataModel) UsingAvatars(userId int64) (avatars []*Item, err error) {
	userItems, err := p.GetPlayerItems(userId)
	if err != nil {
		return nil, err
	}
	for _, it := range userItems.Items {
		if it.Used &&
			it.AvatarPos >= int32(proto.AvatarPosition_AvatarPositionHead) &&
			it.AvatarPos <= int32(proto.AvatarPosition_AvatarPositionWeapon) {
			avatars = append(avatars, it)
		}
	}
	return avatars, err
}

func (p *PlayerDataModel) addUsingNftRecord(item *Item) error {
	if item == nil {
		return fmt.Errorf("add using nft item is nil")
	}

	usingNft := &dbData.UsingNft{
		NftId:     item.Id,
		UserId:    item.Owner,
		Cid:       item.Cid,
		AvatarPos: item.AvatarPos,
	}
	return gameDB.GetGameDB().Create(usingNft).Error
}

func (p *PlayerDataModel) removeUsingNftRecord(userId int64, nftId string) error {
	if nftId == "" {
		return fmt.Errorf("delete using nft id is nil")
	}

	usingNft := dbData.UsingNft{}
	err := gameDB.GetGameDB().Where("nftId = ? ", nftId).First(&usingNft).Error
	if err != nil {
		return err
	}

	return gameDB.GetGameDB().Delete(&usingNft).Error
}

func (p *PlayerDataModel) canLoadAvatar(userId int64, item *Item, pos proto.AvatarPosition) error {
	if item.Attribute == nil {
		return fmt.Errorf("item [%s] attribute not found", item.Id)
	}
	if item.Attribute.Durability < 1 {
		return fmt.Errorf("item [%s] Durability is zero", item.Id)
	}
	slot, err := p.SlotByPosition(userId, pos)
	if err != nil {
		return err
	}
	if slot.Level < int(item.NFTData.UseLevel()) {
		return fmt.Errorf("item socket level < item need level")
	}
	return nil
}

// 穿装备
func (p *PlayerDataModel) LoadAvatar(userId int64, itemId string, pos proto.AvatarPosition) error {
	if pos < proto.AvatarPosition_AvatarPositionHead || pos > proto.AvatarPosition_AvatarPositionWeapon {
		return fmt.Errorf("invalid avatar position [%v]", pos)
	}

	item, err := p.ItemById(userId, itemId)
	if err != nil {
		return err
	}
	if err = p.canLoadAvatar(userId, item, pos); err != nil {
		return err
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
	if err = p.addUsingNftRecord(item); err != nil {
		return err
	}

	p.RPCCallUpdateUserUsingAvatar(userId)
	p.noticePlayerItemMsg(userId, proto.EnvelopeType_BroadCastItemUpdate, []*Item{item})
	return nil
}

// 卸装备
func (p *PlayerDataModel) UnloadAvatar(userId int64, itemId string, callProfileUp bool) error {
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
		p.RPCCallUpdateUserUsingAvatar(userId)
	}
	p.noticePlayerItemMsg(userId, proto.EnvelopeType_BroadCastItemUpdate, []*Item{item})
	return nil
}

func (p *PlayerDataModel) UseItem(userId int64, itemId string) error {
	it, err := p.ItemById(userId, itemId)
	if err != nil {
		return err
	}
	if err = p.canUse(userId, it); err != nil {
		return err
	}

	if err = p.callUseItem(userId, it); err != nil {
		return err
	}

	it.Used = true
	p.noticePlayerItemMsg(userId, proto.EnvelopeType_BroadCastItemUpdate, []*Item{it})
	return nil
}
func (p *PlayerDataModel) canUse(userId int64, it *Item) error {
	if it.Num < 1 {
		return fmt.Errorf("item is empty")
	}
	if it.Used {
		return fmt.Errorf("item is used")
	}

	player, err := p.GetPlayerSceneData(userId)
	if err != nil {
		return err
	}

	if player.Level < it.NFTData.UseLevel() {
		return fmt.Errorf("Insufficient level")
	}
	return nil
}
func (p *PlayerDataModel) callUseItem(userId int64, it *Item) error {
	switch it.NFTType {
	case proto.NFTType_NFTTypeConsumable:
		return p.callUseConsumable(userId, it)
	case proto.NFTType_NFTTypePlaceable, proto.NFTType_NFTTypeThird:
		// entities, err = m.useNFTBuild(userId, it)
	}

	return nil
}

func (p *PlayerDataModel) callUseConsumable(userId int64, item *Item) (err error) {
	isConsumable, conData := item.NFTData.GetConsumableData()
	if !isConsumable || conData == nil {
		return
	}

	err = grpcInvoke.RPCCallUseConsumableToWeb3(userId, item.Id, 0, 0)
	if err != nil {
		return err
	}
	return p.RPCEventUsedConsumable(userId, item)
}

func (p *PlayerDataModel) UpdatePlayerNFTs(userId int64, nfts []message.NFT) {
	needDelNfts := []message.NFT{}
	needUpNfts := []message.NFT{}
	for _, nft := range nfts {
		if nft.Amount == 0 {
			needDelNfts = append(needDelNfts, nft)
		} else {
			needUpNfts = append(needUpNfts, nft)
		}
	}

	delItems := []*Item{}
	for _, nft := range needDelNfts {
		if it := p.deleteNft(userId, nft.Id); it != nil {
			delItems = append(delItems, it)
		}
	}
	if len(delItems) > 0 {
		p.noticePlayerItemMsg(userId, proto.EnvelopeType_BroadCastItemDel, delItems)
	}

	addItems := []*Item{}
	upItems := []*Item{}
	for _, nft := range needUpNfts {
		item, _ := p.ItemById(userId, nft.Id)
		if item == nil {
			it := p.addNft(userId, nft)
			addItems = append(addItems, it)
		} else {
			upIt := p.updateNft(userId, item, nft)
			upItems = append(upItems, upIt)
		}
	}
	if len(upItems) > 0 {
		p.noticePlayerItemMsg(userId, proto.EnvelopeType_BroadCastItemUpdate, upItems)
	}
	if len(addItems) > 0 {
		p.noticePlayerItemMsg(userId, proto.EnvelopeType_BroadCastItemAdd, addItems)
	}
}

func (p *PlayerDataModel) addNft(userId int64, nft message.NFT) *Item {
	item := NFTToItem(userId, nft)
	playerItems, _ := p.GetPlayerItems(userId)
	playerItems.AddItem(item)
	serviceLog.Info("add new  NFT item = %+v", item)
	return item
}

func (p *PlayerDataModel) updateNft(userId int64, item *Item, nft message.NFT) *Item {
	item.Num = int32(nft.Amount)
	item.NFTData = nft
	switch item.NFTType {
	case proto.NFTType_NFTTypeEquipment:
		if _, _, attr := nft.GetEquipmentData(); attr != nil {
			item.Attribute = attr
		}
	case proto.NFTType_NFTTypeWearable:
		if _, _, attr := nft.GetWearablePbData(); attr != nil {
			item.Attribute = attr
		}
	}
	serviceLog.Info("update NFT item = %+v", item)
	return item
}

func (p *PlayerDataModel) deleteNft(userId int64, nftId string) *Item {
	item, err := p.ItemById(userId, nftId)
	if err != nil {
		return nil
	}

	playerItems, _ := p.GetPlayerItems(userId)
	playerItems.DelItem(nftId)
	if item.Used && item.AvatarPos > 0 {
		p.removeUsingNftRecord(userId, item.Id)
	}
	serviceLog.Info("delete NFT item = %+v", item)
	return item
}
