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
	list.List = append(list.List, NewTalentExp(proto.TalentType_Battle))
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
	list.List = append(list.List, NewTalentLevel(proto.TalentType_Battle))
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

// --------------- talent node data ---------------
type TalentNodeData struct {
	NodeId uint32 `json:"nodeId"`
	Level  uint32 `json:"level"`
}

func NewTalentNodeData(nodeId, lv uint32) *TalentNodeData {
	return &TalentNodeData{
		NodeId: nodeId,
		Level:  lv,
	}
}

// --------------- talent node   ---------------
type TalentTree struct {
	TalentType proto.TalentType  `json:"talentType"`
	Nodes      []*TalentNodeData `json:"nodes"`
}

func NewTalentTree(talentType proto.TalentType) *TalentTree {
	return &TalentTree{TalentType: talentType}
}

func (p *TalentTree) GetTalentType() proto.TalentType { return p.TalentType }
func (p *TalentTree) GetNodes() []*TalentNodeData     { return p.Nodes }
func (p *TalentTree) NodeById(nodeId uint32) *TalentNodeData {
	for _, node := range p.Nodes {
		if node.NodeId == nodeId {
			return node
		}
	}
	return nil
}
func (p *TalentTree) AddNode(node *TalentNodeData) {
	if node == nil || node.NodeId <= 0 || node.Level < 1 {
		return
	}
	for _, n := range p.Nodes {
		if n.NodeId == node.NodeId {
			return
		}
	}
	p.Nodes = append(p.Nodes, node)
}
func (p *TalentTree) Reset() {
	p.Nodes = []*TalentNodeData{}
}

// --------------- talent node list ---------------
type TalentTreeList struct {
	List []*TalentTree `json:"list"`
}

func NewTalentNodeList() *TalentTreeList {
	list := &TalentTreeList{}
	list.List = append(list.List, NewTalentTree(proto.TalentType_Farming))
	list.List = append(list.List, NewTalentTree(proto.TalentType_Battle))
	return list
}
func (p *TalentTreeList) GetTalentTreeList() []*TalentTree { return p.List }
func (p *TalentTreeList) TalentTreeByType(talentType proto.TalentType) *TalentTree {
	for _, node := range p.List {
		if node.GetTalentType() == talentType {
			return node
		}
	}

	tree := NewTalentTree(talentType)
	p.List = append(p.List, tree)
	return tree
}
func (p *TalentTreeList) ResetAll() {
	for _, node := range p.List {
		node.Reset()
	}
}

// --------------- talent DB data ---------------
type TalentData struct {
	UserId         int64     `gorm:"primaryKey" json:"userId"`
	ExpJson        string    `json:"expJson"`
	LevelJson      string    `json:"levelJson"`
	TalentTreeJson string    `json:"talentTreeJson"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdateAt       time.Time `json:"updateAt"`

	expData        *TalentExpList   `gorm:"-" json:"-"`
	levelData      *TalentLevelList `gorm:"-" json:"-"`
	talentTreeData *TalentTreeList  `gorm:"-" json:"-"`
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

func (p *TalentData) GetTalentTreeList() *TalentTreeList {
	if p.talentTreeData == nil && len(p.TalentTreeJson) > 2 {
		l := &TalentTreeList{}
		err := json.Unmarshal([]byte(p.TalentTreeJson), l)
		if err == nil {
			p.talentTreeData = l
		}
	}
	return p.talentTreeData
}

func (p *TalentData) SetTalentTreeList(data *TalentTreeList) {
	p.talentTreeData = data
	p.TalentTreeJson = ""
	if p.talentTreeData != nil {
		bs, err := json.Marshal(p.talentTreeData)
		if err == nil {
			p.TalentTreeJson = string(bs)
		}
	}
}

// func (p *TalentData) ToProtoData() *proto.ItemBaseInfo { todo ...  }
// func (p *TalentData) ToGrpcData() *proto.ItemBaseInfo { todo ...  }
