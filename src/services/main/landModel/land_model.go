package land_model

import (
	"fmt"
	"sync"
	"time"

	"github.com/Meland-Inc/game-services/src/global/module"
	"github.com/Meland-Inc/game-services/src/services/main/playerModel"
)

type LandModel struct {
	module.ModuleBase

	mapList          []int32
	mapLandRecordMgr sync.Map
}

func GetLandModel() (*LandModel, error) {
	iLandModel, exist := module.GetModel(module.MODULE_NAME_LAND)
	if !exist {
		return nil, fmt.Errorf("land  model not found")
	}
	landModel, _ := iLandModel.(*LandModel)
	return landModel, nil
}

func NewLandModel() *LandModel {
	p := &LandModel{mapList: []int32{10001}}
	p.InitBaseModel(p, module.MODULE_NAME_LAND)
	return p
}

func (p *LandModel) OnInit() error {
	p.ModuleBase.OnInit()
	for _, mapId := range p.mapList {
		mapRecord := NewMapLandDataRecord(mapId)
		p.mapLandRecordMgr.Store(mapId, mapRecord)
	}
	return nil
}

func (p *LandModel) OnStart() (err error) {
	p.mapLandRecordMgr.Range(func(key, value interface{}) bool {
		mapRecord := value.(*MapLandDataRecord)
		err = mapRecord.OnStart()
		return err == nil
	})
	return err
}

func (p *LandModel) OnTick(utc time.Time) {
	p.ModuleBase.OnTick(utc)
}

func (p *LandModel) Secondly(utc time.Time) {}
func (p *LandModel) Minutely(utc time.Time) {}
func (p *LandModel) Hourly(utc time.Time)   {}
func (p *LandModel) Daily(utc time.Time)    {}

func (p *LandModel) GetMapLandRecord(mapId int32) (*MapLandDataRecord, error) {
	mapRecord, exist := p.mapLandRecordMgr.Load(mapId)
	if !exist {
		return nil, fmt.Errorf("map[%d] land record not found", mapId)
	}
	return mapRecord.(*MapLandDataRecord), nil
}

func (p *LandModel) GetMapLandRecordByUser(userId int64) (*MapLandDataRecord, error) {
	playerDataModel, _ := playerModel.GetPlayerDataModel()
	playerData, err := playerDataModel.GetPlayerSceneData(userId)
	if err != nil {
		return nil, err
	}
	return p.GetMapLandRecord(playerData.MapId)
}
