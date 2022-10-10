package taskModel

import (
	"fmt"
	"game-message-core/grpc/pubsubEventData"
	"game-message-core/proto"
	xlsxTable "game-message-core/xlsxTableData"

	"github.com/Meland-Inc/game-services/src/common/matrix"
	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/common/time_helper"
	"github.com/Meland-Inc/game-services/src/global/configData"
	"github.com/Meland-Inc/game-services/src/global/gameDB"
	dbData "github.com/Meland-Inc/game-services/src/global/gameDB/data"
	"github.com/Meland-Inc/game-services/src/global/grpcAPI/grpcInvoke"
	"github.com/Meland-Inc/game-services/src/global/grpcAPI/grpcPubsubEvent"
)

func (p *TaskModel) givePlayerReward(
	userId int64, tl *dbData.TaskList,
	exp, itemCid, itemNum, itemQuality int32,
	fromTaskList bool,
) {
	if itemCid > 0 && itemNum > 0 {
		err := grpcInvoke.MintNFT(userId, itemCid, itemNum, itemQuality, 0, 0)
		if err != nil {
			serviceLog.Error("WEB3 mint NFT failed err: %v", err)
			return
		}
	}

	env := &pubsubEventData.UserTaskRewardEvent{
		MsgVersion:     time_helper.NowUTCMill(),
		UserId:         userId,
		Exp:            exp,
		TaskListReward: fromTaskList,
	}
	grpcPubsubEvent.RPCPubsubEventTaskReward(env)
	p.broadCastReceiveRewardInfo(userId, tl, exp, itemCid, itemNum, itemQuality, fromTaskList)
}

func (p *TaskModel) randomRewardItem(obj *xlsxTable.TaskObjectList) (cid, num, quality int32) {
	if obj == nil {
		return
	}
	rn := matrix.Random32(0, obj.ChanceSum)
	for _, t := range obj.ParamList {
		if rn <= t.Param2 {
			cid = t.Param1
			num = t.Param3
			quality = 1
			break
		} else {
			rn -= t.Param2
		}
	}
	return
}

func (p *TaskModel) TaskReward(userId int64, taskListKind proto.TaskListType) (*dbData.TaskList, error) {

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

	itemCid, num, quality := p.randomRewardItem(taskCnf.GetRewardItems())
	p.givePlayerReward(userId, tl, taskCnf.RewardExp, itemCid, num, quality, false)

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

	itemRewardList := tlCnf.GetRewardItems()
	itemCid, num, quality := p.randomRewardItem(itemRewardList)
	p.givePlayerReward(userId, tl, tlCnf.RewardExp, itemCid, num, quality, true)

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
