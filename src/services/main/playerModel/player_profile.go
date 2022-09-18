package playerModel

import (
	"fmt"
	"game-message-core/proto"

	"github.com/Meland-Inc/game-services/src/global/configData"
	dbData "github.com/Meland-Inc/game-services/src/global/gameDB/data"
	message "github.com/Meland-Inc/game-services/src/global/web3Message"
)

func (p *PlayerDataModel) GetPlayerProfile(userId int64) (*proto.EntityProfile, error) {
	sceneData, err := p.GetPlayerSceneData(userId)
	if err != nil {
		return nil, err
	}
	avatars, err := p.UsingAvatars(userId)
	if err != nil {
		return nil, err
	}
	itemSlot, err := p.GetPlayerItemSlots(userId)
	if err != nil {
		return nil, err
	}
	return calculatePlayerProfile(sceneData, avatars, itemSlot.Slots.SlotList)
}

// 计算玩家属性(等级,装备...)<清空之前的数据全部从新计算>
func calculatePlayerProfile(
	sceneData *dbData.PlayerSceneData,
	avatars []*Item,
	slotList []message.PlayerItemSlot,
) (*proto.EntityProfile, error) {
	pro := &proto.EntityProfile{
		Lv:        sceneData.Level,
		Exp:       int64(sceneData.Exp),
		HpCurrent: sceneData.Hp,
	}
	if err := profileAddByLv(pro, sceneData.Level); err != nil {
		return pro, err
	}
	profileAddByAvatar(pro, avatars)
	profileByItemSlotLv(pro, slotList)
	return pro, nil
}

func profileAddByLv(pro *proto.EntityProfile, lv int32) error {
	lvSetting := configData.ConfigMgr().RoleLevelCnf(lv)
	if lvSetting == nil {
		return fmt.Errorf("player lv[%d] config not found", lv)
	}

	pro.Lv = lvSetting.Lv
	pro.HpLimit = lvSetting.HpLimit
	pro.HpRecovery = lvSetting.HpRecovery
	pro.Att = lvSetting.Att
	pro.AttSpeed = lvSetting.AttSpeed
	pro.Def = lvSetting.Def
	pro.CritRate = lvSetting.CritRate
	pro.CritDmg = lvSetting.CritDmg
	pro.HitRate = lvSetting.HitRate
	pro.MissRate = lvSetting.MissRate
	pro.MoveSpeed = float32(lvSetting.MoveSpeed)
	return nil
}

func profileAddByAvatar(pro *proto.EntityProfile, avatars []*Item) {
	for _, it := range avatars {
		if it == nil || it.Attribute == nil {
			continue
		}
		// 耐久度=0 装备属性无效
		if it.Attribute.Durability <= 0 {
			continue
		}

		for _, d := range it.Attribute.Data {
			switch d.Type {
			case proto.AttributeType_AttributeTypeHpLimit: // hp limit
				pro.HpLimit += d.Value
			case proto.AttributeType_AttributeTypeHpRecovery: // hp recovery
				pro.HpRecovery += d.Value
			case proto.AttributeType_AttributeTypeAtt: // 增加普通攻击
				pro.Att += d.Value
			case proto.AttributeType_AttributeTypeAttSpeed: // 增加普通攻击speed
				pro.AttSpeed += d.Value
			case proto.AttributeType_AttributeTypeDef: // 增加普通防御
				pro.Def += d.Value
			case proto.AttributeType_AttributeTypeCrit: // 增加暴击率
				pro.CritRate += d.Value
			case proto.AttributeType_AttributeTypeCritDmg: //
				pro.CritDmg += d.Value
			case proto.AttributeType_AttributeTypeHitRate: // 增加暴击伤害
				pro.HitRate += d.Value
			case proto.AttributeType_AttributeTypeMissRate: // 增加闪避率
				pro.MissRate += d.Value
			case proto.AttributeType_AttributeTypeMoveSpeed: // 增加移动速度
				pro.MoveSpeed += float32(d.Value)
			}
		}
	}
}

func profileByItemSlotLv(pro *proto.EntityProfile, slotList []message.PlayerItemSlot) {
	for _, s := range slotList {
		setting := configData.ConfigMgr().GetSlotCnf(int32(s.Position), int32(s.Level))
		if setting == nil {
			fmt.Errorf("slot position[%v]lv[%d] config not found", s.Position, s.Level)
			continue
		}

		pro.HpLimit += setting.HpLimit
		pro.HpRecovery += setting.HpRecovery
		pro.Att += setting.Att
		pro.AttSpeed += setting.AttSpeed
		pro.Def += setting.Def
		pro.CritRate += setting.CritRate
		pro.CritDmg += setting.CritDmg
		pro.HitRate += setting.HitRate
		pro.MissRate += setting.MissRate
		pro.MoveSpeed += float32(setting.MoveSpeed)
	}
}
