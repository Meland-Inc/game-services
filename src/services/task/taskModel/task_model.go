package taskModel

import (
	"fmt"
	"time"

	"github.com/Meland-Inc/game-services/src/common/shardCache"
	"github.com/Meland-Inc/game-services/src/global/component"
)

type TaskModel struct {
	modelMgr  *component.ModelManager
	modelName string
	cache     *shardCache.ShardedCache
	cacheTTL  time.Duration
}

func NewTaskModel() *TaskModel {
	return &TaskModel{}
}

func GetTaskModel() (*TaskModel, error) {
	iModel, exist := component.GetInstance().GetModel(component.MODEL_NAME_TASK)
	if !exist {
		return nil, fmt.Errorf("task model not found")
	}
	dataModel, _ := iModel.(*TaskModel)
	return dataModel, nil
}

func (p *TaskModel) Name() string {
	return p.modelName
}

func (p *TaskModel) ModelMgr() *component.ModelManager {
	return p.modelMgr
}

func (p *TaskModel) OnInit(modelMgr *component.ModelManager) error {
	if modelMgr == nil {
		return fmt.Errorf("task model init service model manager is nil")
	}
	p.modelMgr = modelMgr
	p.modelName = component.MODEL_NAME_TASK
	p.cacheTTL = time.Duration(10) * time.Minute
	p.cache = shardCache.NewSharded(shardCache.NoExpiration, time.Duration(60)*time.Second, 2^4)
	return nil
}

func (p *TaskModel) OnStart() error {
	return nil
}

func (p *TaskModel) OnTick(curMs int64) error {
	return p.taskTick(curMs)
}

func (p *TaskModel) OnStop() error {
	p.modelMgr = nil
	return nil
}

func (p *TaskModel) OnExit() error {
	return nil
}

func (p *TaskModel) EventCall(env *component.ModelEventReq) *component.ModelEventResult {
	return nil
}
func (p *TaskModel) EventCallNoReturn(env *component.ModelEventReq)    {}
func (p *TaskModel) OnEvent(env *component.ModelEventReq, curMs int64) {}
