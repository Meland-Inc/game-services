package message

import (
	"game-message-core/proto"
)

func (n *NFT) IsThird() bool {
	return !n.IsMelandAI

}

func (n *NFT) GetThirdPbData() (isThird bool, data *proto.NftThirdNftInfo) {
	if n.IsMelandAI {
		return false, nil
	}
	data = &proto.NftThirdNftInfo{
		Name:       "",
		ResUrl:     n.TokenURL,
		TokenUrl:   n.TokenURL,
		TokenId:    n.TokenId,
		TimeOutSec: 0,
	}

	return
}
