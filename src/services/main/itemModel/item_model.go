package itemModel

import (
	"fmt"
	"time"

	"github.com/Meland-Inc/game-services/src/common/shardCache"
	"github.com/Meland-Inc/game-services/src/global/component"
)

type ItemModel struct {
	modelMgr  *component.ModelManager
	modelName string
	cache     *shardCache.ShardedCache
	cacheTTL  time.Duration
}

func NewItemModel() *ItemModel {
	return &ItemModel{}
}

func (p *ItemModel) Name() string {
	return p.modelName
}

func (p *ItemModel) ModelMgr() *component.ModelManager {
	return p.modelMgr
}

func (p *ItemModel) OnInit(modelMgr *component.ModelManager) error {
	if modelMgr == nil {
		return fmt.Errorf("player model init service model manager is nil")
	}
	p.modelMgr = modelMgr
	p.modelName = component.MODEL_NAME_ITEM_DATA
	p.cacheTTL = time.Duration(10) * time.Minute
	p.cache = shardCache.NewSharded(shardCache.NoExpiration, time.Duration(60)*time.Second, 2^4)
	return nil
}

func (p *ItemModel) OnStart() error {
	return nil
}

func (p *ItemModel) OnTick(curMs int64) error {
	return p.tick()
}

func (p *ItemModel) OnStop() error {
	p.modelMgr = nil
	return nil
}

func (p *ItemModel) OnExit() error {
	return nil
}

func (p *ItemModel) tick() error {
	return nil
}
