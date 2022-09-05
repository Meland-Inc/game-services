package message

func (n *NFT) isMaterial(value string) (b bool) {
	switch value {
	case string(NFTTraitTypeMaterial):
		b = true
	}
	return b
}

func (n *NFT) IsMaterial() (b bool) {
	if !n.IsMelandAI {
		return false
	}
	for _, na := range n.Metadata.Attributes {
		if na.TraitType == string(NFTTraitTypesType) {
			return n.isMaterial(na.Value)
		}
	}
	return false
}

func (n *NFT) GetMaterialData() (isMaterial bool, quality NFTTraitQuality) {
	if !n.IsMelandAI {
		return false, NFTTraitQualityBasic
	}

	for _, na := range n.Metadata.Attributes {
		switch na.TraitType {
		case string(NFTTraitTypesType):
			if isMaterial = n.isMaterial(na.Value); !isMaterial {
				return false, NFTTraitQualityBasic
			}

		case string(NFTTraitTypesQuality):
			quality = NFTTraitQuality(na.Value)

		default:

		}
	}

	return
}
