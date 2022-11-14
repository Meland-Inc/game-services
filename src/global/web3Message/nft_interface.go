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

	return proto.NFTType_NFTTypeUnknown
}

func ToProtoNftData(nft NFT) *proto.NftData {
	pbNft := &proto.NftData{
		Network:    nft.Network,
		TokenId:    nft.TokenId,
		IsMelandAi: nft.IsMelandAI,
	}

	pbNft.Metadata = &proto.NftMetadata{
		Name:            nft.Metadata.Name,
		Description:     nft.Metadata.Description,
		Image:           *nft.Metadata.Image,
		BackGroundColor: *nft.Metadata.BackgroundColor,
		Attributes:      []*proto.NftAttribute{},
	}

	for _, attr := range nft.Metadata.Attributes {
		pbNft.Metadata.Attributes = append(
			pbNft.Metadata.Attributes,
			&proto.NftAttribute{
				TraitType: attr.TraitType,
				Value:     attr.Value,
			},
		)
	}

	return pbNft
}
