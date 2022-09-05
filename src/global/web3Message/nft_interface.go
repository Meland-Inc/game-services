package message

import (
	"game-message-core/proto"
)

func ParseBigWorldLandState(state LandStatus) (landStatus proto.BigWorldLandState) {
	switch state {
	case LandStatusOccupied:
		landStatus = proto.BigWorldLandState_BigWorldLandStateOccupied
	case LandStatusTicket:
		landStatus = proto.BigWorldLandState_BigWorldLandStateTicket
	case LandStatusUnoccupied:
		landStatus = proto.BigWorldLandState_BigWorldLandStateUnoccupied
	case LandStatusVIP:
		landStatus = proto.BigWorldLandState_BigWorldLandStateVip
	}
	return
}

func ParseBigWorldLandStateList(statesList []string) (landStatus []proto.BigWorldLandState) {
	for _, stateStr := range statesList {
		state := LandStatus(stateStr)
		landStatus = append(landStatus, ParseBigWorldLandState(state))
	}
	return
}

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

func ParseToBigWorldLandState(land string) (landStatus proto.BigWorldLandState) {
	switch land {
	case string(NFTTraitPlaceableLandsVIP):
		landStatus = proto.BigWorldLandState_BigWorldLandStateVip
	case string(NFTTraitPlaceableLandsTicket):
		landStatus = proto.BigWorldLandState_BigWorldLandStateTicket
	case string(NFTTraitPlaceableLandsOccupied):
		landStatus = proto.BigWorldLandState_BigWorldLandStateOccupied
	default:
		landStatus = proto.BigWorldLandState_BigWorldLandStateUnoccupied
	}
	return
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
