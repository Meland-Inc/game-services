package playerModel

import (
	"fmt"
	"time"

	"github.com/Meland-Inc/game-services/src/common/shardCache"
	"github.com/Meland-Inc/game-services/src/global/component"
)

type PlayerDataModel struct {
	component.ModelBase
	modelEvent *component.ModelEvent

	cache    *shardCache.ShardedCache
	cacheTTL time.Duration
}

func GetPlayerDataModel() (*PlayerDataModel, error) {
	iPlayerModel, exist := component.GetInstance().GetModel(component.MODEL_NAME_PLAYER_DATA)
	if !exist {
		return nil, fmt.Errorf("player data model not found")
	}
	dataModel, _ := iPlayerModel.(*PlayerDataModel)
	return dataModel, nil
}

func NewPlayerModel() *PlayerDataModel {
	p := &PlayerDataModel{
		cacheTTL: time.Duration(10) * time.Minute,
		cache:    shardCache.NewSharded(shardCache.NoExpiration, time.Duration(60)*time.Second, 2^4),
	}
	p.InitBaseModel(p, component.MODEL_NAME_PLAYER_DATA)
	p.modelEvent = component.NewModelEvent(p)
	return p
}

func (p *PlayerDataModel) OnInit(modelMgr *component.ModelManager) error {
	if modelMgr == nil {
		return fmt.Errorf("player model init service model manager is nil")
	}
	p.ModelBase.OnInit(modelMgr)
	return nil
}

func (p *PlayerDataModel) OnStart() error {
	return nil
}

func (p *PlayerDataModel) OnTick(utc time.Time) {
	p.ModelBase.OnTick(utc)
	p.modelEvent.ReadEvent(utc.UnixMilli())
}

func (p *PlayerDataModel) EventCall(env *component.ModelEventReq) *component.ModelEventResult {
	return p.modelEvent.EventCall(env)
}
func (p *PlayerDataModel) EventCallNoReturn(env *component.ModelEventReq) {
	p.modelEvent.EventCallNoReturn(env)
}

func (p *PlayerDataModel) Secondly(utc time.Time) {}

func (p *PlayerDataModel) Minutely(utc time.Time) {}

func (p *PlayerDataModel) Hourly(utc time.Time) {}

func (p *PlayerDataModel) Daily(utc time.Time) {}
