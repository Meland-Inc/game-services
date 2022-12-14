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
		serviceLog.Error("RPCLoadNftBuildData err : %+v", err)
		return err
	}
	gameBuildDatas, err := p.loadBuildGameData()
	if err != nil {
		serviceLog.Error("loadBuildGameData err : %+v", err)
		return err
	}

	serviceLog.Info("INitBuildData web3Build len[%d], gameBuildDatas[%d]", len(web3BuildDatas), len(gameBuildDatas))

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

func (p *MapLandDataRecord) getBuildById(buildId int64) *NftBuildData {
	for _, build := range p.buildRecord {
		if build.GetBuildId() == buildId {
			return build
		}
	}
	return nil
}

func (p *MapLandDataRecord) GetNftBuildById(buildId int64) *NftBuildData {
	p.RLock()
	defer p.RUnlock()
	return p.getBuildById(buildId)
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
	p.buildRecord[build.GetNftId()] = build
	p.addUsingLandRecord(build)
	return gameDB.GetGameDB().Create(build.GameData).Error
}

func (p *MapLandDataRecord) removeNftBuildRecord(build *NftBuildData) (err error) {
	if build == nil {
		return fmt.Errorf("delete nft build is nil")
	}

	p.removeUsingLandRecord(build)
	delete(p.buildRecord, build.GetNftId())
	if err1 := gameDB.GetGameDB().Where(
		"build_id=? AND map_id=?", build.Web3Data.BuildId, build.Web3Data.MapId,
	).Delete(&dbData.NftBuild{}).Error; err1 != nil {
		serviceLog.Error(err1.Error())
		err = err1
	}

	err2 := p.playerDataModel.UpdateItemUseState(build.GetOwner(), build.GetNftId(), false, 0)
	if err2 != nil {
		if err2 := gameDB.GetGameDB().Where(
			"nft_id = ? ", build.Web3Data.NftId,
		).First(&dbData.UsingNft{}).Error; err2 != nil {
			serviceLog.Error(err2.Error())
			err = err2
		}
	}

	return err
}

func (p *MapLandDataRecord) canBuild(
	userId int64, nftId string, pos *proto.Vector3, landIds []int32,
) (*playerModel.Item, error) {
	item, err := p.playerDataModel.ItemById(userId, nftId)
	if err != nil {
		return nil, err
	}

	if item.Used {
		return nil, fmt.Errorf("nft is used")
	}

	if pos == nil {
		return nil, fmt.Errorf("pos is nil")
	}

	for _, landId := range landIds {
		landData, exist := p.landRecord[landId]
		if !exist {
			return nil, fmt.Errorf("land[%d] not found", landId)
		}
		if owner := landData.GetOwner(); owner > 0 && owner != userId {
			return nil, fmt.Errorf("can't build other owner land[%d]", landId)
		}
		if _, exist := p.usingLand[landId]; exist {
			return nil, fmt.Errorf("land[%d] is be build", landId)
		}
	}
	return item, nil
}

// ?????? nft???????????????
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

	gameBuildData := dbData.NewNftBuild(userId, int64(web3BuildData.BuildId), nftId, item.Cid, p.MapId, pos, landIds)
	nftBuild := NewNftBuildData(*gameBuildData, *web3BuildData)
	if err = p.addNftBuildRecord(nftBuild); err != nil {
		return nil, err
	}
	grpcPubsubEvent.RPCPubsubEventNftBuildAdd(nftBuild.ToGrpcData())
	return nftBuild, nil
}

// ???????????????
func (p *MapLandDataRecord) Recycling(userId int64, buildId int64) error {
	p.RLock()
	defer p.RUnlock()
	build := p.getBuildById(buildId)
	if build == nil {
		return fmt.Errorf("buildId[%d] build not found", buildId)
	}
	if owner := build.GetOwner(); owner > 0 && owner != userId {
		return fmt.Errorf("can't recycling other owner builds")
	}

	return grpcInvoke.RPCRecyclingBuild(userId, buildId, p.MapId)
}

// ???????????????
func (p *MapLandDataRecord) OnReceiveRecyclingEvent(buildId int64) error {
	p.RLock()
	defer p.RUnlock()
	serviceLog.Info("OnReceiveRecyclingEvent [%d]", buildId)

	build := p.getBuildById(buildId)
	if build == nil {
		return fmt.Errorf("buildId[%d] build not found", buildId)
	}

	if err := p.removeNftBuildRecord(build); err != nil {
		serviceLog.Error(err.Error())
	}
	grpcPubsubEvent.RPCPubsubEventNftBuildRemove(build.ToGrpcData())
	p.BroadcastBuildRecycling(build)
	return nil
}

// ??????
func (p *MapLandDataRecord) BuildCharged(userId int64, nftId string, buildId int64, num int32) error {
	p.RLock()
	defer p.RUnlock()

	if num < 1 {
		return fmt.Errorf("invalid charged num[%d]", num)
	}

	build := p.getBuildById(buildId)
	if build == nil {
		return fmt.Errorf("build[%s] build not found", nftId)
	}
	if owner := build.GetOwner(); owner > 0 && owner != userId {
		return fmt.Errorf("can't charged other owner builds")
	}
	return grpcInvoke.RPCBuildCharged(userId, buildId, p.MapId, num)
}

// ??????(harvest)????????????????????????(?????????????????????)
func (p *MapLandDataRecord) Harvest(userId int64, nftId string, buildId int64) error {
	p.RLock()
	defer p.RUnlock()
	build := p.getBuildById(buildId)
	if build == nil {
		return fmt.Errorf("build[%s]  not found", nftId)
	}
	if owner := build.GetOwner(); owner > 0 && owner != userId {
		return fmt.Errorf("can't Harvest other owner builds")
	}
	return grpcInvoke.RPCHarvest(userId, buildId, p.MapId)
}

// ??????/??????(collection) ???????????????????????????????????????????????????
func (p *MapLandDataRecord) Collection(userId int64, nftId string, buildId int64) error {
	p.RLock()
	defer p.RUnlock()
	build := p.getBuildById(buildId)
	if build == nil {
		return fmt.Errorf("build[%s] not found", nftId)
	}
	return grpcInvoke.RPCCollection(userId, buildId, p.MapId)
}
