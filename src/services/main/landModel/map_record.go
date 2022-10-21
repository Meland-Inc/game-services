package land_model

import (
	"game-message-core/proto"
	"sync"

	"github.com/Meland-Inc/game-services/src/services/main/playerModel"
)

type MapLandDataRecord struct {
	sync.RWMutex
	MapId           int32
	landRecord      map[int32]*proto.LandData // map{landId = landData}
	buildRecord     map[string]*NftBuildData  // map{nftId = NftBuildData}
	usingLand       map[int32]*NftBuildData   // map{landId = NftBuildData}
	playerDataModel *playerModel.PlayerDataModel
}

func NewMapLandDataRecord(mapId int32) *MapLandDataRecord {
	return &MapLandDataRecord{
		MapId:       mapId,
		landRecord:  make(map[int32]*proto.LandData),
		buildRecord: make(map[string]*NftBuildData),
		usingLand:   make(map[int32]*NftBuildData),
	}
}

func (p *MapLandDataRecord) OnStart() error {
	if err := p.InitLandData(); err != nil {
		return err
	}

	if err := p.InitBuildData(); err != nil {
		return err
	}

	playerDataModel, err := playerModel.GetPlayerDataModel()
	if err != nil {
		return err
	}
	p.playerDataModel = playerDataModel
	return nil
}
