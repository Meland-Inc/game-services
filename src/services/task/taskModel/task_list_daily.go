package taskModel

import (
	"fmt"
	"game-message-core/proto"
	"time"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/global/configData"
	"github.com/Meland-Inc/game-services/src/global/gameDB"
	dbData "github.com/Meland-Inc/game-services/src/global/gameDB/data"
	"github.com/Meland-Inc/game-services/src/global/grpcAPI/grpcPubsubEvent"
)

func (p *TaskModel) acceptDailyTask(userId int64) (*dbData.TaskList, error) {
	pt, err := p.GetPlayerTask(userId)
	if err != nil {
		return nil, err
	}

	dtl := pt.GetDailyTaskList()
	if dtl == nil {
		return nil, fmt.Errorf("daily task list not found")
	}

	if dtl.IsFinish() {
		return dtl, fmt.Errorf("task list is finish")
	}
	if dtl.CurTask != nil {
		return dtl, fmt.Errorf("cur task is doing")
	}

	task, err := p.getNextNormalTask(userId, dtl)
	if err != nil {
		return nil, err
	}

	dtl.CurTask = task
	dtl.CanReceive = false
	dtl.Doing = true
	dtl.UpdatedAt = time.Now().UTC()
	pt.SetDailyTaskList(dtl)
	err = gameDB.GetGameDB().Save(pt).Error
	return dtl, err
}

// 放弃任务
func (p *TaskModel) abandonmentDailyTask(userId int64) (*dbData.TaskList, error) {
	pt, err := p.GetPlayerTask(userId)
	if err != nil {
		return nil, err
	}

	tl := pt.GetDailyTaskList()
	if tl == nil || tl.CurTask == nil {
		return nil, fmt.Errorf("task list cur task not found")
	}

	now := time.Now().UTC()
	if now.Unix()-tl.CurTask.CreatedAt.Unix() < 5*60 {
		return nil, fmt.Errorf("task list cur task is Protected")
	}

	tl.CurTask = nil
	err = p.updatePlayerTaskList(userId, tl)
	return tl, err
}

func (p *TaskModel) getDailyTaskReward(userId int64) (*dbData.TaskList, error) {
	pt, err := p.GetPlayerTask(userId)
	if err != nil {
		return nil, err
	}

	tl := pt.GetDailyTaskList()
	if tl == nil || tl.CurTask == nil {
		return tl, fmt.Errorf("daily task list cur task not found")
	}

	if !tl.CurTask.IsFinish() {
		return tl, fmt.Errorf("daily task list cur task not finish")
	}

	taskCnf := configData.ConfigMgr().TaskCnfById(tl.CurTask.TaskId)
	if taskCnf == nil {
		return tl, fmt.Errorf("task config data not found")
	}
	rewardItem := p.givePlayerReward(userId, tl, false, taskCnf.RewardExp, taskCnf.RewardId)

	tl.Rate++
	tl.CurTask = nil

	// 领取任务奖励 自动接取
	if !tl.IsFinish() {
		newTask, err := p.getNextNormalTask(userId, tl)
		if err != nil {
			serviceLog.Error(err.Error())
		} else {
			tl.CurTask = newTask
		}
	}

	err = p.updatePlayerTaskList(userId, tl)
	grpcPubsubEvent.RPCPubsubEventTaskFinish(userId, proto.TaskListType(tl.TaskListType), taskCnf.Id, rewardItem)
	return tl, err
}

func (p *TaskModel) getDailyTaskListReward(userId int64) (*dbData.TaskList, error) {
	pt, err := p.GetPlayerTask(userId)
	if err != nil {
		return nil, err
	}

	tl := pt.GetDailyTaskList()
	if tl == nil {
		return tl, fmt.Errorf("task list not found")
	}

	if !tl.IsFinish() || tl.ReceiveReward > 0 {
		return tl, fmt.Errorf("can't received reward or received reward")
	}

	tlCnf := configData.ConfigMgr().TaskListCnfById(tl.TaskListId)
	if tlCnf == nil {
		return tl, fmt.Errorf("task list cur task not finish")
	}
	rewardItem := p.givePlayerReward(userId, tl, true, tlCnf.RewardExp, tl.ReceiveReward)

	// 阶段奖励记录++
	tl.ReceiveReward++
	tl.Doing = false
	tl.Rate = 0
	tl.ReceiveReward = 0
	err = p.updatePlayerTaskList(userId, tl)
	grpcPubsubEvent.RPCPubsubEventTaskListFinish(userId, proto.TaskListType(tl.TaskListType), rewardItem)
	return tl, err
}
