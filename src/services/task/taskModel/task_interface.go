package taskModel

import (
	"fmt"
	"game-message-core/proto"
	xlsxTable "game-message-core/xlsxTableData"
	"time"

	"github.com/Meland-Inc/game-services/src/common/matrix"
	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/global/configData"
	"github.com/Meland-Inc/game-services/src/global/gameDB"
	dbData "github.com/Meland-Inc/game-services/src/global/gameDB/data"
	"github.com/Meland-Inc/game-services/src/global/grpcAPI/grpcInvoke"
)

func (p *TaskModel) randomTaskOption(taskSetting *xlsxTable.TaskTableRow) ([]*dbData.TaskOption, error) {
	taskCnfOptions, err := taskSetting.GetOptions()
	if err != nil {
		return nil, err
	}

	// 随机获取 任务真实的 完成条件(单条 || 多条)
	taskOption := []*dbData.TaskOption{}
	for _, option := range taskCnfOptions.Options {
		var realOption *xlsxTable.TaskTableOptionParam
		rn := matrix.Random32(0, option.RandomExclusive)
		for _, rl := range option.RandList {
			probability := rl.Param3 // TODO: ... 定义各种任务option 的数据规则。。。@雨越
			if rn <= rl.Param3 {
				realOption = &rl
			} else {
				rn -= probability
			}
		}

		if realOption == nil {
			serviceLog.Error("tasks[%v] option data is invalid", taskSetting.Id)
			continue
		}
		dbOpt := &dbData.TaskOptionCnf{
			TaskOptionType: int32(option.OptionType),
			Param1:         realOption.Param1,
			Param2:         realOption.Param2,
			Param3:         realOption.Param3,
			Param4:         realOption.Param4,
			Param5:         realOption.Param5,
		}
		taskOption = append(taskOption, &dbData.TaskOption{OptionCnf: dbOpt})
	}
	return taskOption, nil
}

func (p *TaskModel) randomTask(tl *dbData.TaskList) (*dbData.Task, error) {
	if tl == nil {
		return nil, fmt.Errorf("task list is nil")
	}

	taskListCnf := configData.ConfigMgr().TaskListCnfById(tl.TaskListId)
	if taskListCnf == nil {
		return nil, fmt.Errorf("TaskList [%d] config not found", tl.TaskListId)
	}
	taskList, err := taskListCnf.GetIncludeTask()
	if err != nil {
		return nil, err
	}

	var taskId int32
	rn := matrix.Random32(0, taskList.ChanceSum)
	for _, parm := range taskList.Param {
		if rn <= parm.Chance {
			taskId = parm.TaskId
			break
		} else {
			rn -= parm.Chance
		}
	}

	taskSetting := configData.ConfigMgr().TaskCnfById(taskId)
	if taskSetting == nil {
		return nil, fmt.Errorf("task[%v] config not found", taskId)
	}

	// 随机获取 任务真实的 完成条件(单条 || 多条)
	taskOption, err := p.randomTaskOption(taskSetting)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	pt := &dbData.Task{
		TaskId:    int32(taskId),
		Options:   taskOption,
		CreatedAt: now,
		UpdatedAt: now,
	}
	return pt, nil
}

func (p *TaskModel) randomTaskList(userId int64, tlType proto.TaskListType) (*dbData.TaskList, error) {
	player, err := p.getPlayerSceneData(userId)
	if err != nil {
		return nil, err
	}

	cnf := configData.ConfigMgr().TaskListCnfByLevel(int32(tlType), player.Level)
	if cnf == nil {
		serviceLog.Warning("task list[%v], lv[%v] config not found", tlType, player.Level)
		return nil, nil
	}

	now := time.Now().UTC()
	curTl := &dbData.TaskList{
		CanReceive:   true,
		Doing:        false,
		TaskListId:   int32(cnf.Id),
		TaskListType: cnf.System,
		CreatedAt:    now,
		UpdatedAt:    now,
		ResetAt:      now,
	}
	if tlType == proto.TaskListType_TaskListTypeRewarded {
		curTl.ResetAt = rewardTaskLastResetTime()
	}

	// init player guide task list auto accept first task
	if tlType == proto.TaskListType_TaskListTypeGuide {
		nextTask, err := p.randomTask(curTl)
		if err != nil {
			return nil, err
		}
		curTl.CurTask = nextTask
	}

	return curTl, nil
}

func (p *TaskModel) InitPlayerTask(userId int64) (*dbData.PlayerTask, error) {
	dtl, err := p.randomTaskList(userId, proto.TaskListType_TaskListTypeDaily)
	if err != nil {
		return nil, err
	}
	rtl, err := p.randomTaskList(userId, proto.TaskListType_TaskListTypeRewarded)
	if err != nil {
		return nil, err
	}
	gtl, err := p.randomTaskList(userId, proto.TaskListType_TaskListTypeGuide)
	if err != nil {
		return nil, err
	}
	pt := dbData.NewPlayerTask(userId, dtl, rtl, gtl)
	if err := gameDB.GetGameDB().Save(pt).Error; err != nil {
		return nil, err
	}
	return pt, nil
}

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

	nextTask, err := p.randomTask(rtl)
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

	task, err := p.randomTask(dtl)
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

func (p *TaskModel) updatePlayerTaskList(userId int64, taskList *dbData.TaskList) error {
	playerTl, err := p.GetPlayerTask(userId)
	if err != nil {
		return err
	}

	switch proto.TaskListType(taskList.TaskListType) {
	case proto.TaskListType_TaskListTypeDaily:
		playerTl.SetDailyTaskList(taskList)
	case proto.TaskListType_TaskListTypeRewarded:
		playerTl.SetRewardTaskList(taskList)
	case proto.TaskListType_TaskListTypeGuide:
		playerTl.SetGuideTaskList(taskList)
	default:
		return fmt.Errorf("up player task list type [%v] not define", taskList.TaskListType)
	}

	err = gameDB.GetGameDB().Save(playerTl).Error
	if err != nil {
		return err
	}
	p.broadCastUpdateTaskListInfo(userId, proto.TaskListType(taskList.TaskListType), taskList)
	return nil
}

func (p *TaskModel) AcceptTask(userId int64, kind proto.TaskListType) (*dbData.TaskList, error) {
	if kind == proto.TaskListType_TaskListTypeRewarded {
		return p.acceptRewardedTask(userId)
	}
	return p.acceptDailyTask(userId)
}

// 放弃任务
func (p *TaskModel) AbandonmentTask(userId int64, kind proto.TaskListType) (*dbData.TaskList, error) {
	pt, err := p.GetPlayerTask(userId)
	if err != nil {
		return nil, err
	}

	tl := pt.GetDailyTaskList()
	if kind == proto.TaskListType_TaskListTypeRewarded {
		tl = pt.GetRewardTaskList()
	}

	if tl == nil || tl.CurTask == nil {
		return nil, fmt.Errorf("task list cur task not found")
	}

	now := time.Now().UTC()
	if now.Unix()-tl.CurTask.CreatedAt.Unix() < 5*60 {
		return nil, fmt.Errorf("task list cur task is Protected")
	}

	tl.CurTask = nil

	if kind == proto.TaskListType_TaskListTypeRewarded {
		pt.SetRewardTaskList(tl)
	} else {
		pt.SetDailyTaskList(tl)
	}

	err = gameDB.GetGameDB().Save(pt).Error
	return tl, err
}
