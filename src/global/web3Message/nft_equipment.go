package message

import (
	"strconv"

	"game-message-core/proto"

	"github.com/spf13/cast"
)

// ----------------- 装备 --------------------------------------------------

func (n *NFT) isEquipment(value string) (b bool) {
	switch value {
	case string(NFTTraitTypeHandsArmor): // "Hands Armor" 头部装备
		b = true
	case string(NFTTraitTypeChestArmor): // "Chest Armor" 胸部装备
		b = true
	case string(NFTTraitTypeHeadArmor): // "Head Armor" 手部装备
		b = true
	case string(NFTTraitTypeLegsArmor): // "Legs Armor" 腿部装备
		b = true
	case string(NFTTraitTypeFeetArmor): // "Feet Armor" 脚部装备
		b = true
	case string(NFTTraitTypeSword): // "Sword" 剑
		b = true
	case string(NFTTraitTypeBow): // "Bow"  弓
		b = true
	case string(NFTTraitTypeDagger): // "Dagger" 匕首
		b = true
	case string(NFTTraitTypeSpear): // "Spear"枪
		b = true

		// case string(NFTTraitTypeConsumable): // "Consumable" 消耗品
		// case string(NFTTraitTypeMaterial): // "Material" 材料
		// case string(NFTTraitTypeMysteryBox): // "MysteryBox" 神秘宝箱
		// case string(NFTTraitTypePlaceable): // "Placeable" 可放置
	}
	return b
}

func (n *NFT) equipmentPosition(value string) (position proto.AvatarPosition) {
	switch value {
	case string(NFTTraitTypeHandsArmor): // "Hands Armor" 头部装备
		position = proto.AvatarPosition_AvatarPositionHead
	case string(NFTTraitTypeChestArmor): // "Chest Armor" 胸部装备
		position = proto.AvatarPosition_AvatarPositionCoat
	case string(NFTTraitTypeHeadArmor): // "Head Armor" 手部装备
		position = proto.AvatarPosition_AvatarPositionHand
	case string(NFTTraitTypeLegsArmor): // "Legs Armor" 腿部装备
		position = proto.AvatarPosition_AvatarPositionPant
	case string(NFTTraitTypeFeetArmor): // "Feet Armor" 脚部装备
		position = proto.AvatarPosition_AvatarPositionShoe
	case string(NFTTraitTypeSword): // "Sword" 剑
		position = proto.AvatarPosition_AvatarPositionWeapon
	case string(NFTTraitTypeBow): // "Bow"  弓
		position = proto.AvatarPosition_AvatarPositionWeapon
	case string(NFTTraitTypeDagger): // "Dagger" 匕首
		position = proto.AvatarPosition_AvatarPositionWeapon
	case string(NFTTraitTypeSpear): // "Spear"枪
		position = proto.AvatarPosition_AvatarPositionWeapon

		// case string(NFTTraitTypeConsumable): // "Consumable" 消耗品
		// case string(NFTTraitTypeMaterial): // "Material" 材料
		// case string(NFTTraitTypeMysteryBox): // "MysteryBox" 神秘宝箱
		// case string(NFTTraitTypePlaceable): // "Placeable" 可放置
	}
	return position
}

func (n *NFT) IsEquipment() bool {
	if !n.IsMelandAI {
		return false
	}

	for _, na := range n.Metadata.Attributes {
		if na.TraitType == string(NFTTraitTypesType) {
			return n.isEquipment(na.Value)
		}
	}
	return false
}

func (n *NFT) EquipmentPosition() proto.AvatarPosition {
	if !n.IsMelandAI {
		return proto.AvatarPosition_AvatarPositionNone
	}

	for _, na := range n.Metadata.Attributes {
		if na.TraitType == string(NFTTraitTypesWearingPosition) {
			return n.equipmentPosition(na.Value)
		}
	}

	return proto.AvatarPosition_AvatarPositionNone
}

func (n *NFT) GetEquipmentSkill() (skillId int32) {
	if !n.IsMelandAI {
		return 0
	}

	for _, na := range n.Metadata.Attributes {
		switch na.TraitType {
		case string(NFTTraitTypesType):
			if !n.isEquipment(na.Value) {
				return 0
			}

		case string(NFTTraitTypesCoreSkillId):
			skillId = cast.ToInt32(na.Value)
			return

		}
	}

	return
}

func (n *NFT) GetEquipmentData() (isEquipment bool, position proto.AvatarPosition, attribute *proto.AvatarAttribute) {
	if !n.IsMelandAI {
		return isEquipment, position, attribute
	}

	attribute = &proto.AvatarAttribute{Durability: 200}
	for _, na := range n.Metadata.Attributes {
		switch na.TraitType {
		case string(NFTTraitTypesType):
			if isEquipment = n.isEquipment(na.Value); !isEquipment {
				return false, position, nil
			}

		case string(NFTTraitTypesQuality):
			attribute.Rarity = na.Value
		case string(NFTTraitTypesWearingPosition):
			position = n.equipmentPosition(na.Value)

		case string(NFTTraitTypesMaxHP): // "MaxHP"
			value := cast.ToInt32(na.Value)
			attribute.Data = append(attribute.Data,
				&proto.AttributeData{Type: proto.AttributeType_AttributeTypeHpLimit, Value: value},
			)
		case string(NFTTraitTypesHPRecovery): // "HP Recovery"
			value := cast.ToInt32(na.Value)
			attribute.Data = append(attribute.Data,
				&proto.AttributeData{Type: proto.AttributeType_AttributeTypeHpRecovery, Value: value},
			)
		case string(NFTTraitTypesAttack):
			value := cast.ToInt32(na.Value)
			attribute.Data = append(attribute.Data,
				&proto.AttributeData{Type: proto.AttributeType_AttributeTypeAtt, Value: value},
			)
		case string(NFTTraitTypesAttackSpeed):
			value := cast.ToInt32(na.Value)
			attribute.Data = append(attribute.Data,
				&proto.AttributeData{Type: proto.AttributeType_AttributeTypeAttSpeed, Value: value},
			)
		case string(NFTTraitTypesDefence): // "Defence"
			value := cast.ToInt32(na.Value)
			attribute.Data = append(attribute.Data,
				&proto.AttributeData{Type: proto.AttributeType_AttributeTypeDef, Value: value},
			)

		case string(NFTTraitTypesCritPoints): // "Defence"
			value := cast.ToInt32(na.Value)
			attribute.Data = append(attribute.Data,
				&proto.AttributeData{Type: proto.AttributeType_AttributeTypeCrit, Value: value},
			)
		case string(NFTTraitTypesCritDamage): // "Crit Damage"
			value := cast.ToInt32(na.Value)
			attribute.Data = append(attribute.Data,
				&proto.AttributeData{Type: proto.AttributeType_AttributeTypeCritDmg, Value: value},
			)
		case string(NFTTraitTypesHitPoints): // "Hit Points"
			value := cast.ToInt32(na.Value)
			attribute.Data = append(attribute.Data,
				&proto.AttributeData{Type: proto.AttributeType_AttributeTypeHitRate, Value: value},
			)
		case string(NFTTraitTypesDodgePoints): // "Dodge Points"
			value := cast.ToInt32(na.Value)
			attribute.Data = append(attribute.Data,
				&proto.AttributeData{Type: proto.AttributeType_AttributeTypeMissRate, Value: value},
			)
		case string(NFTTraitTypesMoveSpeed): // "Move Speed"
			value := cast.ToInt32(na.Value)
			attribute.Data = append(attribute.Data,
				&proto.AttributeData{Type: proto.AttributeType_AttributeTypeMoveSpeed, Value: value},
			)

		case string(NFTTraitTypesCoreSkillId): // "CoreSkillId"
			value := cast.ToInt32(na.Value)
			attribute.Data = append(attribute.Data,
				&proto.AttributeData{Type: proto.AttributeType_AttributeTypeSkillId, Value: value},
			)
		// case string(NFTTraitTypesSkillLevel): // "SkillLevel"
		// case string(NFTTraitTypesPlaceableLands): // "Placeable Lands"
		// case string(NFTTraitTypesRarity): // "Rarity"
		// case string(NFTTraitTypesRestoreHP): // "Restore HP"
		// case string(NFTTraitTypesSeries): // "Series"

		// case string(NFTTraitTypesWearingPosition): // "Wearing Position"

		default:

		}
	}

	return isEquipment, position, attribute
}

func (n *NFT) UseLevel() (level int32) {
	for _, na := range n.Metadata.Attributes {
		if na.TraitType == string(NFTTraitTypesLevel) {
			lv, err := strconv.ParseInt(na.Value, 10, 64)
			if err == nil {
				level = int32(lv)
			}
			break
		}
	}
	return
}