package taskModel

import (
	"fmt"
	"game-message-core/proto"

	"github.com/Meland-Inc/game-services/src/common/matrix"
	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	dbData "github.com/Meland-Inc/game-services/src/global/gameDB/data"
	"github.com/Meland-Inc/game-services/src/global/grpcAPI/grpcInvoke"
)

func (p *TaskModel) getPlayerTaskList(
	userId int64, taskListKind proto.TaskListType,
) (checkTls []*dbData.TaskList, err error) {
	pt, err := p.GetPlayerTask(userId)
	if err != nil {
		return nil, err
	}
	gtl := pt.GetGuideTaskList()
	dtl := pt.GetDailyTaskList()
	rtl := pt.GetRewardTaskList()

	switch taskListKind {
	case proto.TaskListType_TaskListTypeRewarded:
		if rtl != nil {
			checkTls = append(checkTls, rtl)
		}
	case proto.TaskListType_TaskListTypeDaily:
		if dtl != nil {
			checkTls = append(checkTls, dtl)
		}
	case proto.TaskListType_TaskListTypeGuide:
		if gtl != nil {
			checkTls = append(checkTls, gtl)
		}
	default:
		if rtl != nil {
			checkTls = append(checkTls, rtl)
		}
		if dtl != nil {
			checkTls = append(checkTls, dtl)
		}
		if gtl != nil {
			checkTls = append(checkTls, gtl)
		}
	}
	return
}

func (p *TaskModel) upgradeTaskOption(
	userId int64, taskListKind proto.TaskListType,
	upOptionF func(*dbData.TaskOption) (upgrade bool),
) error {
	checkTaskLists, err := p.getPlayerTaskList(userId, taskListKind)
	if err != nil {
		return err
	}

	for _, checkTl := range checkTaskLists {
		if checkTl.CurTask == nil {
			continue
		}
		upgrade := false
		for _, opt := range checkTl.CurTask.Options {
			if opt == nil || opt.IsFinish() {
				continue
			}
			if opt.OptionCnf == nil {
				serviceLog.Error("user[%d] taskList[%v] task[%d] optionCnf is nil", userId, taskListKind, checkTl.CurTask.TaskId)
				continue
			}

			// upgrade task rate logic
			upgrade = upOptionF(opt)
		}

		if upgrade {
			p.updatePlayerTaskList(userId, checkTl)
		}
	}
	return nil
}

func (p *TaskModel) HandInItemHandler(
	userId int64, taskListKind proto.TaskListType, handInItems []*proto.TaskOptionItem,
) error {
	for _, it := range handInItems {
		err := grpcInvoke.BurnNFT(userId, it.NftId, it.Num)
		if err != nil {
			serviceLog.Error("handInItem web3 burn nft [%d][%s][%d] fail, error: %v", userId, userId, it.NftId, it.Num, err)
			return err
		}
	}

	return p.upgradeTaskOption(
		userId,
		taskListKind,
		func(taskOption *dbData.TaskOption) (upgrade bool) {
			// upgrade task rate logic
			if taskOption.OptionCnf.TaskOptionType != int32(proto.TaskOptionType_HandInItem) {
				return false
			}
			for _, item := range handInItems {
				if item.Num < 1 {
					continue
				}
				if item.ItemCid != taskOption.OptionCnf.Param1 {
					continue
				}
				rateOffset := taskOption.OptionCnf.Param2 - taskOption.Rate
				if rateOffset < 1 {
					continue
				}

				rateAdd := matrix.LimitInt32Min(item.Num, rateOffset)
				taskOption.Rate += rateAdd
				item.Num -= rateAdd
				upgrade = true
			}
			return false
		})
}

func (p *TaskModel) UseItemHandler(
	userId int64, taskListKind proto.TaskListType, usedItem *proto.TaskOptionItem,
) error {
	if usedItem.Num < 1 || usedItem.ItemCid < 1 {
		return fmt.Errorf("invalid item cid%d] num[%d]", usedItem.ItemCid, usedItem.Num)
	}

	return p.upgradeTaskOption(
		userId,
		taskListKind,
		func(taskOption *dbData.TaskOption) (upgrade bool) {
			// upgrade task rate logic
			if taskOption.OptionCnf.TaskOptionType != int32(proto.TaskOptionType_UseItem) {
				return false
			}
			if usedItem.ItemCid != taskOption.OptionCnf.Param1 {
				return false
			}
			taskOption.Rate += usedItem.Num
			return true
		})
}

func (p *TaskModel) PickUpItemHandler(
	userId int64, taskListKind proto.TaskListType, pickItems []*proto.TaskOptionItem,
) error {
	return p.upgradeTaskOption(
		userId,
		taskListKind,
		func(taskOption *dbData.TaskOption) (upgrade bool) {
			// upgrade task rate logic
			if taskOption.OptionCnf.TaskOptionType != int32(proto.TaskOptionType_PickUpItem) {
				return false
			}
			for _, pickItem := range pickItems {

				if pickItem.ItemCid != taskOption.OptionCnf.Param1 {
					return false
				}
				taskOption.Rate += pickItem.Num
			}
			return true
		})
}

func (p *TaskModel) KillMonsterHandler(
	userId int64, taskListKind proto.TaskListType, killMon *proto.TaskOptionKillMonster,
) error {
	if killMon.MonCid < 1 || killMon.Num < 1 {
		return fmt.Errorf("invalid killMonster data cid%d] num[%d]", killMon.MonCid, killMon.Num)
	}

	return p.upgradeTaskOption(
		userId,
		taskListKind,
		func(taskOption *dbData.TaskOption) (upgrade bool) {
			// upgrade task rate logic
			if taskOption.OptionCnf.TaskOptionType != int32(proto.TaskOptionType_KillMonster) {
				return false
			}
			if killMon.MonCid != taskOption.OptionCnf.Param1 {
				return false
			}
			taskOption.Rate += killMon.Num
			return true
		})
}

func (p *TaskModel) UserLevelHandler(
	userId int64, taskListKind proto.TaskListType, userLv int32,
) error {
	return p.upgradeTaskOption(
		userId,
		taskListKind,
		func(taskOption *dbData.TaskOption) (upgrade bool) {
			// upgrade task rate logic
			if taskOption.OptionCnf.TaskOptionType != int32(proto.TaskOptionType_UserLevel) {
				return false
			}
			taskOption.Rate = userLv
			return true
		})
}

func (p *TaskModel) TargetSlotLevelHandler(
	userId int64, taskListKind proto.TaskListType, slotInfo *proto.TaskOptionTargetSlotLevel,
) error {
	return p.upgradeTaskOption(
		userId,
		taskListKind,
		func(taskOption *dbData.TaskOption) (upgrade bool) {
			// upgrade task rate logic
			if taskOption.OptionCnf.TaskOptionType != int32(proto.TaskOptionType_TargetSlotLevel) {
				return false
			}
			if slotInfo.SlotPos != taskOption.OptionCnf.Param1 {
				return false
			}
			taskOption.Rate = slotInfo.Level
			return true
		})
}

func (p *TaskModel) SlotLevelCountHandler(
	userId int64, taskListKind proto.TaskListType,
) error {
	playerSlotData, err := p.getPlayerSlotData(userId)
	if err != nil {

	}

	return p.upgradeTaskOption(
		userId,
		taskListKind,
		func(taskOption *dbData.TaskOption) (upgrade bool) {
			// upgrade task rate logic
			if taskOption.OptionCnf.TaskOptionType != int32(proto.TaskOptionType_SlotLevelCount) {
				return false
			}
			targetLv := taskOption.OptionCnf.Param1
			var targetLvCount int32
			for _, slot := range playerSlotData.GetSlotList().SlotList {
				if int32(slot.Level) >= targetLv {
					targetLvCount++
				}
			}
			if targetLvCount != taskOption.Rate {
				taskOption.Rate = targetLvCount
				return true
			}
			return false
		})
}

func (p *TaskModel) CraftSkillLevelHandler(
	userId int64, taskListKind proto.TaskListType, craftInfo *proto.TaskOptionCraftSkillLevel,
) error {
	return p.upgradeTaskOption(
		userId,
		taskListKind,
		func(taskOption *dbData.TaskOption) (upgrade bool) {
			// upgrade task rate logic
			if taskOption.OptionCnf.TaskOptionType != int32(proto.TaskOptionType_CraftSkillLevel) {
				return false
			}
			if craftInfo.SkillId != taskOption.OptionCnf.Param1 {
				return false
			}
			taskOption.Rate = craftInfo.Level
			return true
		})
}

func (p *TaskModel) UseRecipeHandler(
	userId int64, taskListKind proto.TaskListType, recipeInfo *proto.TaskOptionUseRecipe,
) error {
	return p.upgradeTaskOption(
		userId,
		taskListKind,
		func(taskOption *dbData.TaskOption) (upgrade bool) {
			// upgrade task rate logic
			if taskOption.OptionCnf.TaskOptionType != int32(proto.TaskOptionType_UseRecipe) {
				return false
			}
			if recipeInfo.RecipeId != taskOption.OptionCnf.Param1 {
				return false
			}
			taskOption.Rate += recipeInfo.Times
			return true
		})
}

func (p *TaskModel) TaskTypeCountHandler(
	userId int64, taskListKind proto.TaskListType, taskInfo *proto.TaskOptionTaskTypeCount,
) error {
	return p.upgradeTaskOption(
		userId,
		taskListKind,
		func(taskOption *dbData.TaskOption) (upgrade bool) {
			// upgrade task rate logic
			if taskOption.OptionCnf.TaskOptionType != int32(proto.TaskOptionType_TaskTypeCount) {
				return false
			}
			if int32(taskInfo.Kind) != taskOption.OptionCnf.Param1 {
				return false
			}
			taskOption.Rate += taskInfo.Count // TODO: 任务类型已经移除只作为任务的子项类型， @雨越
			return true
		})
}

func (p *TaskModel) TargetPositionHandler(
	userId int64, taskListKind proto.TaskListType, position *proto.TaskOptionTargetPosition,
) error {
	return p.upgradeTaskOption(
		userId,
		taskListKind,
		func(taskOption *dbData.TaskOption) (upgrade bool) {
			// upgrade task rate logic
			if taskOption.OptionCnf.TaskOptionType != int32(proto.TaskOptionType_TargetPosition) {
				return false
			}
			// 此处需要check player current position
			taskOption.Rate = 1 // 1 = option is finish
			return true
		})
}
