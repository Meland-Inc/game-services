package message

import (
	"game-message-core/proto"
)

// ----------------- 穿戴 --------------------------------------------------

func (n *NFT) isWearable(value string) (b bool) {
	switch value {
	case string(NFTTraitTypeWearable): // "Wearable"
		b = true
	}
	return b
}

func (n *NFT) IsWearable() (b bool) {
	if !n.IsMelandAI {
		return false
	}

	for _, na := range n.Metadata.Attributes {
		if na.TraitType == string(NFTTraitTypesType) {
			return n.isWearable(na.Value)
		}
	}
	return false
}

func (n *NFT) wearablePosition(value string) (position proto.AvatarPosition) {
	switch value {
	case string(NFTTraitWearingPositionHead):
		position = proto.AvatarPosition_AvatarPositionHead

	case string(NFTTraitWearingPositionGloves):
		position = proto.AvatarPosition_AvatarPositionHand

	case string(NFTTraitWearingPositionUpperBody):
		position = proto.AvatarPosition_AvatarPositionCoat

	case string(NFTTraitWearingPositionLowerBody):
		position = proto.AvatarPosition_AvatarPositionPant

	case string(NFTTraitWearingPositionShoes):
		position = proto.AvatarPosition_AvatarPositionShoe

	}
	return position
}

func (n *NFT) GetWearablePbData() (isWearable bool, position proto.AvatarPosition, attribute *proto.AvatarAttribute) {
	if !n.IsMelandAI {
		return
	}

	attribute = &proto.AvatarAttribute{Durability: 200}
	for _, na := range n.Metadata.Attributes {
		switch na.TraitType {
		case string(NFTTraitTypesType):
			if isWearable = n.isWearable(na.Value); !isWearable {
				return false, position, nil
			}

		case string(NFTTraitTypesRarity):
			attribute.Rarity = na.Value

		case string(NFTTraitTypesWearingPosition):
			position = n.wearablePosition(na.Value)

		default:

		}
	}

	return isWearable, position, attribute
}
