package playerModel

import (
	"fmt"
	"time"

	"github.com/Meland-Inc/game-services/src/common/shardCache"
	"github.com/Meland-Inc/game-services/src/global/component"
)

type PlayerModel struct {
	modelMgr  *component.ModelManager
	modelName string
	cache     *shardCache.ShardedCache
	cacheTTL  time.Duration
}

func NewPlayerModel() *PlayerModel {
	return &PlayerModel{}
}

func (p *PlayerModel) Name() string {
	return p.modelName
}

func (p *PlayerModel) ModelMgr() *component.ModelManager {
	return p.modelMgr
}

func (p *PlayerModel) OnInit(modelMgr *component.ModelManager) error {
	if modelMgr == nil {
		return fmt.Errorf("player model init service model manager is nil")
	}
	p.modelMgr = modelMgr
	p.modelName = component.MODEL_NAME_PLAYER_DATA
	p.cacheTTL = time.Duration(10) * time.Minute
	p.cache = shardCache.NewSharded(shardCache.NoExpiration, time.Duration(60)*time.Second, 2^4)
	return nil
}

func (p *PlayerModel) OnStart() error {
	return nil
}

func (p *PlayerModel) OnTick(curMs int64) error {
	return p.tick()
}

func (p *PlayerModel) OnStop() error {
	p.modelMgr = nil
	return nil
}

func (p *PlayerModel) OnExit() error {
	return nil
}

func (p *PlayerModel) tick() error {
	return nil
}
