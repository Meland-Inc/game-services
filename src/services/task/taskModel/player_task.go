package taskModel

import (
	"fmt"

	"github.com/Meland-Inc/game-services/src/common/time_helper"
	"github.com/Meland-Inc/game-services/src/global/gameDB"
	dbData "github.com/Meland-Inc/game-services/src/global/gameDB/data"
	"gorm.io/gorm"
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

func (p *TaskModel) taskTick(curMs int64) error {
	now := time_helper.NowUTC()
	if now.Hour() != 0 || now.Minute() != 0 {
		return nil
	}
	p.checkAndResetPlayerTask(now)
	return nil
}
