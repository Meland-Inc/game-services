package taskModel

import (
	"fmt"
	"game-message-core/proto"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/global/gameDB"
	dbData "github.com/Meland-Inc/game-services/src/global/gameDB/data"
	"gorm.io/gorm"
)

const (
	TASK_PLAYER_CACHE_KEY = "task_player_cache_key_%d"
)

func (p *TaskModel) getPlayerTaskCacheKey(userId int64) string {
	return fmt.Sprintf(TASK_PLAYER_CACHE_KEY, userId)

}

func (p *TaskModel) GetPlayerTask(userId int64) (*dbData.PlayerTask, error) {
	cacheKey := p.getPlayerTaskCacheKey(userId)
	rv, err := p.cache.GetOrStore(
		cacheKey,
		func() (interface{}, error) {
			playerTask := &dbData.PlayerTask{}
			err := gameDB.GetGameDB().Where("user_id = ?", userId).First(playerTask).Error
			if err != nil {
				if err == gorm.ErrRecordNotFound {
					playerTask, err = p.InitPlayerTask(userId)
				} else {
					return nil, err
				}
			}
			p.tryRestTask(playerTask)

			if p.refreshPlayerTasksRate(userId, playerTask) {
				err := gameDB.GetGameDB().Save(playerTask).Error
				if err != nil {
					serviceLog.Error(err.Error())
				}
			}
			return playerTask, err
		},
		p.cacheTTL)

	if err != nil {
		return nil, err
	}

	p.cache.Touch(cacheKey, p.cacheTTL)
	pt := rv.(*dbData.PlayerTask)
	p.refreshPlayerTasks(userId, pt)
	return pt, nil
}

func (p *TaskModel) refreshPlayerTasks(userId int64, pt *dbData.PlayerTask) {
	if pt == nil {
		return
	}

	changed := false
	if dtl := pt.GetDailyTaskList(); dtl == nil {
		dtl, _ := p.initTaskList(userId, proto.TaskListType_TaskListTypeDaily)
		if dtl != nil {
			pt.SetDailyTaskList(dtl)
			changed = true
		}
	}

	if rtl := pt.GetRewardTaskList(); rtl == nil {
		rtl, _ := p.initTaskList(userId, proto.TaskListType_TaskListTypeRewarded)
		if rtl != nil {
			pt.SetRewardTaskList(rtl)
			changed = true
		}
	}

	if gtl := pt.GetGuideTaskList(); gtl == nil {
		gtl, _ := p.initTaskList(userId, proto.TaskListType_TaskListTypeGuide)
		if gtl != nil {
			pt.SetGuideTaskList(gtl)
			changed = true
		}
	} else if gtl.CurTask == nil {
		if nextTask, _ := p.getNextGuideTask(userId, gtl); nextTask != nil {
			gtl.CurTask = nextTask
			changed = true
		}
	}

	if changed {
		if err := gameDB.GetGameDB().Save(pt).Error; err != nil {
			serviceLog.Error(err.Error())
		}
	}
}

func (p *TaskModel) refreshPlayerTasksRate(userId int64, pt *dbData.PlayerTask) (upgrade bool) {
	sceneData, err := p.getPlayerSceneData(userId)
	if err != nil {
		return false
	}
	slotData, err := p.getPlayerSlotData(userId)
	if err != nil {
		return false
	}
	if dtl := pt.GetDailyTaskList(); dtl != nil && dtl.CurTask != nil {
		for _, opt := range dtl.CurTask.Options {
			if p.updateTaskOptionRate(sceneData, slotData, opt) {
				pt.SetDailyTaskList(dtl)
				upgrade = true
			}
		}
	}
	if rtl := pt.GetRewardTaskList(); rtl != nil && rtl.CurTask != nil {
		for _, opt := range rtl.CurTask.Options {
			if p.updateTaskOptionRate(sceneData, slotData, opt) {
				pt.SetRewardTaskList(rtl)
				upgrade = true
			}
		}
	}
	if gtl := pt.GetGuideTaskList(); gtl != nil && gtl.CurTask != nil {
		for _, opt := range gtl.CurTask.Options {
			if p.updateTaskOptionRate(sceneData, slotData, opt) {
				pt.SetGuideTaskList(gtl)
				upgrade = true
			}
		}
	}
	return upgrade
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
