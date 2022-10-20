package land_model

import (
	"game-message-core/proto"
	"sync"
)

type MapLandDataRecord struct {
	sync.RWMutex
	MapId      int32
	landRecord map[int32]*proto.LandData // map{landId = landData}
}

func NewMapLandDataRecord(mapId int32) *MapLandDataRecord {
	return &MapLandDataRecord{
		MapId:      mapId,
		landRecord: nil,
	}
}
