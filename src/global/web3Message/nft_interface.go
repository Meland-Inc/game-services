package message

import (
	"game-message-core/proto"
)

func ParsePlaceableRarity(daprRarity string) (rarity proto.NFTRarity) {
	switch daprRarity {
	case string(NFTTraitRarityEpic):
		rarity = proto.NFTRarity_NFTRarityEpic
	case string(NFTTraitRarityMythic):
		rarity = proto.NFTRarity_NFTRarityMythic
	case string(NFTTraitRarityRare):
		rarity = proto.NFTRarity_NFTRarityRare
	case string(NFTTraitRarityUnique):
		rarity = proto.NFTRarity_NFTRarityUnique
	default:
		rarity = proto.NFTRarity_NFTRarityCommon
	}
	return rarity
}

func NFTPbType(nft NFT) proto.NFTType {
	if nft.IsThird() {
		return proto.NFTType_NFTTypeThird
	}

	if nft.IsConsumable() {
		return proto.NFTType_NFTTypeConsumable
	}

	if nft.IsEquipment() {
		return proto.NFTType_NFTTypeEquipment
	}

	if nft.IsWearable() {
		return proto.NFTType_NFTTypeWearable
	}

	if nft.IsMaterial() {
		return proto.NFTType_NFTTypeMaterial
	}

	if nft.IsPlaceable() {
		return proto.NFTType_NFTTypePlaceable
	}

	return proto.NFTType_NFTTypeUnknown
}
