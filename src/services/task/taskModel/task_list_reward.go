package taskModel

import (
	"fmt"
	"game-message-core/proto"
	"time"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/global/configData"
	"github.com/Meland-Inc/game-services/src/global/gameDB"
	dbData "github.com/Meland-Inc/game-services/src/global/gameDB/data"
	"github.com/Meland-Inc/game-services/src/global/grpcAPI/grpcInvoke"
	"github.com/Meland-Inc/game-services/src/global/grpcAPI/grpcPubsubEvent"
)

func (p *TaskModel) acceptRewardedTask(userId int64) (*dbData.TaskList, error) {
	pt, err := p.GetPlayerTask(userId)
	if err != nil {
		return nil, err
	}

	rtl := pt.GetRewardTaskList()
	if rtl == nil {
		return nil, fmt.Errorf("rewarded task list is nil")
	}
	if rtl.IsFinish() {
		return rtl, fmt.Errorf("task list is finish")
	}
	if rtl.CurTask != nil {
		return rtl, fmt.Errorf("task cur task is doing")
	}

	cnf := configData.ConfigMgr().TaskListCnfById(rtl.TaskListId)
	if cnf == nil {
		return nil, fmt.Errorf("reward task list [%d] config not found", rtl.TaskListId)
	}

	nextTask, err := p.getNextNormalTask(userId, rtl)
	if err != nil {
		return nil, err
	}

	// 放弃和第一次领取 悬赏任务 需要花费 MELD
	if cnf.NeedMELD > 0 {
		grpcInvoke.BurnUserMELD(userId, int(cnf.NeedMELD))
		if err != nil {
			return rtl, fmt.Errorf("can not find accept Rewarded Task need MELD")
		}
	}

	if !rtl.Doing {
		rtl.CanReceive = false
	}
	rtl.Doing = true
	rtl.CurTask = nextTask
	rtl.UpdatedAt = time.Now().UTC()
	pt.SetRewardTaskList(rtl)
	err = gameDB.GetGameDB().Save(pt).Error
	return rtl, err
}

func (p *TaskModel) abandonmentRewardTask(userId int64) (*dbData.TaskList, error) {
	pt, err := p.GetPlayerTask(userId)
	if err != nil {
		return nil, err
	}

	tl := pt.GetRewardTaskList()
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

func (p *TaskModel) getRewardTaskReward(userId int64) (*dbData.TaskList, error) {
	pt, err := p.GetPlayerTask(userId)
	if err != nil {
		return nil, err
	}

	tl := pt.GetRewardTaskList()
	if tl == nil || tl.CurTask == nil {
		return tl, fmt.Errorf("task list cur task not found")
	}
	if !tl.CurTask.IsFinish() {
		return tl, fmt.Errorf("task list cur task not finish")
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

func (p *TaskModel) getRewardTaskListReward(userId int64) (*dbData.TaskList, error) {
	pt, err := p.GetPlayerTask(userId)
	if err != nil {
		return nil, err
	}

	tl := pt.GetRewardTaskList()
	if tl == nil {
		return tl, fmt.Errorf("task list cur task not found")
	}

	rewardRateMin := configData.ConfigMgr().RewardTaskListRateMin()
	receivedRewardRate := tl.Rate / rewardRateMin
	if tl.ReceiveReward >= receivedRewardRate {
		return tl, fmt.Errorf("can't received reward")
	}

	tlCnf := configData.ConfigMgr().TaskListCnfById(tl.TaskListId)
	if tlCnf == nil {
		return tl, fmt.Errorf("task list cur task not finish")
	}
	rewardItem := p.givePlayerReward(userId, tl, true, tlCnf.RewardExp, tl.ReceiveReward)

	// 阶段奖励记录++
	tl.ReceiveReward++
	if tl.ReceiveReward == 2 {
		tl.Doing = false
		tl.Rate = 0
		tl.ReceiveReward = 0
	}
	err = p.updatePlayerTaskList(userId, tl)
	grpcPubsubEvent.RPCPubsubEventTaskListFinish(userId, proto.TaskListType(tl.TaskListType), rewardItem)
	return tl, err
}
