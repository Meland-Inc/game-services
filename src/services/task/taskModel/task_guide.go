package taskModel

import (
	"fmt"
	"game-message-core/proto"
	"time"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/global/configData"
	"github.com/Meland-Inc/game-services/src/global/gameDB"
	dbData "github.com/Meland-Inc/game-services/src/global/gameDB/data"
)

func (p *TaskModel) RefreshGuideTask(userId int64, broadcast bool) {
	pt, err := p.GetPlayerTask(userId)
	if err != nil {
		return
	}

	gtl := pt.GetGuideTaskList()
	if gtl == nil {
		serviceLog.Error("refresh [%d] guide task list not found", userId)
		return
	}
	if gtl.CurTask != nil {
		return
	}

	nextTask, err := p.getNextGuideTask(userId, gtl.Rate)
	if err != nil {
		return
	}

	gtl.CurTask = nextTask
	pt.SetDailyTaskList(gtl)
	err = gameDB.GetGameDB().Save(pt).Error
	if err != nil {
		serviceLog.Error(err.Error())
	}
	if broadcast {
		p.broadCastUpdateTaskListInfo(userId, proto.TaskListType_TaskListTypeGuide, gtl)
	}
}

func (p *TaskModel) getNextGuideTask(userId int64, curTaskId int32) (*dbData.Task, error) {
	player, err := p.getPlayerSceneData(userId)
	if err != nil {
		return nil, err
	}

	preTaskCnf := configData.ConfigMgr().TaskCnfById(curTaskId)
	if preTaskCnf == nil {
		return nil, fmt.Errorf("task config data not found")
	}
	if preTaskCnf.NextTaskId < 1 {
		return nil, nil
	}

	// 自动接取下一个任务
	nextTaskCnf := configData.ConfigMgr().TaskCnfById(preTaskCnf.NextTaskId)
	if nextTaskCnf == nil {
		return nil, fmt.Errorf("task[%d] config data not found", preTaskCnf.NextTaskId)
	}

	if player.Level < nextTaskCnf.Level {
		return nil, fmt.Errorf("userLv < next task need lv")
	}

	var nextTask *dbData.Task
	taskOption, err := p.randomTaskOption(nextTaskCnf)
	if err != nil {
		now := time.Now().UTC()
		nextTask = &dbData.Task{
			TaskId:    nextTaskCnf.Id,
			Options:   taskOption,
			CreatedAt: now,
			UpdatedAt: now,
		}
	} else {
		serviceLog.Error(err.Error())
	}
	return nextTask, err
}

func (p *TaskModel) acceptGuideTask(userId int64) (*dbData.TaskList, error) {
	pt, err := p.GetPlayerTask(userId)
	if err != nil {
		return nil, err
	}

	gtl := pt.GetGuideTaskList()
	if gtl == nil {
		return nil, fmt.Errorf("guide task list not found")
	}

	if gtl.CurTask != nil {
		return gtl, fmt.Errorf("cur task is doing")
	}
	gtl.CurTask, err = p.getNextGuideTask(userId, gtl.Rate)
	return gtl, err
}

func (p *TaskModel) guideTaskReward(userId int64, taskListKind proto.TaskListType) (*dbData.TaskList, error) {
	pt, err := p.GetPlayerTask(userId)
	if err != nil {
		return nil, err
	}

	gtl := pt.GetGuideTaskList()
	if gtl == nil || gtl.CurTask == nil {
		return gtl, fmt.Errorf("task list cur task not found")
	}
	if !gtl.CurTask.IsFinish() {
		return gtl, fmt.Errorf("task list cur task not finish")
	}

	taskCnf := configData.ConfigMgr().TaskCnfById(gtl.CurTask.TaskId)
	if taskCnf == nil {
		return gtl, fmt.Errorf("task config data not found")
	}
	p.givePlayerReward(userId, gtl, false, taskCnf.RewardExp, taskCnf.RewardId)

	gtl.Rate = gtl.CurTask.TaskId
	gtl.CurTask = nil

	gtl, err = p.acceptGuideTask(userId)
	if err != nil {
		return nil, err
	}
	pt.SetDailyTaskList(gtl)
	err = gameDB.GetGameDB().Save(pt).Error
	p.broadCastUpdateTaskListInfo(userId, taskListKind, gtl)
	return gtl, err
}
