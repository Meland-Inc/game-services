package land_model

import (
	"fmt"
	"game-message-core/proto"

	"github.com/Meland-Inc/game-services/src/global/grpcAPI/grpcInvoke"
	message "github.com/Meland-Inc/game-services/src/global/web3Message"
)

func (p *MapLandDataRecord) InitLandData() error {
	// load all land data for land-service
	lands, err := grpcInvoke.RPCLoadLandData(p.MapId)
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
