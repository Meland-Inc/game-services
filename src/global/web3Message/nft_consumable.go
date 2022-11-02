package message

import (
	"game-message-core/proto"

	"github.com/spf13/cast"
)

func (n *NFT) isConsumable(value string) (b bool) {
	switch value {
	case string(NFTTraitTypeConsumable): // "Consumable" 消耗品
		b = true
	}
	return b
}

func (n *NFT) IsConsumable() (b bool) {
	if !n.IsMelandAI {
		return false
	}

	for _, na := range n.Metadata.Attributes {
		if na.TraitType == string(NFTTraitTypesType) {
			return n.isConsumable(na.Value)
		}
	}
	return false
}

func (n *NFT) GetConsumableData() (isConsumable bool, data *proto.NFTConsumableInfo) {
	if !n.IsMelandAI {
		return false, nil
	}
	data = &proto.NFTConsumableInfo{}
	for _, na := range n.Metadata.Attributes {
		switch na.TraitType {
		case string(NFTTraitTypesType):
			if isConsumable = n.isConsumable(na.Value); !isConsumable {
				return false, nil
			}

		case string(NFTTraitTypesQuality):
			data.Quality = na.Value

		case string(NFTTraitTypesRestoreHP):
			data.ConsumableType = proto.NFTConsumableType_NFTConsumableTypeRestoreHP
			data.Value = cast.ToInt32(na.Value)

		case string(NFTTraitTypesLearnRecipe):
			data.ConsumableType = proto.NFTConsumableType_NFTConsumableTypeLearnRecipe
			data.Value = cast.ToInt32(na.Value)

		case string(NFTTraitTypesGetBuff):
			data.ConsumableType = proto.NFTConsumableType_NFTConsumableTypeAddBuff
			data.Value = cast.ToInt32(na.Value)

		case string(NFTTraitTypesOccupyLand):
			data.ConsumableType = proto.NFTConsumableType_NFTConsumableTypeOccupyLand

		default:

		}
	}

	return isConsumable, data
}
