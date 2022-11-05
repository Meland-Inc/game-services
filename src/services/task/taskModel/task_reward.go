package taskModel

import (
	"fmt"
	"game-message-core/grpc/pubsubEventData"
	"game-message-core/proto"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/common/time_helper"
	"github.com/Meland-Inc/game-services/src/global/configData"
	"github.com/Meland-Inc/game-services/src/global/gameDB"
	dbData "github.com/Meland-Inc/game-services/src/global/gameDB/data"
	"github.com/Meland-Inc/game-services/src/global/grpcAPI/grpcInvoke"
	"github.com/Meland-Inc/game-services/src/global/grpcAPI/grpcPubsubEvent"
)

func (p *TaskModel) givePlayerReward(
	userId int64, tl *dbData.TaskList, fromTaskList bool,
	exp, itemRewardId int32,
) {

	rewardItems, err := configData.RandomRewardItems(itemRewardId)
	if err != nil {
		serviceLog.Error(err.Error())
	}
	if len(rewardItems) > 0 {
		go func() {
			for _, item := range rewardItems {
				if item.Cid > 0 && item.Num > 0 {
					err := grpcInvoke.MintNFT(userId, item.Cid, item.Num, item.Quality, 0, 0)
					if err != nil {
						serviceLog.Error("WEB3 mint NFT failed err: %v", err)
					}
				}
			}
		}()
	}

	env := &pubsubEventData.UserTaskRewardEvent{
		MsgVersion:     time_helper.NowUTCMill(),
		UserId:         userId,
		Exp:            exp,
		TaskListReward: fromTaskList,
	}
	grpcPubsubEvent.RPCPubsubEventTaskReward(env)
	p.broadCastReceiveRewardInfo(userId, tl, fromTaskList, exp, rewardItems)
}

func (p *TaskModel) TaskReward(userId int64, taskListKind proto.TaskListType) (*dbData.TaskList, error) {
	if taskListKind == proto.TaskListType_TaskListTypeGuide {
		return p.guideTaskReward(userId, taskListKind)
	}
	return p.normalTaskReward(userId, taskListKind)
}

func (p *TaskModel) normalTaskReward(userId int64, taskListKind proto.TaskListType) (*dbData.TaskList, error) {
	pt, err := p.GetPlayerTask(userId)
	if err != nil {
		return nil, err
	}
	var tl *dbData.TaskList
	switch taskListKind {
	case proto.TaskListType_TaskListTypeRewarded:
		tl = pt.GetRewardTaskList()
	case proto.TaskListType_TaskListTypeDaily:
		tl = pt.GetDailyTaskList()
	case proto.TaskListType_TaskListTypeGuide:
		tl = pt.GetGuideTaskList()
	}
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
	p.givePlayerReward(userId, tl, false, taskCnf.RewardExp, taskCnf.RewardId)

	tl.Rate++
	tl.CurTask = nil
	// 领取任务奖励 自动接取
	if !tl.IsFinish() {
		newTask, err := p.randomTask(tl)
		if err != nil {
			serviceLog.Error(err.Error())
		} else {
			tl.CurTask = newTask
		}
	}

	switch taskListKind {
	case proto.TaskListType_TaskListTypeRewarded:
		pt.SetRewardTaskList(tl)
	case proto.TaskListType_TaskListTypeDaily:
		pt.SetDailyTaskList(tl)
	}

	err = gameDB.GetGameDB().Save(pt).Error
	p.broadCastUpdateTaskListInfo(userId, taskListKind, tl)
	return tl, err
}

func (p *TaskModel) TaskListReward(userId int64, taskListKind proto.TaskListType) (*dbData.TaskList, error) {
	pt, err := p.GetPlayerTask(userId)
	if err != nil {
		return nil, err
	}
	var tl *dbData.TaskList
	switch taskListKind {
	case proto.TaskListType_TaskListTypeRewarded:
		tl = pt.GetRewardTaskList()
	case proto.TaskListType_TaskListTypeDaily:
		tl = pt.GetDailyTaskList()
	}
	if tl == nil {
		return tl, fmt.Errorf("task list not found")
	}

	switch taskListKind {
	case proto.TaskListType_TaskListTypeRewarded:
		rewardRateMin := configData.ConfigMgr().RewardTaskListRateMin()
		receivedRewardRate := tl.Rate / rewardRateMin
		if tl.ReceiveReward >= receivedRewardRate {
			return tl, fmt.Errorf("can't received reward")
		}

	case proto.TaskListType_TaskListTypeDaily:
		if !tl.IsFinish() || tl.ReceiveReward > 0 {
			return tl, fmt.Errorf("can't received reward or received reward")
		}
	}

	tlCnf := configData.ConfigMgr().TaskListCnfById(tl.TaskListId)
	if tlCnf == nil {
		return tl, fmt.Errorf("task list cur task not finish")
	}
	p.givePlayerReward(userId, tl, true, tlCnf.RewardExp, tl.TaskListId)

	// 阶段奖励记录++
	tl.ReceiveReward++

	switch taskListKind {
	case proto.TaskListType_TaskListTypeRewarded:
		if tl.ReceiveReward == 2 {
			tl.Doing = false
			tl.Rate = 0
			tl.ReceiveReward = 0
		}
		pt.SetRewardTaskList(tl)
	case proto.TaskListType_TaskListTypeDaily:
		tl.Doing = false
		tl.Rate = 0
		tl.ReceiveReward = 0
		pt.SetDailyTaskList(tl)
	}

	err = gameDB.GetGameDB().Save(pt).Error
	return tl, err
}
