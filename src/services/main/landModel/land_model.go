package land_model

import (
	"fmt"
	"sync"

	"github.com/Meland-Inc/game-services/src/global/component"
)

type LandModel struct {
	modelMgr         *component.ModelManager
	modelName        string
	mapList          []int32
	mapLandRecordMgr sync.Map
}

func GetLandModel() (*LandModel, error) {
	iLandModel, exist := component.GetInstance().GetModel(component.MODEL_NAME_LAND)
	if !exist {
		return nil, fmt.Errorf("player data model not found")
	}
	landModel, _ := iLandModel.(*LandModel)
	return landModel, nil
}

func NewLandModel() *LandModel {
	return &LandModel{mapList: []int32{10001}}
}

func (p *LandModel) Name() string {
	return p.modelName
}

func (p *LandModel) ModelMgr() *component.ModelManager {
	return p.modelMgr
}

func (p *LandModel) OnInit(modelMgr *component.ModelManager) error {
	if modelMgr == nil {
		return fmt.Errorf("land model init service model manager is nil")
	}
	p.modelMgr = modelMgr
	p.modelName = component.MODEL_NAME_LAND
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
		if err != nil {
			return false
		}
		return true
	})
	return err
}

func (p *LandModel) OnTick(curMs int64) error {
	return p.tick()
}

func (p *LandModel) OnStop() error {
	p.modelMgr = nil
	return nil
}

func (p *LandModel) OnExit() error {
	return nil
}

func (p *LandModel) tick() error {
	return nil
}

func (p *LandModel) GetMapLandRecord(mapId int32) (*MapLandDataRecord, error) {
	mapRecord, exist := p.mapLandRecordMgr.Load(mapId)
	if !exist {
		return nil, fmt.Errorf("map[%d] record not found", mapId)
	}
	return mapRecord.(*MapLandDataRecord), nil
}
