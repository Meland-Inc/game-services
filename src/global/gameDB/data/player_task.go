package dbData

import (
	"encoding/json"
	"game-message-core/proto"
	"time"

	"github.com/Meland-Inc/game-services/src/global/configData"
)

type PlayerTask struct {
	Id                   uint      `gorm:"primaryKey;autoIncrement" json:"id,string"`
	UserId               int64     `gorm:"not null" json:"userId"`
	DailyTaskListJson    string    `gorm:"type:text" json:"dailyTaskListJson"`
	RewardedTaskListJson string    `gorm:"type:text" json:"rewardedTaskListJson"`
	CreatedAt            time.Time `gorm:"not null" json:"createdAt"`
	DailyTaskList        *TaskList `gorm:"-" json:"-"`
	RewardedTaskList     *TaskList `gorm:"-" json:"-"`
}

func NewPlayerTask(playerId int64, dailyTaskList, rewardedTaskList *TaskList) *PlayerTask {
	pt := &PlayerTask{
		UserId:    playerId,
		CreatedAt: time.Now().UTC(),
	}
	pt.SetDailyTaskList(dailyTaskList)
	pt.SetRewardTaskList(rewardedTaskList)
	return pt
}

func (pt *PlayerTask) GetDailyTaskList() *TaskList {
	if pt.DailyTaskList == nil && len(pt.DailyTaskListJson) > 2 {
		dtl := &TaskList{}
		err := json.Unmarshal([]byte(pt.DailyTaskListJson), dtl)
		if err == nil {
			pt.DailyTaskList = dtl
		}
	}

	return pt.DailyTaskList
}

func (pt *PlayerTask) SetDailyTaskList(tl *TaskList) {
	pt.DailyTaskList = tl
	pt.DailyTaskListJson = ""
	if pt.DailyTaskList != nil {
		bs, err := json.Marshal(pt.DailyTaskList)
		if err == nil {
			pt.DailyTaskListJson = string(bs)
		}
	}
}

func (pt *PlayerTask) GetRewardTaskList() *TaskList {
	if pt.RewardedTaskList == nil && len(pt.RewardedTaskListJson) > 2 {
		dtl := &TaskList{}
		err := json.Unmarshal([]byte(pt.RewardedTaskListJson), dtl)
		if err == nil {
			pt.RewardedTaskList = dtl
		}
	}

	return pt.RewardedTaskList
}

func (pt *PlayerTask) SetRewardTaskList(tl *TaskList) {
	pt.RewardedTaskList = tl
	pt.RewardedTaskListJson = ""
	if pt.RewardedTaskList != nil {
		bs, err := json.Marshal(pt.RewardedTaskList)
		if err == nil {
			pt.RewardedTaskListJson = string(bs)
		}
	}
}

func (pt *PlayerTask) ToProtoData() (pbPt *proto.PlayerTask) {
	pbPt = &proto.PlayerTask{}
	if dtl := pt.GetDailyTaskList(); dtl != nil {
		pbPt.TaskLists = append(pbPt.TaskLists, dtl.ToPbData())
	}
	if rtl := pt.GetRewardTaskList(); rtl != nil {
		pbPt.TaskLists = append(pbPt.TaskLists, rtl.ToPbData())
	}
	return
}

type TaskOptionCnf struct {
	TaskType int32 `json:"taskType"` // 等价于 proto.TaskType_TaskType
	Param1   int32 `json:"param1"`
	Param2   int32 `json:"param2"`
	Param3   int32 `json:"param3"`
}

func (p *TaskOptionCnf) ToPbData() *proto.TaskOptionCnf {
	pbSetting := &proto.TaskOptionCnf{
		Kind: proto.TaskType(p.TaskType),
	}

	switch pbSetting.Kind {
	case proto.TaskType_TaskTypeGetItem, proto.TaskType_TaskTypeUseItem:
		pbSetting.Data = &proto.TaskOptionCnf_Item{
			Item: &proto.TaskOptionItem{
				ItemCid: p.Param1,
				Num:     p.Param2,
			},
		}
	case proto.TaskType_TaskTypeKillMonster:
		pbSetting.Data = &proto.TaskOptionCnf_MonInfo{
			MonInfo: &proto.TaskOptionMonster{
				MonCid: p.Param1,
				Num:    p.Param2,
			},
		}
	case proto.TaskType_TaskTypeMoveTo:
		pbSetting.Data = &proto.TaskOptionCnf_TarPos{
			TarPos: &proto.TaskOptionMoveTo{
				R: p.Param1,
				C: p.Param2,
			},
		}
	case proto.TaskType_TaskTypeQuiz:
		pbSetting.Data = &proto.TaskOptionCnf_QuizInfo{
			QuizInfo: &proto.TaskOptionQuiz{
				QuizType: p.Param1,
				QuizNum:  p.Param2,
			},
		}
	case proto.TaskType_TaskTypeOccupiedLand:
		pbSetting.Data = &proto.TaskOptionCnf_Num{
			Num: p.Param1,
		}

	}

	return pbSetting
}

type TaskOption struct {
	OptionCnf *TaskOptionCnf `json:"optionCnf"`
	Rate      int32          `json:"rate"`
}

func (p *TaskOption) ToPbData() *proto.TaskOption {
	if p.OptionCnf == nil {
		return nil
	}
	to := &proto.TaskOption{
		Rate:      p.Rate,
		OptionCnf: p.OptionCnf.ToPbData(),
	}

	return to
}

type Task struct {
	TaskId    int32         `json:"taskId"`
	TaskType  int32         `json:"taskType"` // 等价于 proto.TaskType_TaskType
	Options   []*TaskOption `json:"options"`
	CreatedAt time.Time     `json:"createdAt"`
	UpdatedAt time.Time     `json:"updatedAt"`
}

func (p *Task) ToPbData() *proto.Task {
	pt := &proto.Task{
		TaskId:       p.TaskId,
		TaskKind:     proto.TaskType(p.TaskType),
		CreatedAtSec: p.CreatedAt.Unix(),
	}
	for _, opt := range p.Options {
		if pb := opt.ToPbData(); pb != nil {
			pt.Options = append(pt.Options, pb)
		}
	}

	return pt
}

func (p *Task) IsFinish() bool {
	for _, opt := range p.Options {
		if opt == nil || opt.OptionCnf == nil {
			continue
		}

		switch proto.TaskType(opt.OptionCnf.TaskType) {
		case proto.TaskType_TaskTypeGetItem, proto.TaskType_TaskTypeUseItem:
			if opt.Rate < opt.OptionCnf.Param2 {
				return false
			}
		case proto.TaskType_TaskTypeKillMonster:
			if opt.Rate < opt.OptionCnf.Param2 {
				return false
			}
		case proto.TaskType_TaskTypeMoveTo:
			if opt.Rate == 0 {
				return false
			}
		case proto.TaskType_TaskTypeQuiz:
			if opt.Rate < opt.OptionCnf.Param2 {
				return false
			}
		case proto.TaskType_TaskTypeOccupiedLand:
			if opt.Rate < opt.OptionCnf.Param1 {
				return false
			}
		}
	}

	return true
}

type TaskList struct {
	TaskListId    int32     `json:"taskListId"`
	TaskListType  int32     `json:"taskListType"` // 等价于 proto.TaskListType
	CanReceive    bool      `json:"canReceive"`
	Doing         bool      `json:"doing"` // 任务链是否正在进行
	Rate          int32     `json:"rate"`
	CurTask       *Task     `json:"curTask"`
	ReceiveReward int32     `json:"receiveReward"` // 已经领取了的阶段奖励 0：未领取
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
	ResetAt       time.Time `json:"resetAt"`
}

func (p *TaskList) ToPbData() *proto.TaskList {
	pbTl := &proto.TaskList{
		Id:            p.TaskListId,
		Kind:          proto.TaskListType(p.TaskListType),
		CanReceive:    p.CanReceive,
		Doing:         p.Doing,
		Rate:          p.Rate,
		ReceiveReward: p.ReceiveReward,
	}
	if p.CurTask != nil {
		pbTl.CurTask = p.CurTask.ToPbData()
	}
	return pbTl
}

func (p *TaskList) IsFinish() (finish bool) {
	switch proto.TaskListType(p.TaskListType) {
	case proto.TaskListType_TaskListTypeDaily:
		finish = p.Rate >= configData.ConfigMgr().DailyTaskRateLimit()
	case proto.TaskListType_TaskListTypeRewarded:
		finish = p.Rate >= configData.ConfigMgr().RewardTaskListRateMax()
	}
	return
}
