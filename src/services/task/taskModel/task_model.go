package taskModel

import (
	"fmt"
	"time"

	"github.com/Meland-Inc/game-services/src/common/shardCache"
	"github.com/Meland-Inc/game-services/src/global/component"
)

type TaskModel struct {
	component.ModelBase

	cache    *shardCache.ShardedCache
	cacheTTL time.Duration
}

func NewTaskModel() *TaskModel {
	p := &TaskModel{}
	p.InitBaseModel(p, component.MODEL_NAME_TASK)
	return p
}

func GetTaskModel() (*TaskModel, error) {
	iModel, exist := component.GetInstance().GetModel(component.MODEL_NAME_TASK)
	if !exist {
		return nil, fmt.Errorf("task model not found")
	}
	dataModel, _ := iModel.(*TaskModel)
	return dataModel, nil
}

func (p *TaskModel) OnInit(modelMgr *component.ModelManager) error {
	if modelMgr == nil {
		return fmt.Errorf("task model init service model manager is nil")
	}
	p.ModelBase.OnInit(modelMgr)
	p.cacheTTL = time.Duration(10) * time.Minute
	p.cache = shardCache.NewSharded(shardCache.NoExpiration, time.Duration(60)*time.Second, 2^4)
	return nil
}

func (p *TaskModel) OnTick(utc time.Time) {
	p.ModelBase.OnTick(utc)
}

func (p *TaskModel) EventCall(env *component.ModelEventReq) *component.ModelEventResult {
	return nil
}
func (p *TaskModel) EventCallNoReturn(env *component.ModelEventReq)    {}
func (p *TaskModel) OnEvent(env *component.ModelEventReq, curMs int64) {}

func (p *TaskModel) Secondly(utc time.Time) {}
func (p *TaskModel) Minutely(utc time.Time) {}
func (p *TaskModel) Hourly(utc time.Time)   {}
func (p *TaskModel) Daily(utc time.Time) {
	p.checkAndResetPlayerTask(utc)
}
