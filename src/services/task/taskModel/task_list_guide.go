package taskModel

import (
	"fmt"
	"game-message-core/proto"
	"time"

	"github.com/Meland-Inc/game-services/src/global/configData"
	"github.com/Meland-Inc/game-services/src/global/gameDB"
	dbData "github.com/Meland-Inc/game-services/src/global/gameDB/data"
	"github.com/Meland-Inc/game-services/src/global/grpcAPI/grpcPubsubEvent"
)

func (p *TaskModel) nextGuideTaskId(taskListId, curTaskId int32) (int32, error) {
	taskListCnf := configData.ConfigMgr().TaskListCnfById(taskListId)
	if taskListCnf == nil {
		return 0, fmt.Errorf("TaskList [%d] config not found", taskListId)
	}
	taskSeq, err := taskListCnf.GetSequence()
	if err != nil {
		return 0, err
	}
	if len(taskSeq.Sequence) < 1 {
		return 0, fmt.Errorf("guide taskList  task Sequence empty")
	}

	// 第一次接取任务
	if curTaskId == 0 {
		return taskSeq.Sequence[0], nil
	}

	nextTaskIdx := -1
	for idx, id := range taskSeq.Sequence {
		if id == curTaskId {
			nextTaskIdx = idx + 1
			break
		}
	}
	if nextTaskIdx == -1 {
		return 0, fmt.Errorf("guide taskList not found cur task [%d]", curTaskId)
	}

	if nextTaskIdx >= len(taskSeq.Sequence) {
		return 0, nil
	}

	return taskSeq.Sequence[nextTaskIdx], nil
}

func (p *TaskModel) getNextGuideTask(userId int64, tl *dbData.TaskList) (*dbData.Task, error) {
	player, err := p.getPlayerSceneData(userId)
	if err != nil {
		return nil, err
	}

	nextTaskId, err := p.nextGuideTaskId(tl.TaskListId, tl.Rate)
	if err != nil {
		return nil, err
	}
	if nextTaskId == 0 {
		return nil, nil
	}

	taskSetting := configData.ConfigMgr().TaskCnfById(nextTaskId)
	if taskSetting == nil {
		return nil, fmt.Errorf("task[%v] config not found", nextTaskId)
	}
	if taskSetting.Level > player.Level {
		return nil, nil
	}
	taskOpts, err := p.getTaskOptions(nextTaskId)
	if err != nil {
		return nil, err
	}
	p.initTaskOptionsRate(userId, taskOpts)
	now := time.Now().UTC()
	nextTask := &dbData.Task{
		TaskId:    int32(nextTaskId),
		Options:   taskOpts,
		CreatedAt: now,
		UpdatedAt: now,
	}
	return nextTask, nil
}

func (p *TaskModel) getGuideTaskReward(userId int64) (*dbData.TaskList, error) {
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

	gtl.Rate = gtl.CurTask.TaskId
	nextTask, err := p.getNextGuideTask(userId, gtl)
	if err != nil {
		return nil, err
	}

	rewardItem := p.givePlayerReward(userId, gtl, false, taskCnf.RewardExp, taskCnf.RewardId)
	gtl.CurTask = nextTask
	pt.SetGuideTaskList(gtl)
	err = gameDB.GetGameDB().Save(pt).Error
	grpcPubsubEvent.RPCPubsubEventTaskFinish(userId, proto.TaskListType(gtl.TaskListType), taskCnf.Id, rewardItem)
	p.broadCastUpdateTaskListInfo(userId, proto.TaskListType_TaskListTypeGuide, gtl)
	return gtl, err
}
