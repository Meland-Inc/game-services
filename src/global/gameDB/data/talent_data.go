package dbData

import (
	"encoding/json"
	"game-message-core/proto"
	"time"

	"github.com/Meland-Inc/game-services/src/common/time_helper"
)

// --------------- talent EXP ---------------
type TalentExp struct {
	TalentType proto.TalentType `json:"talentType"`
	CurExp     uint32           `json:"curExp"`   // 当前真是经验
	ExpCount   uint32           `json:"expCount"` // 总经验
}

func NewTalentExp(talentType proto.TalentType) *TalentExp {
	return &TalentExp{
		TalentType: talentType,
		CurExp:     0,
		ExpCount:   0,
	}
}

func (p *TalentExp) GetTalentType() proto.TalentType { return p.TalentType }
func (p *TalentExp) GetCurExp() uint32               { return p.CurExp }
func (p *TalentExp) GetExpCount() uint32             { return p.ExpCount }
func (p *TalentExp) ResetCurExp()                    { p.CurExp = p.ExpCount }
func (p *TalentExp) AddExp(exp uint32) {
	p.CurExp += exp
	p.ExpCount += exp
}
func (p *TalentExp) TakeExp(exp uint32) {
	p.CurExp -= exp
	if p.CurExp < 0 {
		p.CurExp = 0
	}
}

// --------------- talent EXP list ---------------
type TalentExpList struct {
	List []*TalentExp `json:"list"`
}

func NewTalentExpList() *TalentExpList {
	list := &TalentExpList{}
	list.List = append(list.List, NewTalentExp(proto.TalentType_Farming))
	list.List = append(list.List, NewTalentExp(proto.TalentType_Fighting))
	return list
}
func (p *TalentExpList) GetExpList() []*TalentExp { return p.List }
func (p *TalentExpList) TalentExpByType(talentType proto.TalentType) *TalentExp {
	for _, talentExp := range p.List {
		if talentExp.GetTalentType() == talentType {
			return talentExp
		}
	}

	talentExp := NewTalentExp(talentType)
	p.List = append(p.List, talentExp)
	return talentExp
}

// --------------- talent level ---------------
type TalentLevel struct {
	TalentType proto.TalentType `json:"talentType"`
	Level      uint32           `json:"level"`
}

func NewTalentLevel(talentType proto.TalentType) *TalentLevel {
	return &TalentLevel{TalentType: talentType, Level: 0}
}

func (p *TalentLevel) GetTalentType() proto.TalentType { return p.TalentType }
func (p *TalentLevel) GetLevel() uint32                { return p.Level }
func (p *TalentLevel) ResetLevel()                     { p.Level = 0 }
func (p *TalentLevel) SetLevel(lv uint32) {
	p.Level = lv
	if p.Level < 0 {
		p.Level = 0
	}
}

// --------------- talent level List ---------------
type TalentLevelList struct {
	List []*TalentLevel `json:"list"`
}

func NewTalentLevelList() *TalentLevelList {
	list := &TalentLevelList{}
	list.List = append(list.List, NewTalentLevel(proto.TalentType_Farming))
	list.List = append(list.List, NewTalentLevel(proto.TalentType_Fighting))
	return list
}
func (p *TalentLevelList) GetExpList() []*TalentLevel { return p.List }
func (p *TalentLevelList) TalentExpByType(talentType proto.TalentType) *TalentLevel {
	for _, talentLv := range p.List {
		if talentLv.GetTalentType() == talentType {
			return talentLv
		}
	}

	talentLv := NewTalentLevel(talentType)
	p.List = append(p.List, talentLv)
	return talentLv
}
func (p *TalentLevelList) ResetAll() {
	for _, talentLv := range p.List {
		talentLv.ResetLevel()
	}
}

// --------------- talent node   ---------------
type TalentNodeData struct {
	TalentType proto.TalentType `json:"talentType"`
	NodeIds    []uint32         `json:"nodeIds"`
}

func NewTalentNodeData(talentType proto.TalentType) *TalentNodeData {
	return &TalentNodeData{TalentType: talentType}
}

func (p *TalentNodeData) GetTalentType() proto.TalentType { return p.TalentType }
func (p *TalentNodeData) GetNodeIds() []uint32            { return p.NodeIds }
func (p *TalentNodeData) AddNode(nodeId uint32) {
	if nodeId <= 0 {
		return
	}
	p.NodeIds = append(p.NodeIds, nodeId)
}
func (p *TalentNodeData) Reset() {
	p.NodeIds = []uint32{}
}

// --------------- talent node list ---------------
type TalentNodeList struct {
	List []*TalentNodeData `json:"list"`
}

func NewTalentNodeList() *TalentNodeList {
	list := &TalentNodeList{}
	list.List = append(list.List, NewTalentNodeData(proto.TalentType_Farming))
	list.List = append(list.List, NewTalentNodeData(proto.TalentType_Fighting))
	return list
}
func (p *TalentNodeList) GetNodeList() []*TalentNodeData { return p.List }
func (p *TalentNodeList) TalentNodesByType(talentType proto.TalentType) *TalentNodeData {
	for _, node := range p.List {
		if node.GetTalentType() == talentType {
			return node
		}
	}

	node := NewTalentNodeData(talentType)
	p.List = append(p.List, node)
	return node
}
func (p *TalentNodeList) ResetAll() {
	for _, node := range p.List {
		node.Reset()
	}
}

// --------------- talent DB data ---------------
type TalentData struct {
	UserId    int64     `gorm:"primaryKey" json:"userId"`
	ExpJson   string    `json:"expJson"`
	LevelJson string    `json:"levelJson"`
	NodeJson  string    `json:"nodeJson"`
	CreatedAt time.Time `json:"createdAt"`
	UpdateAt  time.Time `json:"updateAt"`

	expData   *TalentExpList   `gorm:"-" json:"-"`
	levelData *TalentLevelList `gorm:"-" json:"-"`
	nodeData  *TalentNodeList  `gorm:"-" json:"-"`
}

func NewTalentData(userId int64) *TalentData {
	row := &TalentData{}
	row.UserId = userId
	row.CreatedAt = time_helper.NowUTC()
	row.UpdateAt = row.CreatedAt
	return row
}

func (p *TalentData) GetExpData() *TalentExpList {
	if p.expData == nil && len(p.ExpJson) > 2 {
		l := &TalentExpList{}
		err := json.Unmarshal([]byte(p.ExpJson), l)
		if err == nil {
			p.expData = l
		}
	}
	return p.expData
}

func (p *TalentData) SetExpData(data *TalentExpList) {
	p.expData = data
	p.ExpJson = ""
	if p.expData != nil {
		bs, err := json.Marshal(p.expData)
		if err == nil {
			p.ExpJson = string(bs)
		}
	}
}

func (p *TalentData) GetLevelData() *TalentLevelList {
	if p.levelData == nil && len(p.LevelJson) > 2 {
		l := &TalentLevelList{}
		err := json.Unmarshal([]byte(p.LevelJson), l)
		if err == nil {
			p.levelData = l
		}
	}
	return p.levelData
}

func (p *TalentData) SetLevelData(data *TalentLevelList) {
	p.levelData = data
	p.LevelJson = ""
	if p.levelData != nil {
		bs, err := json.Marshal(p.levelData)
		if err == nil {
			p.LevelJson = string(bs)
		}
	}
}

func (p *TalentData) GetNodeData() *TalentNodeList {
	if p.nodeData == nil && len(p.NodeJson) > 2 {
		l := &TalentNodeList{}
		err := json.Unmarshal([]byte(p.NodeJson), l)
		if err == nil {
			p.nodeData = l
		}
	}
	return p.nodeData
}

func (p *TalentData) SetNodeData(data *TalentNodeList) {
	p.nodeData = data
	p.NodeJson = ""
	if p.nodeData != nil {
		bs, err := json.Marshal(p.nodeData)
		if err == nil {
			p.NodeJson = string(bs)
		}
	}
}

// func (p *TalentData) ToProtoData() *proto.ItemBaseInfo { todo ...  }
// func (p *TalentData) ToGrpcData() *proto.ItemBaseInfo { todo ...  }
