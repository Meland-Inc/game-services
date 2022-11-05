package playerModel

import (
	"fmt"
	"game-message-core/proto"
	xlsxTable "game-message-core/xlsxTableData"

	"github.com/Meland-Inc/game-services/src/common/time_helper"
	"github.com/Meland-Inc/game-services/src/global/configData"
	"github.com/Meland-Inc/game-services/src/global/gameDB"
	dbData "github.com/Meland-Inc/game-services/src/global/gameDB/data"
	"github.com/Meland-Inc/game-services/src/global/grpcAPI/grpcInvoke"
	"github.com/Meland-Inc/game-services/src/global/grpcAPI/grpcPubsubEvent"
	message "github.com/Meland-Inc/game-services/src/global/web3Message"
	"gorm.io/gorm"
)

func (p *PlayerDataModel) initPlayerItemSlot(userId int64) (*dbData.ItemSlot, error) {
	slot := &dbData.ItemSlot{
		UserId:    userId,
		CreatedAt: time_helper.NowUTC(),
		UpdateAt:  time_helper.NowUTC(),
	}
	err := slot.InitSlotList()
	return slot, err
}

func (p *PlayerDataModel) GetPlayerItemSlots(userId int64) (*dbData.ItemSlot, error) {
	cacheKey := p.playerItemsSlotCacheKey(userId)
	iData, err := p.cache.GetOrStore(
		cacheKey,
		func() (interface{}, error) {
			playerSlot := &dbData.ItemSlot{}
			err := gameDB.GetGameDB().Where("user_id = ?", userId).First(playerSlot).Error
			if err != nil && err == gorm.ErrRecordNotFound {
				playerSlot, err = p.initPlayerItemSlot(userId)
				err = nil
			}
			return playerSlot, err
		},
		p.cacheTTL,
	)
	if err != nil {
		return nil, err
	}
	p.cache.Touch(cacheKey, p.cacheTTL)
	return iData.(*dbData.ItemSlot), nil
}

func (p *PlayerDataModel) SlotByPosition(userId int64, pos proto.AvatarPosition) (*message.PlayerItemSlot, error) {
	userSlot, err := p.GetPlayerItemSlots(userId)
	if err != nil {
		return nil, err
	}
	for _, s := range userSlot.GetSlotList().SlotList {
		if s.Position == int(pos) {
			return s, nil
		}
	}
	return nil, fmt.Errorf("position [%v] slot not found", pos)
}

func (p *PlayerDataModel) setPlayerItemSlotLevel(
	userId int64, pos proto.AvatarPosition, lv int32, broadcast bool,
) (*dbData.ItemSlot, error) {
	if pos < proto.AvatarPosition_AvatarPositionHead || pos > proto.AvatarPosition_AvatarPositionWeapon {
		return nil, fmt.Errorf("invalid slot avatar position")
	}
	playerSlot, err := p.GetPlayerItemSlots(userId)
	if err != nil {
		return nil, err
	}
	playerSlot.SetSlotLevel(pos, lv)
	if err := gameDB.GetGameDB().Save(playerSlot).Error; err != nil {
		return playerSlot, err
	}
	p.RPCCallUpdateUserProfile(userId)
	if broadcast {
		p.noticeUpdatePlayerItemSlot(playerSlot)
	}
	return playerSlot, nil
}

func (p *PlayerDataModel) UpgradeItemSlots(
	userId int64, pos proto.AvatarPosition, broadcast bool,
) (*dbData.ItemSlot, error) {
	userSlot, err := p.GetPlayerItemSlots(userId)
	if err != nil {
		return nil, err
	}

	player, err := p.GetPlayerSceneData(userId)
	if err != nil {
		return userSlot, err
	}

	var curSocket *message.PlayerItemSlot
	for _, s := range userSlot.GetSlotList().SlotList {
		if proto.AvatarPosition(s.Position) == pos {
			curSocket = s
			break
		}
	}

	err = p.canUpgradeItemSlots(player, pos, curSocket.Level)
	if err != nil {
		return nil, err
	}

	// Deduction upgrade item socket need used player exp
	setting := configData.ConfigMgr().GetSlotCnf(int32(pos), player.Level)
	p.setLevelAndExp(userId, player.Level, player.Exp-setting.UpExp)

	newLevel := int32(curSocket.Level + 1)
	_, err = p.setPlayerItemSlotLevel(userId, pos, newLevel, broadcast)
	if err == nil {
		grpcPubsubEvent.RPCPubsubEventSlotLevelUpgrade(userId, int32(pos), newLevel)
	}
	return userSlot, err
}

func (p *PlayerDataModel) canUpgradeItemSlots(player *dbData.PlayerSceneData, pos proto.AvatarPosition, curLv int) error {
	slotMaxLvSetting, err := configData.GameValueById(1000002)
	if err != nil {
		return err
	}
	if int32(curLv) >= slotMaxLvSetting.Value {
		return fmt.Errorf("slot is max level")
	}

	slotLvOffsetSetting, err := configData.GameValueById(1000004)
	if err != nil {
		return err
	}
	if int32(curLv) >= player.Level+slotLvOffsetSetting.Value {
		return fmt.Errorf("item slot position [%v] is current max level", pos)
	}

	setting := configData.ConfigMgr().GetSlotCnf(int32(pos), int32(curLv))
	if setting == nil {
		return fmt.Errorf("item slot position:[%v] Lv:[%d] config not found", pos, curLv)
	}

	if player.Exp < setting.UpExp {
		return fmt.Errorf("cur exp can't up to next level")
	}

	// burn upgrade item socket need used meld
	if setting.UpMeld > 0 {
		return grpcInvoke.BurnUserMELD(player.UserId, int(setting.UpMeld))
	}
	return nil
}

func (p *PlayerDataModel) allUsingItemSlotAttributes(userId int64) (allSettings []*xlsxTable.SlotLvTableRow) {
	userSlot, err := p.GetPlayerItemSlots(userId)
	if err != nil {
		return
	}

	configMgr := configData.ConfigMgr()
	for _, s := range userSlot.GetSlotList().SlotList {
		setting := configMgr.GetSlotCnf(int32(s.Position), int32(s.Level))
		if setting != nil {
			allSettings = append(allSettings, setting)
		}
	}

	return allSettings
}

func (p *PlayerDataModel) playerMaxLevelItemSlot(userId int64) (proto.ItemSlot, error) {
	userSlot, err := p.GetPlayerItemSlots(userId)
	if err != nil {
		return proto.ItemSlot{}, err
	}

	levelMaxSlot := proto.ItemSlot{}
	for _, s := range userSlot.GetSlotList().SlotList {
		if levelMaxSlot.Level < int32(s.Level) {
			levelMaxSlot.Level = int32(s.Level)
			levelMaxSlot.Position = proto.AvatarPosition(s.Position)
		}
	}
	return levelMaxSlot, err
}
