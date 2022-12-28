package playerModel

import (
	"fmt"
	"game-message-core/proto"

	"github.com/Meland-Inc/game-services/src/common/matrix"
	"github.com/Meland-Inc/game-services/src/global/configData"
	dbData "github.com/Meland-Inc/game-services/src/global/gameDB/data"
	"github.com/Meland-Inc/game-services/src/global/grpcAPI/grpcPubsubEvent"
)

func (p *PlayerDataModel) canUpgradeLevel(player *dbData.PlayerSceneData) error {
	maxLvSetting, err := configData.GameValueById(1000001)
	if err != nil {
		return err
	}
	if player.Level >= maxLvSetting.Value {
		return fmt.Errorf("player is max level")
	}
	curLvCnf := configData.ConfigMgr().RoleLevelCnf(player.Level)
	if curLvCnf == nil {
		return fmt.Errorf("role level[%d] not found", player.Level)
	}

	if player.Exp < curLvCnf.Exp {
		return fmt.Errorf("current Exp not upgreade to next level")
	}

	playerItemSlot, err := p.GetPlayerItemSlots(player.UserId)
	if err != nil {
		return err
	}

	// 角色升级要求是有4个或以上的插槽等级大于角色当前等级-5
	slotLvOffsetSetting, err := configData.GameValueById(1000004)
	if err != nil {
		return err
	}
	slotLvOffsetCountSetting, err := configData.GameValueById(1000005)
	if err != nil {
		return err
	}

	var slotLvOffsetCount int32
	for _, s := range playerItemSlot.GetSlotList().SlotList {
		if int32(s.Level) > player.Level-slotLvOffsetSetting.Value {
			slotLvOffsetCount++
		}
	}

	if slotLvOffsetCount < slotLvOffsetCountSetting.Value {
		return fmt.Errorf("can used item socket level < 4")
	}

	return nil
}

func (p *PlayerDataModel) UpgradePlayerLevel(userId int64) (lv int32, exp int32, err error) {
	player, err := p.GetPlayerSceneData(userId)
	if err != nil {
		return 0, 0, err
	}
	if err = p.canUpgradeLevel(player); err != nil {
		return player.Level, player.Exp, err
	}

	curLv := player.Level
	curExp := player.Exp
	curLvCnf := configData.ConfigMgr().RoleLevelCnf(curLv)
	if curLvCnf == nil {
		return curLv, curExp, fmt.Errorf("role level[%d] not found", player.Level)
	}

	newExp := player.Exp - curLvCnf.Exp
	newLv := player.Level + 1
	err = p.setLevelAndExp(userId, newLv, newExp)
	if err != nil {
		return curLv, curExp, err
	}

	grpcPubsubEvent.RPCPubsubEventUserLevelUpgrade(userId, newLv)
	p.RPCCallUpdateUserProfile(userId)
	return player.Level, player.Exp, nil
}

func (p *PlayerDataModel) setLevelAndExp(userId int64, lv, exp int32) error {
	if exp < 0 || lv < 0 {
		return fmt.Errorf("level [%d] exp[%v] invalid", lv, exp)
	}
	player, err := p.GetPlayerSceneData(userId)
	if err != nil {
		return err
	}
	if player.Level == lv && player.Exp == exp {
		return nil
	}

	upProfiles := []*proto.EntityProfileUpdate{}
	if player.Level != lv {
		player.Level = lv
		upProfiles = append(upProfiles, &proto.EntityProfileUpdate{
			Field:    proto.EntityProfileField_EntityProfileFieldLv,
			CurValue: lv,
		})
	}
	if player.Exp != exp {
		player.Exp = matrix.LimitInt32(exp, 0, configData.RoleCurrentExpLimit())
		upProfiles = append(upProfiles, &proto.EntityProfileUpdate{
			Field:    proto.EntityProfileField_EntityProfileFieldExp,
			CurValue: exp,
		})
	}
	if err = p.UpPlayerSceneData(player); err != nil {
		return err
	}
	if len(upProfiles) > 0 {
		p.noticePlayerProfileUpdate(userId, upProfiles)
	}
	return nil
}

func (p *PlayerDataModel) AddExp(userId int64, exp int32) error {
	if exp < 1 {
		return fmt.Errorf("invalid add exp [%d]", exp)
	}
	player, err := p.GetPlayerSceneData(userId)
	if err != nil {
		return err
	}
	return p.setLevelAndExp(userId, player.Level, player.Exp+exp)
}

func (p *PlayerDataModel) DeductExp(userId int64, deductExp int32) error {
	if deductExp < 1 {
		return fmt.Errorf("invalid Deduct exp [%d]", deductExp)
	}
	sceneData, err := p.GetPlayerSceneData(userId)
	if err != nil {
		return err
	}

	if sceneData.Exp < deductExp {
		return fmt.Errorf("curExp[%d] cannot deductExp[%d]", sceneData.Exp, deductExp)
	}

	return p.setLevelAndExp(userId, sceneData.Level, sceneData.Exp-deductExp)
}
