package land_model

import (
	"fmt"
	"sync"
	"time"

	"github.com/Meland-Inc/game-services/src/global/component"
	"github.com/Meland-Inc/game-services/src/services/main/playerModel"
)

type LandModel struct {
	component.ModelBase

	mapList          []int32
	mapLandRecordMgr sync.Map
	playerDataModel  *playerModel.PlayerDataModel
}

func GetLandModel() (*LandModel, error) {
	iLandModel, exist := component.GetInstance().GetModel(component.MODEL_NAME_LAND)
	if !exist {
		return nil, fmt.Errorf("land  model not found")
	}
	landModel, _ := iLandModel.(*LandModel)
	return landModel, nil
}

func NewLandModel() *LandModel {
	return &LandModel{mapList: []int32{10001}}
}

func (p *LandModel) OnInit(modelMgr *component.ModelManager) error {
	if modelMgr == nil {
		return fmt.Errorf("land model init service model manager is nil")
	}
	p.ModelBase.OnInit(modelMgr)
	p.InitBaseModel(p, component.MODEL_NAME_LAND)
	for _, mapId := range p.mapList {
		mapRecord := NewMapLandDataRecord(mapId)
		p.mapLandRecordMgr.Store(mapId, mapRecord)
	}
	return nil
}

func (p *LandModel) OnStart() (err error) {
	p.ModelBase.OnStart()
	p.playerDataModel, err = playerModel.GetPlayerDataModel()
	if err != nil {
		return err
	}

	p.mapLandRecordMgr.Range(func(key, value interface{}) bool {
		mapRecord := value.(*MapLandDataRecord)
		err = mapRecord.OnStart()
		return err == nil
	})
	return err
}

func (p *LandModel) OnTick(utc time.Time) {
}

func (p *LandModel) EventCall(env *component.ModelEventReq) *component.ModelEventResult {
	return nil
}
func (p *LandModel) EventCallNoReturn(env *component.ModelEventReq)    {}
func (p *LandModel) OnEvent(env *component.ModelEventReq, curMs int64) {}

func (p *LandModel) tick() error {
	return nil
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
	playerData, err := p.playerDataModel.GetPlayerSceneData(userId)
	if err != nil {
		return nil, err
	}
	return p.GetMapLandRecord(playerData.MapId)
}
