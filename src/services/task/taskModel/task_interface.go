package taskModel

import (
	"fmt"
	"game-message-core/grpc/pubsubEventData"
	"game-message-core/proto"
	xlsxTable "game-message-core/xlsxTableData"
	"time"

	"github.com/Meland-Inc/game-services/src/common/matrix"
	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/common/time_helper"
	"github.com/Meland-Inc/game-services/src/global/configData"
	"github.com/Meland-Inc/game-services/src/global/gameDB"
	dbData "github.com/Meland-Inc/game-services/src/global/gameDB/data"
	"github.com/Meland-Inc/game-services/src/global/grpcAPI/grpcInvoke"
	"github.com/Meland-Inc/game-services/src/global/grpcAPI/grpcPubsubEvent"
)

func (p *TaskModel) updateTaskOptionRate(
	sceneData *dbData.PlayerSceneData,
	slotData *dbData.ItemSlot,
	opt *dbData.TaskOption,
) (upgrade bool) {
	switch proto.TaskOptionType(opt.OptionCnf.TaskOptionType) {
	case proto.TaskOptionType_UserLevel:
		if opt.Rate != sceneData.Level {
			opt.Rate = sceneData.Level
			upgrade = true
		}
	case proto.TaskOptionType_TargetSlotLevel:
		for _, slot := range slotData.GetSlotList().SlotList {
			if int32(slot.Position) == opt.OptionCnf.Param1 {
				if opt.Rate != int32(slot.Level) {
					opt.Rate = int32(slot.Level)
					upgrade = true
				}
				break
			}
		}
	case proto.TaskOptionType_SlotLevelCount:
		count := int32(0)
		for _, slot := range slotData.GetSlotList().SlotList {
			if int32(slot.Level) >= opt.OptionCnf.Param1 {
				count++
			}
		}
		if opt.Rate != count {
			opt.Rate = count
			upgrade = true
		}
	}

	return upgrade
}

func (p *TaskModel) initTaskOptionsRate(userId int64, options []*dbData.TaskOption) {
	sceneData, err := p.getPlayerSceneData(userId)
	if err != nil {
		serviceLog.Error("getSceneData err: %+v", err)
	}
	slotData, err := p.getPlayerSlotData(userId)
	if err != nil {
		serviceLog.Error("getSlotData err: %+v", err)
	}
	for _, opt := range options {
		p.updateTaskOptionRate(sceneData, slotData, opt)
	}
}

// 获取任务按权重随机的任务项
func (p *TaskModel) getTaskChanceOptions(taskSetting *xlsxTable.TaskTableRow) (*dbData.TaskOptionCnf, error) {
	taskCnfOptions, err := taskSetting.GetChanceOptions()
	if err != nil {
		return nil, err
	}

	// 随机获取 任务权重选项
	if taskCnfOptions != nil {
		rn := matrix.Random32(0, taskCnfOptions.ChanceSum)
		for _, option := range taskCnfOptions.Options {
			if rn <= option.Chance {
				dbOpt := &dbData.TaskOptionCnf{
					TaskOptionType: int32(option.OptionType),
					Param1:         option.Param1,
					Param2:         option.Param2,
					Param3:         option.Param3,
					Param4:         option.Param4,
				}
				return dbOpt, nil
			} else {
				rn -= option.Chance
			}
		}
	}
	return nil, nil
}

// 获取任务指定的任务项
func (p *TaskModel) getTaskDesignateOptions(taskSetting *xlsxTable.TaskTableRow) ([]*dbData.TaskOptionCnf, error) {
	taskCnfOptions, err := taskSetting.GetDesignateOptions()
	if err != nil {
		return nil, err
	}

	optCnfs := []*dbData.TaskOptionCnf{}
	if taskCnfOptions != nil {
		for _, option := range taskCnfOptions.Options {
			dbOpt := &dbData.TaskOptionCnf{
				TaskOptionType: int32(option.OptionType),
				Param1:         option.Param1,
				Param2:         option.Param2,
				Param3:         option.Param3,
				Param4:         option.Param4,
			}
			optCnfs = append(optCnfs, dbOpt)
		}
	}
	return optCnfs, nil
}

func (p *TaskModel) getTaskOptions(taskId int32) ([]*dbData.TaskOption, error) {
	taskSetting := configData.ConfigMgr().TaskCnfById(taskId)
	if taskSetting == nil {
		return nil, fmt.Errorf("task[%v] config not found", taskId)
	}

	// 权重随机获取任务真实的完成条件
	taskChaOpts, err := p.getTaskChanceOptions(taskSetting)
	if err != nil {
		return nil, err
	}
	// 获取完成任务所需的固定的完成条件
	taskDesiOpts, err := p.getTaskDesignateOptions(taskSetting)
	if err != nil {
		return nil, err
	}

	taskOpts := []*dbData.TaskOption{}
	if taskChaOpts != nil {
		taskOpts = append(taskOpts, &dbData.TaskOption{
			OptionCnf: &dbData.TaskOptionCnf{
				TaskOptionType: taskChaOpts.TaskOptionType,
				Param1:         taskChaOpts.Param1,
				Param2:         taskChaOpts.Param2,
				Param3:         taskChaOpts.Param3,
				Param4:         taskChaOpts.Param4,
			},
		})
	}
	for _, opt := range taskDesiOpts {
		taskOpts = append(taskOpts, &dbData.TaskOption{
			OptionCnf: &dbData.TaskOptionCnf{
				TaskOptionType: opt.TaskOptionType,
				Param1:         opt.Param1,
				Param2:         opt.Param2,
				Param3:         opt.Param3,
				Param4:         opt.Param4,
			},
		})
	}
	return taskOpts, nil
}

func (p *TaskModel) getNextNormalTask(userId int64, tl *dbData.TaskList) (*dbData.Task, error) {
	taskListCnf := configData.ConfigMgr().TaskListCnfById(tl.TaskListId)
	if taskListCnf == nil {
		return nil, fmt.Errorf("TaskList [%d] config not found", tl.TaskListId)
	}

	taskList, err := taskListCnf.GetTaskPool()
	if err != nil {
		return nil, err
	}

	var nextTaskId int32
	rn := matrix.Random32(0, taskList.ChanceSum)
	for _, parm := range taskList.Param {
		if rn <= parm.Chance {
			nextTaskId = parm.TaskId
			break
		} else {
			rn -= parm.Chance
		}
	}
	if nextTaskId == 0 {
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

func (p *TaskModel) initTaskList(userId int64, tlType proto.TaskListType) (*dbData.TaskList, error) {
	player, err := p.getPlayerSceneData(userId)
	if err != nil {
		return nil, err
	}

	taskListCnf := configData.ConfigMgr().TaskListCnfByLevel(int32(tlType), player.Level)
	if taskListCnf == nil {
		serviceLog.Warning("task list[%v], lv[%v] config not found", tlType, player.Level)
		return nil, nil
	}

	now := time.Now().UTC()
	curTl := &dbData.TaskList{
		CanReceive:   true,
		Doing:        false,
		TaskListId:   int32(taskListCnf.Id),
		TaskListType: taskListCnf.System,
		CreatedAt:    now,
		UpdatedAt:    now,
		ResetAt:      now,
	}
	if tlType == proto.TaskListType_TaskListTypeRewarded {
		curTl.ResetAt = rewardTaskLastResetTime()
	}

	// init player guide task list auto accept first task
	if tlType == proto.TaskListType_TaskListTypeGuide {
		nextTask, err := p.getNextGuideTask(userId, curTl)
		if err != nil {
			return nil, err
		}
		curTl.CurTask = nextTask
		curTl.Doing = true
		curTl.CanReceive = false
	}

	return curTl, nil
}

func (p *TaskModel) InitPlayerTask(userId int64) (*dbData.PlayerTask, error) {
	dtl, err := p.initTaskList(userId, proto.TaskListType_TaskListTypeDaily)
	if err != nil {
		return nil, err
	}
	rtl, err := p.initTaskList(userId, proto.TaskListType_TaskListTypeRewarded)
	if err != nil {
		return nil, err
	}
	gtl, err := p.initTaskList(userId, proto.TaskListType_TaskListTypeGuide)
	if err != nil {
		return nil, err
	}
	pt := dbData.NewPlayerTask(userId, dtl, rtl, gtl)
	if err := gameDB.GetGameDB().Save(pt).Error; err != nil {
		return nil, err
	}
	return pt, nil
}

func (p *TaskModel) AcceptTask(userId int64, kind proto.TaskListType) (*dbData.TaskList, error) {
	switch kind {
	case proto.TaskListType_TaskListTypeRewarded:
		return p.acceptRewardedTask(userId)
	case proto.TaskListType_TaskListTypeDaily:
		return p.acceptDailyTask(userId)
	default:
		return nil, fmt.Errorf("task list [%v] not found", kind)
	}
}

// 放弃任务
func (p *TaskModel) AbandonmentTask(userId int64, kind proto.TaskListType) (*dbData.TaskList, error) {
	switch kind {
	case proto.TaskListType_TaskListTypeRewarded:
		return p.abandonmentRewardTask(userId)
	case proto.TaskListType_TaskListTypeDaily:
		return p.abandonmentDailyTask(userId)
	default:
		return nil, fmt.Errorf("task list [%v] not found", kind)
	}
}

func (p *TaskModel) givePlayerReward(
	userId int64, tl *dbData.TaskList, fromTaskList bool, exp, itemRewardId int32,
) []*proto.ItemBaseInfo {
	rewardItems, err := configData.RandomRewardItems(itemRewardId)
	if err != nil {
		serviceLog.Error(err.Error())
	}
	if len(rewardItems) > 0 {
		go func() {
			for _, item := range rewardItems {
				if item.Cid > 0 && item.Num > 0 {
					err := grpcInvoke.Web3MintNFT(userId, item.Cid, item.Num, item.Quality, 0, 0)
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
	return rewardItems
}

// 领取任务奖励
func (p *TaskModel) TaskReward(userId int64, kind proto.TaskListType) (*dbData.TaskList, error) {
	switch kind {
	case proto.TaskListType_TaskListTypeRewarded:
		return p.getRewardTaskReward(userId)
	case proto.TaskListType_TaskListTypeDaily:
		return p.getDailyTaskReward(userId)
	case proto.TaskListType_TaskListTypeGuide:
		return p.getGuideTaskReward(userId)
	default:
		return nil, fmt.Errorf("task list [%v] not found", kind)
	}
}

// 领取任务链奖励
func (p *TaskModel) TaskListReward(userId int64, kind proto.TaskListType) (*dbData.TaskList, error) {
	switch kind {
	case proto.TaskListType_TaskListTypeRewarded:
		return p.getRewardTaskListReward(userId)
	case proto.TaskListType_TaskListTypeDaily:
		return p.getDailyTaskListReward(userId)
	default:
		return nil, fmt.Errorf("task list [%v] not found", kind)
	}
}
