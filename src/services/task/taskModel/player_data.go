package taskModel

import (
	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/global/gameDB"
	dbData "github.com/Meland-Inc/game-services/src/global/gameDB/data"
	"github.com/Meland-Inc/game-services/src/global/grpcAPI/grpcInvoke"
	message "github.com/Meland-Inc/game-services/src/global/web3Message"
)

func (p *TaskModel) getPlayerSceneData(userId int64) (*dbData.PlayerSceneData, error) {
	data := &dbData.PlayerSceneData{}
	err := gameDB.GetGameDB().Where("user_id = ?", userId).First(data).Error
	return data, err
}

func (p *TaskModel) getPlayerNFT(userId int64) ([]message.NFT, error) {
	output, err := grpcInvoke.RPCLoadUserNFTS(userId)
	if err != nil {
		serviceLog.Error("loadNft User[%v] NFTS err : %+v", userId, err)
		return nil, err
	}

	return output.Nfts, nil
}
