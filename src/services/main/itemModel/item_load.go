package itemModel

import (
	"encoding/json"
	"fmt"

	"github.com/Meland-Inc/game-services/src/common/daprInvoke"
	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/common/time_helper"
	dbData "github.com/Meland-Inc/game-services/src/global/gameDB/data"
	message "github.com/Meland-Inc/game-services/src/global/web3Message"
)

func (p *ItemModel) LoadItems(userId int64) (items []*Item, err error) {
	nfts, err := p.loadNFTS(userId)
	if err != nil {
		return nil, err
	}

	usingNfts, err := p.loadUsingNfts(userId)
	if err != nil {
		return nil, err
	}

	for _, nft := range nfts {
		used, avatarPos := checkUsed(nft, usingNfts)
		nft.Used = used
		nft.AvatarPos = avatarPos
		items = append(items, nft)
	}

	return items, err
}

func (p *ItemModel) loadNFTS(userId int64) ([]*Item, error) {
	beginMs := time_helper.NowUTCMill()
	defer func() {
		serviceLog.Info("web3 load player[%d] NFT use time MS[%V]", userId, time_helper.NowMill()-beginMs)
	}()

	userNfts, err := p.RPCLoadUserNFTS(userId)
	if err != nil {
		serviceLog.Error("loadItemsByDapr User[%v] NFTS err : %+v", userId, err)
		return nil, err
	}

	items, err := p.parseUserNft(userId, userNfts)
	serviceLog.Info("user NFT list = %+v", userNfts)
	serviceLog.Info("user NFT list len(items)=%+v, err: %+v", len(items), err)
	return items, err
}

func (p *ItemModel) RPCLoadUserNFTS(userId int64) (*message.GetUserNFTsOutput, error) {
	input := message.GetUserNFTsInput{UserId: fmt.Sprint(userId)}
	inputBytes, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	outBytes, err := daprInvoke.InvokeMethod(
		string(message.AppIdMelandService),
		string(message.MelandServiceActionGetUserNFTs),
		inputBytes,
	)
	if err != nil {
		serviceLog.Error("load web3 user NFT failed err:%+v", err)
	}

	nfts := &message.GetUserNFTsOutput{}
	err = nfts.UnmarshalJSON(outBytes)
	if err != nil {
		serviceLog.Error("UserPlaceablesOutput Unmarshal : err : %+v", err)
		return nil, err
	}
	return nfts, err
}

func (p *ItemModel) parseUserNft(userId int64, userNfts *message.GetUserNFTsOutput) ([]*Item, error) {
	var items []*Item
	for _, nft := range userNfts.Nfts {
		item := NFTToItem(userId, nft)
		for _, out := range userNfts.PlaceableTimeouts {
			if out.NftId == item.Id {
				item.TimeOut = out
			}
		}
		items = append(items, item)
	}
	serviceLog.Info("user [%v] NFT item Length = %+v", userId, len(items))
	return items, nil

}

func checkUsed(item *Item, usingNfts []dbData.UsingNft) (used bool, avatarPos int32) {
	for _, info := range usingNfts {
		if info.NftId == item.Id {
			return true, int32(info.AvatarPos)
		}
	}
	return false, 0
}
