package land_model

import (
	"fmt"
	"game-message-core/proto"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/global/grpcAPI/grpcInvoke"
	message "github.com/Meland-Inc/game-services/src/global/web3Message"
)

func (p *MapLandDataRecord) InitLandData() error {
	// load all land data for land-service
	lands, err := grpcInvoke.RPCLoadLandData(p.MapId)
	serviceLog.Info("RPC load land length[%v], err: %+v", len(lands), err)
	if err != nil {
		return err
	}

	p.landRecord = make(map[int32]*proto.LandData)
	for _, land := range lands {
		pbLd := message.ToProtoLandData(land)
		p.landRecord[pbLd.Id] = pbLd
	}
	return nil
}

func (p *MapLandDataRecord) AllLandData() (lands []*proto.LandData, err error) {
	p.RLock()
	defer p.RUnlock()

	for _, land := range p.landRecord {
		lands = append(lands, land)
	}
	return lands, nil
}

// return land data if not found return nil
func (p *MapLandDataRecord) LandById(id int32) (*proto.LandData, error) {
	p.RLock()
	defer p.RUnlock()

	land, exist := p.landRecord[id]
	if !exist {
		return nil, fmt.Errorf("land[%d] not found", id)
	}
	return land, nil
}

func (p *MapLandDataRecord) MultiUpdateLandData(upLands []*proto.LandData) {
	if len(upLands) == 0 {
		return
	}

	p.RLock()
	defer p.RUnlock()

	for _, land := range upLands {
		p.landRecord[land.Id] = land
	}
	p.BroadcastLandDataUpdate(upLands)
}

// 占领地格
func (p *MapLandDataRecord) OccupyLand(userId int64, landId, landPosX, landPosZ int32) error {
	p.RLock()
	defer p.RUnlock()

	land, exist := p.landRecord[landId]
	if !exist {
		return fmt.Errorf("land[%d] not found", landId)
	}
	if land.Owner > 0 {
		return fmt.Errorf("land[%d] is occupied by[%d]", landId, land.Owner)
	}

	userItems, err := p.playerDataModel.GetPlayerItems(userId)
	if err != nil {
		return err
	}
	var occupySeedCount int32
	for _, item := range userItems.Items {
		if item.Cid == 3010204 { //占地种子id
			occupySeedCount += item.Num
			break
		}
	}
	if occupySeedCount < 1 {
		return fmt.Errorf("land seed not found")
	}

	return grpcInvoke.RPCOccupyLand(userId, p.MapId, landId)
}
