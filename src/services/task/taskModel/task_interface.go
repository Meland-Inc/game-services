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

func (p *TaskModel) randomTask(tl *dbData.TaskList) (*dbData.Task, error) {
	randTaskF := func(obj xlsxTable.TaskObjectList) (param1, param2 int32) {
		rn := matrix.Random32(0, obj.ChanceSum)
		for _, t := range obj.ParamList {
			if rn <= t.Param3 {
				param1 = t.Param1
				param2 = t.Param2
				break
			} else {
				rn -= t.Param3
			}
		}
		return
	}

	if tl == nil {
		return nil, fmt.Errorf("task list is nil")
	}

	taskListCnf := configData.ConfigMgr().TaskListCnfById(tl.TaskListId)
	if taskListCnf == nil {
		return nil, fmt.Errorf("TaskList [%d] not found", tl.TaskListId)
	}

	var taskId int32
	taskList := taskListCnf.GetIncludeTask()
	rn := matrix.Random32(0, taskList.ChanceSum)
	for _, t := range taskList.ParamList {
		if rn <= t.Param2 {
			taskId = t.Param1
			break
		} else {
			rn -= t.Param2
		}
	}

	cnf := configData.ConfigMgr().TaskCnfById(taskId)
	if cnf == nil {
		return nil, fmt.Errorf("task[%v] config not found", taskId)
	}

	var taskObjectList *xlsxTable.TaskObjectList
	switch proto.TaskType(cnf.Kind) {
	case proto.TaskType_TaskTypeGetItem:
		taskObjectList = cnf.GetNeedItem()
	case proto.TaskType_TaskTypeUseItem:
		taskObjectList = cnf.GetUseItem()
	case proto.TaskType_TaskTypeKillMonster:
		taskObjectList = cnf.GetKillMonster()
	case proto.TaskType_TaskTypeMoveTo:
		taskObjectList = cnf.GetTargetPos()
	case proto.TaskType_TaskTypeQuiz:
		taskObjectList = cnf.GetQuiz()
	case proto.TaskType_TaskTypeOccupiedLand:
		taskObjectList = &xlsxTable.TaskObjectList{
			ChanceSum: 10000,
			ParamList: []xlsxTable.TaskObject{
				xlsxTable.TaskObject{
					Param1: cnf.RequestLand,
					Param3: 10000,
				},
			},
		}
	default:
		return nil, fmt.Errorf("task[%v] config invalid", taskId)
	}

	optionCnf := &dbData.TaskOptionCnf{TaskType: cnf.Kind}
	optionCnf.Param1, optionCnf.Param2 = randTaskF(*taskObjectList)
	opt := &dbData.TaskOption{
		OptionCnf: optionCnf,
		Rate:      0,
	}

	now := time.Now().UTC()
	pt := &dbData.Task{
		TaskId:    int32(taskId),
		TaskType:  cnf.Kind,
		Options:   []*dbData.TaskOption{opt},
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

	if proto.TaskListType_TaskListTypeRewarded == tlType && player.Level < TASK_REWARDED_LEVEL_MIN {
		return nil, nil
	}

	cnf := configData.ConfigMgr().TaskListCnfByLevel(int32(tlType), player.Level)
	if cnf == nil {
		serviceLog.Error("task list lv[%dv] config not found", player.Level)
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

	return curTl, nil
}

func (p *TaskModel) refreshPlayerTasks(userId int64, pt *dbData.PlayerTask) {
	if pt == nil {
		return
	}

	changed := false
	if dtl := pt.GetDailyTaskList(); dtl == nil {
		dtl, _ := p.randomTaskList(userId, proto.TaskListType_TaskListTypeDaily)
		if dtl != nil {
			pt.SetDailyTaskList(dtl)
			changed = true
			p.broadCastUpdateTaskListInfo(userId, proto.TaskListType_TaskListTypeDaily, dtl)
		}
	}

	if rtl := pt.GetRewardTaskList(); rtl == nil {
		rtl, _ := p.randomTaskList(userId, proto.TaskListType_TaskListTypeRewarded)
		if rtl != nil {
			pt.SetRewardTaskList(rtl)
			changed = true
			p.broadCastUpdateTaskListInfo(userId, proto.TaskListType_TaskListTypeRewarded, rtl)
		}
	}

	if changed {
		if err := gameDB.GetGameDB().Save(pt).Error; err != nil {
			serviceLog.Error(err.Error())
		}
	}
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

	pt := dbData.NewPlayerTask(userId, dtl, rtl)
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
	player, err := p.getPlayerSceneData(userId)
	if err != nil {
		return nil, err
	}
	if player.Level < TASK_REWARDED_LEVEL_MIN {
		return nil, fmt.Errorf("player level < 50, can't accept rewarded task")
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

func (p *TaskModel) UpGradeTaskProgress(
	userId int64,
	taskListKind proto.TaskListType,
	items []*proto.TaskOptionItem,
	pos *proto.TaskOptionMoveTo,
	quiz *proto.TaskOptionQuiz,
	monsterCid int32,
	usedItemCid, usedItemNum, optionLandNum int32,
) (*dbData.TaskList, error) {
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
		return nil, fmt.Errorf("task list cur task not found")
	}
	if tl.CurTask.IsFinish() {
		return nil, fmt.Errorf("task list cur task is finish")
	}

	upgrade := false

	for _, opt := range tl.CurTask.Options {
		if opt == nil || opt.OptionCnf == nil {
			continue
		}

		switch t := proto.TaskType(opt.OptionCnf.TaskType); t {
		case proto.TaskType_TaskTypeGetItem:
			var giveCount = opt.Rate
			nfts, err := p.getPlayerNFT(userId)
			if err != nil {
				return tl, err
			}
			for _, oit := range items {
				for _, nft := range nfts {
					if nft.Id != oit.NftId {
						continue
					}
					if int32(nft.Amount) < oit.Num {
						return tl, fmt.Errorf("not found [%v] item", oit.Num)
					}

					giveCount += oit.Num
					if giveCount > opt.OptionCnf.Param2 {
						return tl, fmt.Errorf("give too much item")
					}
				}
			}

			for _, oit := range items {
				grpcInvoke.BurnNFT(userId, oit.NftId, oit.Num)
				opt.Rate += oit.Num
				upgrade = true
			}
		case proto.TaskType_TaskTypeMoveTo:
			if pos == nil || pos.R != opt.OptionCnf.Param1 || pos.C != opt.OptionCnf.Param2 {
				return nil, fmt.Errorf("invalid pos")
			}
			opt.Rate = 1
			upgrade = true

		case proto.TaskType_TaskTypeQuiz:
			if quiz == nil || quiz.QuizType != opt.OptionCnf.Param1 {
				return nil, fmt.Errorf("invalid quiz")
			}
			opt.Rate = matrix.MinInt32(opt.Rate+quiz.QuizNum, opt.OptionCnf.Param2)
			upgrade = true

		case proto.TaskType_TaskTypeUseItem:
			if usedItemCid < 1 || usedItemNum < 1 {
				return tl, nil
			}
			for _, opt := range tl.CurTask.Options {
				if opt.OptionCnf == nil || opt.OptionCnf.Param1 != usedItemCid {
					continue
				}
				opt.Rate += usedItemNum
				upgrade = true
			}

		case proto.TaskType_TaskTypeKillMonster:
			if opt.OptionCnf != nil &&
				opt.OptionCnf.Param1 == monsterCid &&
				opt.Rate < opt.OptionCnf.Param2 {
				opt.Rate++
				upgrade = true
			}
		case proto.TaskType_TaskTypeOccupiedLand:
			opt.Rate = matrix.MinInt32(opt.Rate+optionLandNum, opt.OptionCnf.Param1)
			upgrade = true

		}
	}
	if !upgrade {
		return tl, nil
	}

	switch taskListKind {
	case proto.TaskListType_TaskListTypeDaily:
		pt.SetDailyTaskList(tl)
	case proto.TaskListType_TaskListTypeRewarded:
		pt.SetRewardTaskList(tl)
	}
	err = gameDB.GetGameDB().Save(pt).Error
	if err != nil {
		return tl, err
	}

	p.broadCastUpdateTaskListInfo(userId, taskListKind, tl)
	return tl, nil
}
