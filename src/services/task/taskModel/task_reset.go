package taskModel

import (
	"game-message-core/proto"
	"time"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/global/gameDB"
	dbData "github.com/Meland-Inc/game-services/src/global/gameDB/data"
)

// 悬赏任务最近的一次重置的时间节点
func rewardTaskLastResetTime() time.Time {
	now := time.Now().UTC()
	offset := int(-(now.Weekday() + 1))
	if offset < -6 {
		offset += 7
	}
	lastSaturday := now.AddDate(0, 0, offset)

	return time.Date(
		lastSaturday.Year(),
		lastSaturday.Month(),
		lastSaturday.Day(),
		0, 0, 0, 0,
		lastSaturday.Location(),
	)
}

func (this *TaskModel) resetPlayerTask(now time.Time, pt *dbData.PlayerTask, resetDTl, resetRTL, broadCast bool) {
	if resetDTl {
		dtl, err := this.randomTaskList(pt.UserId, proto.TaskListType_TaskListTypeDaily)
		if err == nil {
			pt.SetDailyTaskList(dtl)
			if broadCast {
				this.broadCastUpdateTaskListInfo(pt.UserId, proto.TaskListType_TaskListTypeDaily, dtl)
			}
		}
	}

	if resetRTL {
		rtl := pt.GetRewardTaskList()
		if rtl == nil || !rtl.Doing {
			rtl, _ = this.randomTaskList(pt.UserId, proto.TaskListType_TaskListTypeRewarded)
		}
		if rtl != nil {
			rtl.CanReceive = true
			rtl.ResetAt = rewardTaskLastResetTime()
		}
		pt.SetRewardTaskList(rtl)
		if broadCast {
			this.broadCastUpdateTaskListInfo(pt.UserId, proto.TaskListType_TaskListTypeRewarded, rtl)
		}
	}

	if resetDTl || resetRTL {
		err := gameDB.GetGameDB().Save(pt).Error
		if err != nil {
			serviceLog.Error("resetPlayerTask DB SAVE err: %v", err)
		}
	}
}

// 从db 获取玩家任务信息后需要检测 任务是否需要重置
func (this *TaskModel) tryRestTask(pt *dbData.PlayerTask) {
	now := time.Now().UTC()
	needResetDtl := false
	if dtl := pt.GetDailyTaskList(); dtl != nil {
		if dtl.ResetAt.Day() != now.Day() ||
			dtl.ResetAt.Month() != now.Month() ||
			dtl.ResetAt.Year() != now.Year() {
			needResetDtl = true
		}
	}

	needResetRtl := false
	if rtl := pt.GetRewardTaskList(); rtl != nil {
		preSaturday := rewardTaskLastResetTime()
		if preSaturday.Day() != rtl.ResetAt.Day() ||
			preSaturday.Month() != rtl.ResetAt.Month() ||
			preSaturday.Year() != rtl.ResetAt.Year() {
			needResetRtl = true
		}
	}

	this.resetPlayerTask(now, pt, needResetDtl, needResetRtl, false)
}

func (this *TaskModel) checkAndResetPlayerTask(now time.Time) {
	needResetRewardTl := now.Weekday() == time.Saturday
	needResetDailyTl := true

	cacheObjects := this.cache.Items()
	for _, it := range cacheObjects {
		for _, v := range it {
			if pt, ok := v.Object.(*dbData.PlayerTask); ok {
				this.resetPlayerTask(now, pt, needResetDailyTl, needResetRewardTl, true)
			}
		}
	}
}
