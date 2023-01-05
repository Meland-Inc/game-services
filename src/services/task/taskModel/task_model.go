package taskModel

import (
	"fmt"
	"time"

	"github.com/Meland-Inc/game-services/src/common/shardCache"
	"github.com/Meland-Inc/game-services/src/global/contract"
	"github.com/Meland-Inc/game-services/src/global/module"
)

type TaskModel struct {
	module.ModuleBase

	cache    *shardCache.ShardedCache
	cacheTTL time.Duration
}

func GetTaskModel() (*TaskModel, error) {
	iModel, exist := module.GetModel(module.MODULE_NAME_TASK)
	if !exist {
		return nil, fmt.Errorf("task model not found")
	}
	dataModel, _ := iModel.(*TaskModel)
	return dataModel, nil
}

func NewTaskModel() *TaskModel {
	p := &TaskModel{}
	p.InitBaseModel(p, module.MODULE_NAME_TASK)
	return p
}

func (p *TaskModel) OnInit() error {
	p.ModuleBase.OnInit()
	p.cacheTTL = time.Duration(10) * time.Minute
	p.cache = shardCache.NewSharded(shardCache.NoExpiration, time.Duration(60)*time.Second, 2^4)
	return nil
}

func (p *TaskModel) OnTick(utc time.Time) {
	p.ModuleBase.OnTick(utc)
}

func (p *TaskModel) Secondly(utc time.Time) {}
func (p *TaskModel) Minutely(utc time.Time) {}
func (p *TaskModel) Hourly(utc time.Time)   {}
func (p *TaskModel) Daily(utc time.Time) {
	p.checkAndResetPlayerTask(utc)
}

func (p *TaskModel) EventCall(env contract.IModuleEventReq) contract.IModuleEventResult {
	return nil
}
func (p *TaskModel) EventCallNoReturn(env contract.IModuleEventReq) {}
func (p *TaskModel) ReadEvent() contract.IModuleEventReq {
	return nil
}
