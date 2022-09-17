package playerModel

import (
	"fmt"
	"time"

	"github.com/Meland-Inc/game-services/src/common/shardCache"
	"github.com/Meland-Inc/game-services/src/global/component"
)

type PlayerDataModel struct {
	modelMgr  *component.ModelManager
	modelName string
	cache     *shardCache.ShardedCache
	cacheTTL  time.Duration
}

func NewPlayerModel() *PlayerDataModel {
	return &PlayerDataModel{}
}

func (p *PlayerDataModel) Name() string {
	return p.modelName
}

func (p *PlayerDataModel) ModelMgr() *component.ModelManager {
	return p.modelMgr
}

func (p *PlayerDataModel) OnInit(modelMgr *component.ModelManager) error {
	if modelMgr == nil {
		return fmt.Errorf("player model init service model manager is nil")
	}
	p.modelMgr = modelMgr
	p.modelName = component.MODEL_NAME_PLAYER_DATA
	p.cacheTTL = time.Duration(10) * time.Minute
	p.cache = shardCache.NewSharded(shardCache.NoExpiration, time.Duration(60)*time.Second, 2^4)
	return nil
}

func (p *PlayerDataModel) OnStart() error {
	return nil
}

func (p *PlayerDataModel) OnTick(curMs int64) error {
	return p.tick()
}

func (p *PlayerDataModel) OnStop() error {
	p.modelMgr = nil
	return nil
}

func (p *PlayerDataModel) OnExit() error {
	return nil
}

func (p *PlayerDataModel) tick() error {
	return nil
}
