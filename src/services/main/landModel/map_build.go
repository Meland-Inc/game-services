package land_model

import (
	"fmt"
	"game-message-core/proto"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/global/gameDB"
	dbData "github.com/Meland-Inc/game-services/src/global/gameDB/data"
	"github.com/Meland-Inc/game-services/src/global/grpcAPI/grpcInvoke"
	"github.com/Meland-Inc/game-services/src/global/grpcAPI/grpcPubsubEvent"
	message "github.com/Meland-Inc/game-services/src/global/web3Message"
	"github.com/Meland-Inc/game-services/src/services/main/playerModel"
)

func (p *MapLandDataRecord) addUsingLandRecord(build *NftBuildData) {
	if build == nil {
		return
	}
	for _, landId := range build.InLandIds() {
		p.usingLand[landId] = build
	}
}

func (p *MapLandDataRecord) removeUsingLandRecord(build *NftBuildData) {
	if build == nil {
		return
	}
	for _, landId := range build.InLandIds() {
		delete(p.usingLand, landId)
	}
}

func (p *MapLandDataRecord) loadBuildGameData() ([]dbData.NftBuild, error) {
	var builds []dbData.NftBuild
	err := gameDB.GetGameDB().Where("map_id = ?", p.MapId).Find(&builds).Error
	return builds, err
}

func (p *MapLandDataRecord) InitBuildData() error {
	// load all land data for land-service
	web3BuildDatas, err := grpcInvoke.RPCLoadNftBuildData(p.MapId)
	if err != nil {
		return err
	}
	gameBuildDatas, err := p.loadBuildGameData()
	if err != nil {
		return err
	}

	for _, web3build := range web3BuildDatas {
		exist := false
		for _, gameBuild := range gameBuildDatas {
			if gameBuild.NftId == web3build.NftId {
				nftBuild := NewNftBuildData(gameBuild, web3build)
				p.buildRecord[nftBuild.GetNftId()] = nftBuild
				p.addUsingLandRecord(nftBuild)
				exist = true
				break
			}
		}
		if !exist {
			serviceLog.Error("[]web3build[%s] gameBuildData not found", web3build.UserId, web3build.NftId)
			// todo ... call recycling not found game data nftBuild
		}
	}

	return nil
}

func (p *MapLandDataRecord) GetAllNftBuild() map[string]*NftBuildData {
	p.RLock()
	defer p.RUnlock()
	return p.buildRecord
}

func (p *MapLandDataRecord) GetUserNftBuilds(userId int64) (builds []*NftBuildData) {
	p.RLock()
	defer p.RUnlock()
	for _, build := range p.buildRecord {
		if build.GetOwner() == userId {
			builds = append(builds, build)
		}
	}
	return
}

func (p *MapLandDataRecord) getBuildByEntityId(entityId int64) *NftBuildData {
	for _, build := range p.buildRecord {
		if build.GetEntityId() == entityId {
			return build
		}
	}
	return nil
}

func (p *MapLandDataRecord) GetNftBuildByEntityId(entityId int64) *NftBuildData {
	p.RLock()
	defer p.RUnlock()
	return p.getBuildByEntityId(entityId)
}

// if not found return nil
func (p *MapLandDataRecord) getBuildByNftId(nftId string) *NftBuildData {
	build, _ := p.buildRecord[nftId]
	return build
}

func (p *MapLandDataRecord) GetNftBuildByNftId(nftId string) *NftBuildData {
	p.RLock()
	defer p.RUnlock()
	return p.getBuildByNftId(nftId)
}

func (p *MapLandDataRecord) UpdateNftBuildWeb3Data(data message.BuildData) error {
	p.RLock()
	defer p.RUnlock()
	build := p.getBuildByNftId(data.NftId)
	if build == nil {
		return fmt.Errorf("NftId[%s] build not found", data.NftId)
	}
	build.SetWeb3Data(&data)
	grpcPubsubEvent.RPCPubsubEventNftBuildUpdate(build.ToGrpcData())
	p.BroadcastBuildUpdate(build)
	return nil
}

func (p *MapLandDataRecord) addNftBuildRecord(build *NftBuildData) error {
	if build == nil {
		return fmt.Errorf("add build record build s is nil")
	}

	err := p.playerDataModel.UpdateItemUseState(build.GetOwner(), build.GetNftId(), true, 0)
	if err != nil {
		return err
	}
	err = gameDB.GetGameDB().Create(build.GameData).Error
	if err != nil {
		return err
	}
	p.buildRecord[build.GetNftId()] = build
	p.addUsingLandRecord(build)
	return nil
}

func (p *MapLandDataRecord) removeNftBuildRecord(build *NftBuildData) error {
	if build == nil {
		return fmt.Errorf("delete nft build is nil")
	}
	err := p.playerDataModel.UpdateItemUseState(build.GetOwner(), build.GetNftId(), false, 0)
	if err != nil {
		return err
	}
	delete(p.buildRecord, build.GetNftId())
	p.removeUsingLandRecord(build)
	return gameDB.GetGameDB().Delete(build.GameData).Error
}

func (p *MapLandDataRecord) canBuild(
	userId int64, nftId string, pos *proto.Vector3, landIds []int32,
) (*playerModel.Item, error) {
	item, err := p.playerDataModel.ItemById(userId, nftId)
	if err != nil {
		return nil, err
	}

	if pos == nil {
		return nil, fmt.Errorf("pos is nil")
	}

	for _, landId := range landIds {
		landData, exist := p.landRecord[landId]
		if !exist {
			return nil, fmt.Errorf("land[%d] not found", landId)
		}
		if landData.GetOwner() != userId {
			return nil, fmt.Errorf("can't build other owner land[%d]", landId)
		}
		if _, exist := p.usingLand[landId]; exist {
			return nil, fmt.Errorf("land[%d] is be build", landId)
		}
	}
	return item, nil
}

// 使用 nft建造建筑物
func (p *MapLandDataRecord) Build(
	userId int64, nftId string, pos *proto.Vector3, landIds []int32,
) (*NftBuildData, error) {
	p.RLock()
	defer p.RUnlock()
	item, err := p.canBuild(userId, nftId, pos, landIds)
	if err != nil {
		return nil, err
	}

	web3BuildData, err := grpcInvoke.RPCBuild(userId, nftId, p.MapId, landIds)
	if err != nil {
		return nil, err
	}

	gameBuildData := dbData.NewNftBuild(userId, nftId, item.Cid, p.MapId, pos, landIds)
	nftBuild := NewNftBuildData(*gameBuildData, *web3BuildData)
	p.addNftBuildRecord(nftBuild)
	grpcPubsubEvent.RPCPubsubEventNftBuildAdd(nftBuild.ToGrpcData())
	return nftBuild, nil
}

// 拆除建筑物
func (p *MapLandDataRecord) Recycling(userId int64, nftId string, buildId int64) error {
	p.RLock()
	defer p.RUnlock()
	build := p.getBuildByNftId(nftId)
	if build == nil {
		return fmt.Errorf("nft[%s] build not found", nftId)
	}
	if build.GetOwner() != userId {
		return fmt.Errorf("can't recycling other owner builds")
	}

	err := grpcInvoke.RPCRecyclingBuild(userId, nftId, p.MapId)
	if err != nil {
		return err
	}

	p.removeNftBuildRecord(build)
	grpcPubsubEvent.RPCPubsubEventNftBuildRemove(build.ToGrpcData())
	return nil
}

// 充电
func (p *MapLandDataRecord) BuildCharged(userId int64, nftId string, buildId int64, num int32) error {
	p.RLock()
	defer p.RUnlock()
	build := p.getBuildByNftId(nftId)
	if build == nil {
		return fmt.Errorf("nft[%s] build not found", nftId)
	}
	if build.GetOwner() != userId {
		return fmt.Errorf("can't charged other owner builds")
	}
	return grpcInvoke.RPCBuildCharged(userId, nftId, p.MapId, num)
}

// 收获(harvest)自己建造物的产出(有电量的建造物)
func (p *MapLandDataRecord) Harvest(userId int64, nftId string, buildId int64) error {
	p.RLock()
	defer p.RUnlock()
	build := p.getBuildByNftId(nftId)
	if build == nil {
		return fmt.Errorf("nft[%s] build not found", nftId)
	}
	if build.GetOwner() != userId {
		return fmt.Errorf("can't Harvest other owner builds")
	}
	return grpcInvoke.RPCHarvest(userId, nftId, p.MapId)
}

// 采集/偷取(collection) 他人的或者自己的没电量的建造物产出
func (p *MapLandDataRecord) Collection(userId int64, nftId string, buildId int64) error {
	p.RLock()
	defer p.RUnlock()
	build := p.getBuildByNftId(nftId)
	if build == nil {
		return fmt.Errorf("nft[%s] build not found", nftId)
	}
	if build.GetOwner() != userId {
		return fmt.Errorf("can't Harvest other owner builds")
	}
	return grpcInvoke.RPCCollection(userId, nftId, p.MapId)
}
