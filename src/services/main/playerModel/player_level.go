package playerModel

import (
	"fmt"
	"game-message-core/proto"

	"github.com/Meland-Inc/game-services/src/common/matrix"
	"github.com/Meland-Inc/game-services/src/global/configData"
	dbData "github.com/Meland-Inc/game-services/src/global/gameDB/data"
)

func (p *PlayerDataModel) canUpgradeLevel(player *dbData.PlayerSceneData) error {
	maxLv := configData.ConfigMgr().RoleMaxLevel()
	if player.Level == maxLv {
		return fmt.Errorf("player is max level")
	}
	curLvCnf := configData.ConfigMgr().RoleLevelCnf(player.Level)
	if curLvCnf == nil {
		return fmt.Errorf("role level[%d] not found", player.Level)
	}

	if player.Exp < curLvCnf.Exp {
		return fmt.Errorf("current Exp not upgreade to next level")
	}

	// playerItemSlot, err := p.getPlayerItemSlots(player.UserId)
	// if err != nil {
	// 	return err
	// }

	// 角色升级要求是有4个或以上的插槽等级大于角色当前等级-5
	// var count int32
	// for _, s := range playerItemSlot.ItemSlots {
	// 	if s.Level > player.Lv-5 {
	// 		count++
	// 	}
	// }
	// if count < 4 {
	// 	return fmt.Errorf("can used item socket level < 4")
	// }

	return nil
}

func (p *PlayerDataModel) UpgradeLevel(userId int64) (lv int32, exp int32, err error) {
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

	// TODO: call scene service update player profile

	return player.Level, player.Exp, err
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

	exp = matrix.LimitInt32(exp, 0, configData.ConfigMgr().RoleCurrentExpLimit())
	if err = p.UpPlayerSceneData(
		userId, player.Hp, lv, exp, player.MapId, player.X,
		player.Y, player.Z, player.DirX, player.DirY, player.DirZ,
	); err != nil {
		return err
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
		player.Exp = exp
		upProfiles = append(upProfiles, &proto.EntityProfileUpdate{
			Field:    proto.EntityProfileField_EntityProfileFieldExp,
			CurValue: exp,
		})
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

func (p *PlayerDataModel) DeductExp(userId int64, exp int32) error {
	if exp < 1 {
		return fmt.Errorf("invalid Deduct exp [%d]", exp)
	}
	player, err := p.GetPlayerSceneData(userId)
	if err != nil {
		return err
	}
	return p.setLevelAndExp(userId, player.Level, player.Exp-exp)
}
