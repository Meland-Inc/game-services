package playerModel

import (
	"fmt"
	"time"

	"github.com/Meland-Inc/game-services/src/common/shardCache"
	"github.com/Meland-Inc/game-services/src/global/contract"
	"github.com/Meland-Inc/game-services/src/global/module"
)

type PlayerDataModel struct {
	module.ModuleBase
	modelEvent *module.ModuleEvent

	cache    *shardCache.ShardedCache
	cacheTTL time.Duration
}

func GetPlayerDataModel() (*PlayerDataModel, error) {
	iPlayerModel, exist := module.GetModel(module.MODULE_NAME_PLAYER_DATA)
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
	p.InitBaseModel(p, module.MODULE_NAME_PLAYER_DATA)
	p.modelEvent = module.NewModelEvent()
	return p
}

func (p *PlayerDataModel) OnInit() error {
	p.ModuleBase.OnInit()
	return nil
}

func (p *PlayerDataModel) OnStart() error {
	return nil
}

func (p *PlayerDataModel) OnTick(utc time.Time) {
	p.ModuleBase.OnTick(utc)
	if env := p.ReadEvent(); env != nil {
		p.OnEvent(env, utc.UnixMilli())
	}
}

func (p *PlayerDataModel) Secondly(utc time.Time) {}

func (p *PlayerDataModel) Minutely(utc time.Time) {}

func (p *PlayerDataModel) Hourly(utc time.Time) {}

func (p *PlayerDataModel) Daily(utc time.Time) {}

func (p *PlayerDataModel) EventCall(env contract.IModuleEventReq) contract.IModuleEventResult {
	return p.modelEvent.EventCall(env)
}

func (p *PlayerDataModel) EventCallNoReturn(env contract.IModuleEventReq) {
	p.modelEvent.EventCallNoReturn(env)
}
func (p *PlayerDataModel) ReadEvent() contract.IModuleEventReq {
	return p.modelEvent.ReadEvent()
}
