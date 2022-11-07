package dbData

import (
	"encoding/json"
	"game-message-core/proto"
	"time"

	"github.com/Meland-Inc/game-services/src/global/configData"
)

type TaskOptionCnf struct {
	TaskOptionType int32 `json:"taskOptionType"` // 等价于 proto.TaskOptionType
	Param1         int32 `json:"param1"`
	Param2         int32 `json:"param2"`
	Param3         int32 `json:"param3"`
	Param4         int32 `json:"param4"`
}

func (p *TaskOptionCnf) ToPbData() *proto.TaskOptionCnf {
	pbSetting := &proto.TaskOptionCnf{
		Kind: proto.TaskOptionType(p.TaskOptionType),
	}

	switch pbSetting.Kind {
	case proto.TaskOptionType_HandInItem:
		pbSetting.Data = &proto.TaskOptionCnf_HandInItem{
			HandInItem: &proto.TaskOptionItem{
				ItemCid: p.Param1,
				Num:     p.Param2,
			},
		}
	case proto.TaskOptionType_UseItem:
		pbSetting.Data = &proto.TaskOptionCnf_UseItem{
			UseItem: &proto.TaskOptionItem{
				ItemCid: p.Param1,
				Num:     p.Param2,
			},
		}
	case proto.TaskOptionType_GetItem:
		pbSetting.Data = &proto.TaskOptionCnf_GetItem{
			GetItem: &proto.TaskOptionItem{
				ItemCid: p.Param1,
				Num:     p.Param2,
			},
		}
	case proto.TaskOptionType_KillMonster:
		pbSetting.Data = &proto.TaskOptionCnf_KillMonster{
			KillMonster: &proto.TaskOptionKillMonster{
				MonCid: p.Param1,
				Num:    p.Param2,
			},
		}
	case proto.TaskOptionType_UserLevel:
		pbSetting.Data = &proto.TaskOptionCnf_UserLevel{
			UserLevel: p.Param1,
		}
	case proto.TaskOptionType_TargetSlotLevel:
		pbSetting.Data = &proto.TaskOptionCnf_TargetSlotLevel{
			TargetSlotLevel: &proto.TaskOptionTargetSlotLevel{
				SlotPos: p.Param1,
				Level:   p.Param2,
			},
		}
	case proto.TaskOptionType_SlotLevelCount:
		pbSetting.Data = &proto.TaskOptionCnf_SlotLevelCount{
			SlotLevelCount: &proto.TaskOptionSlotLevelCount{
				SlotLevel: p.Param1,
				SlotCount: p.Param2,
			},
		}
	case proto.TaskOptionType_CraftSkillLevel:
		pbSetting.Data = &proto.TaskOptionCnf_CraftSkillLevel{
			CraftSkillLevel: &proto.TaskOptionCraftSkillLevel{
				SkillId: p.Param1,
				Level:   p.Param2,
			},
		}
	case proto.TaskOptionType_UseRecipe:
		pbSetting.Data = &proto.TaskOptionCnf_UseRecipe{
			UseRecipe: &proto.TaskOptionUseRecipe{
				RecipeId: p.Param1,
				Times:    p.Param2,
			},
		}
	case proto.TaskOptionType_RecipeUseCount:
		pbSetting.Data = &proto.TaskOptionCnf_RecipeUseCount{
			RecipeUseCount: p.Param1,
		}
	case proto.TaskOptionType_TaskCount:
		pbSetting.Data = &proto.TaskOptionCnf_FinishTaskCount{
			FinishTaskCount: &proto.TaskOptionFinishTaskCount{
				Kind:  proto.TaskListType(p.Param1),
				Count: p.Param2,
			},
		}
	case proto.TaskOptionType_TaskListTypeCount:
		pbSetting.Data = &proto.TaskOptionCnf_FinishTaskListCount{
			FinishTaskListCount: &proto.TaskOptionFinishTaskListCount{
				Kind:  proto.TaskListType(p.Param1),
				Count: p.Param2,
			},
		}
	case proto.TaskOptionType_TargetPosition:
		pbSetting.Data = &proto.TaskOptionCnf_TargetPosition{
			TargetPosition: &proto.TaskOptionTargetPosition{
				X:          p.Param1,
				Y:          p.Param2,
				Z:          p.Param3,
				DistOffset: p.Param4,
			},
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
	return &proto.TaskOption{
		Rate:      p.Rate,
		OptionCnf: p.OptionCnf.ToPbData(),
	}
}

func (p *TaskOption) IsFinish() bool {
	switch proto.TaskOptionType(p.OptionCnf.TaskOptionType) {
	case proto.TaskOptionType_HandInItem,
		proto.TaskOptionType_UseItem,
		proto.TaskOptionType_GetItem:
		if p.Rate < p.OptionCnf.Param2 {
			return false
		}
	case proto.TaskOptionType_KillMonster:
		if p.Rate < p.OptionCnf.Param2 {
			return false
		}
	case proto.TaskOptionType_UserLevel:
		if p.Rate < p.OptionCnf.Param1 {
			return false
		}
	case proto.TaskOptionType_TargetSlotLevel:
		if p.Rate < p.OptionCnf.Param2 {
			return false
		}
	case proto.TaskOptionType_SlotLevelCount:
		if p.Rate < p.OptionCnf.Param2 {
			return false
		}
	case proto.TaskOptionType_CraftSkillLevel:
		if p.Rate < p.OptionCnf.Param2 {
			return false
		}
	case proto.TaskOptionType_UseRecipe:
		if p.Rate < p.OptionCnf.Param2 {
			return false
		}
	case proto.TaskOptionType_RecipeUseCount:
		if p.Rate < p.OptionCnf.Param1 {
			return false
		}
	case proto.TaskOptionType_TaskCount:
		if p.Rate < p.OptionCnf.Param2 {
			return false
		}
	case proto.TaskOptionType_TaskListTypeCount:
		if p.Rate < p.OptionCnf.Param2 {
			return false
		}
	case proto.TaskOptionType_TargetPosition:
		if p.Rate < 1 { // 客户端上报坐标 检测通过后 标记为1
			return false
		}
	}
	return true
}

type Task struct {
	TaskId    int32         `json:"taskId"`
	Options   []*TaskOption `json:"options"`
	CreatedAt time.Time     `json:"createdAt"`
	UpdatedAt time.Time     `json:"updatedAt"`
}

func (p *Task) ToPbData() *proto.Task {
	pt := &proto.Task{
		TaskId:       p.TaskId,
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
		if opt == nil {
			continue
		}
		if !opt.IsFinish() {
			return false
		}
	}
	return true
}

type TaskList struct {
	TaskListId    int32     `json:"taskListId"`
	TaskListType  int32     `json:"taskListType"` // 等价于 proto.TaskListType
	CanReceive    bool      `json:"canReceive"`   // 任务链是否可以在接取 子任务
	Doing         bool      `json:"doing"`        // 任务链是否正在进行
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
	case proto.TaskListType_TaskListTypeGuide:
		return false
	}
	return
}

type PlayerTask struct {
	Id                   uint      `gorm:"primaryKey;autoIncrement" json:"id,string"`
	UserId               int64     `gorm:"not null" json:"userId"`
	DailyTaskListJson    string    `gorm:"type:text" json:"dailyTaskListJson"`
	RewardedTaskListJson string    `gorm:"type:text" json:"rewardedTaskListJson"`
	GuideTaskListJson    string    `gorm:"type:text" json:"guideTaskList"`
	CreatedAt            time.Time `gorm:"not null" json:"createdAt"`

	DailyTaskList    *TaskList `gorm:"-" json:"-"`
	RewardedTaskList *TaskList `gorm:"-" json:"-"`
	GuideTaskList    *TaskList `gorm:"-" json:"-"`
}

func NewPlayerTask(userId int64,
	dailyTaskList *TaskList,
	rewardedTaskList *TaskList,
	guideTaskList *TaskList,
) *PlayerTask {
	pt := &PlayerTask{
		UserId:    userId,
		CreatedAt: time.Now().UTC(),
	}
	pt.SetDailyTaskList(dailyTaskList)
	pt.SetRewardTaskList(rewardedTaskList)
	pt.SetGuideTaskList(guideTaskList)
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

func (pt *PlayerTask) GetGuideTaskList() *TaskList {
	if pt.GuideTaskList == nil && len(pt.GuideTaskListJson) > 2 {
		dtl := &TaskList{}
		err := json.Unmarshal([]byte(pt.GuideTaskListJson), dtl)
		if err == nil {
			pt.GuideTaskList = dtl
		}
	}
	return pt.GuideTaskList
}

func (pt *PlayerTask) SetGuideTaskList(tl *TaskList) {
	pt.GuideTaskList = tl
	pt.GuideTaskListJson = ""
	if pt.GuideTaskList != nil {
		bs, err := json.Marshal(pt.GuideTaskList)
		if err == nil {
			pt.GuideTaskListJson = string(bs)
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
	if gtl := pt.GetGuideTaskList(); gtl != nil {
		pbPt.TaskLists = append(pbPt.TaskLists, gtl.ToPbData())
	}
	return
}
