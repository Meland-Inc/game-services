package message

import (
	"game-message-core/proto"

	"github.com/spf13/cast"
)

func (n *NFT) isPlaceable(value string) (b bool) {
	switch value {
	case string(NFTTraitTypePlaceable):
		b = true
	}
	return b
}

func (n *NFT) IsPlaceable() bool {
	if !n.IsMelandAI {
		return false
	}
	for _, na := range n.Metadata.Attributes {
		if na.TraitType == string(NFTTraitTypesType) {
			return n.isPlaceable(na.Value)
		}
	}
	return false
}

func (n *NFT) GetPlaceableSkill() (skills []*proto.NftSkill) {
	if !n.IsMelandAI {
		return skills
	}

	for _, na := range n.Metadata.Attributes {
		switch na.TraitType {
		case string(NFTTraitTypesType):
			if !n.isPlaceable(na.Value) {
				return []*proto.NftSkill{}
			}

		case string(NFTTraitTypesCoreSkillId):
			skillId := cast.ToInt32(na.Value)
			if len(skills) > 0 {
				skills[0].SkillId = skillId
			} else {
				skills = append(skills, &proto.NftSkill{SkillId: skillId})
			}
		case string(NFTTraitTypesSkillLevel):
			lv := cast.ToInt32(na.Value)
			if len(skills) > 0 {
				skills[0].SkillLevel = lv
			} else {
				skills = append(skills, &proto.NftSkill{SkillLevel: lv})
			}
		}
	}
	return skills
}

func (n *NFT) GetPlaceablePbData() (isPlaceable bool, pbPlaceable *proto.NftPlaceableInfo) {
	if !n.IsMelandAI {
		return false, nil
	}

	pbPlaceable = &proto.NftPlaceableInfo{
		Token:     n.Id,
		ObjectCid: cast.ToInt32(n.ItemId),
	}
	for _, na := range n.Metadata.Attributes {
		switch na.TraitType {
		case string(NFTTraitTypesType):
			if isPlaceable = n.isPlaceable(na.Value); !isPlaceable {
				return false, nil
			}

		case string(NFTTraitTypesCoreSkillId):
			skillId := cast.ToInt32(na.Value)
			if len(pbPlaceable.Skills) > 0 {
				pbPlaceable.Skills[0].SkillId = skillId
			} else {
				pbPlaceable.Skills = append(pbPlaceable.Skills, &proto.NftSkill{SkillId: skillId})
			}
		case string(NFTTraitTypesSkillLevel):
			lv := cast.ToInt32(na.Value)
			if len(pbPlaceable.Skills) > 0 {
				pbPlaceable.Skills[0].SkillLevel = lv
			} else {
				pbPlaceable.Skills = append(pbPlaceable.Skills, &proto.NftSkill{SkillLevel: lv})
			}

		case string(NFTTraitTypesRarity):
			pbPlaceable.PlaceableRarity = ParsePlaceableRarity(na.Value)

		default:

		}
	}

	return isPlaceable, pbPlaceable
}
