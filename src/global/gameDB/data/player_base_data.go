package dbData

import (
	"encoding/json"
	"game-message-core/proto"
	"time"
)

type PlayerBaseData struct {
	UId         uint      `gorm:"primaryKey;autoIncrement" json:"uid,string"`
	UserId      int64     `json:"userId"`
	Name        string    `json:"name"`
	RoleId      int32     `json:"roleId"`
	RoleIcon    string    `json:"roleIcon"`
	FeatureJson string    `gorm:"type:text" json:"featureJson"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdateAt    time.Time `json:"updateAt"`

	Feature *proto.PlayerFeature `gorm:"-" json:"-"`
}

func (p *PlayerBaseData) SetFeature(feature *proto.PlayerFeature) error {
	if feature == nil {
		p.FeatureJson = ""
		p.Feature = nil
		return nil
	}

	bs, err := json.Marshal(feature)
	if err != nil {
		return err
	}
	p.FeatureJson = string(bs)
	p.Feature = feature
	return nil
}
func (p *PlayerBaseData) GetFeature() *proto.PlayerFeature {
	if p.Feature == nil && len(p.FeatureJson) >= 2 {
		feature := &proto.PlayerFeature{}
		err := json.Unmarshal([]byte(p.FeatureJson), feature)
		if err != nil {
			p.Feature = feature
		}
	}
	return p.Feature
}
func (p *PlayerBaseData) ToNetPlayerBaseData() *proto.PlayerBaseData {
	if p.UserId == 0 {
		return nil
	}
	return &proto.PlayerBaseData{
		UserId:   p.UserId,
		Name:     p.Name,
		RoleId:   p.RoleId,
		RoleIcon: p.RoleIcon,
		Feature:  p.GetFeature(),
	}
}
