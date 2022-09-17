package itemModel

import (
	"fmt"

	"github.com/Meland-Inc/game-services/src/global/gameDB"
	dbData "github.com/Meland-Inc/game-services/src/global/gameDB/data"
)

func (p *ItemModel) loadUsingNfts(userId int64) ([]dbData.UsingNft, error) {
	if userId < 1 {
		return nil, fmt.Errorf("all using nft invalid user id [%d]", userId)
	}

	var usingNfts []dbData.UsingNft
	err := gameDB.GetGameDB().Where("user_id = ?", userId).Find(&usingNfts).Error
	return usingNfts, err
}

func (p *ItemModel) addUsingNftRecord(item *Item) error {
	if item == nil {
		return fmt.Errorf("add using nft item is nil")
	}

	usingNft := &dbData.UsingNft{
		NftId:     item.Id,
		UserId:    item.Owner,
		Cid:       item.Cid,
		AvatarPos: item.AvatarPos,
	}
	return gameDB.GetGameDB().Create(usingNft).Error
}

func (p *ItemModel) removeUsingNftRecord(userId int64, nftId string) error {
	if nftId == "" {
		return fmt.Errorf("delete using nft id is nil")
	}

	usingNft := dbData.UsingNft{}
	err := gameDB.GetGameDB().Where("nftId = ? ", nftId).First(&usingNft).Error
	if err != nil {
		return err
	}

	return gameDB.GetGameDB().Delete(&usingNft).Error
}
