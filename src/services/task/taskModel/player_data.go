package taskModel

import (
	"errors"
	"game-message-core/grpc/methodData"
	"game-message-core/proto"

	"github.com/Meland-Inc/game-services/src/global/gameDB"
	dbData "github.com/Meland-Inc/game-services/src/global/gameDB/data"
	"github.com/Meland-Inc/game-services/src/global/grpcAPI/grpcInvoke"
)

func (p *TaskModel) getPlayerSceneData(userId int64) (*dbData.PlayerSceneData, error) {
	data := &dbData.PlayerSceneData{}
	err := gameDB.GetGameDB().Where("user_id = ?", userId).First(data).Error
	return data, err
}

func (p *TaskModel) takeUserNft(userId int64, items []*proto.TaskOptionItem) error {
	if len(items) < 1 {
		return nil
	}
	if userId < 1 {
		return errors.New("invalid user id zero")
	}

	takeNfts := []methodData.TakeNftData{}
	for _, it := range items {
		if it.Num < 1 || (it.ItemCid < 1 && it.NftId == "") {
			continue
		}
		takeNfts = append(takeNfts, methodData.TakeNftData{
			NftId:   it.NftId,
			ItemCid: it.ItemCid,
			Num:     it.Num,
		})
	}
	return grpcInvoke.RPCMainServiceTakeNFT(userId, takeNfts)
}
