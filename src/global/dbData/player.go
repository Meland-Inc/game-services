package dbData

import (
	"encoding/json"
	"game-message-core/proto"
	"time"
)

type PlayerRow struct {
	UId         uint      `gorm:"primaryKey;autoIncrement" json:"uid,string"`
	UserId      int64     `json:"userId"`
	Name        string    `json:"name"`
	RoleId      int32     `json:"roleId"`
	RoleIcon    string    `json:"roleIcon"`
	FeatureJson string    `gorm:"type:text" json:"featureJson"`
	Level       int32     `json:"level"`
	Exp         int32     `json:"exp"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdateAt    time.Time `json:"updateAt"`
	LastLogin   time.Time `json:"lastLogin"`

	Feature *proto.PlayerFeature `gorm:"-" json:"-"`
}

func (p *PlayerRow) SetFeature(feature *proto.PlayerFeature) error {
	bs, err := json.Marshal(feature)
	if err != nil {
		return err
	}

	p.FeatureJson = string(bs)
	p.Feature = feature
	return err
}
func (p *PlayerRow) GetFeature() *proto.PlayerFeature {
	if p.Feature == nil && len(p.FeatureJson) >= 2 {
		feature := &proto.PlayerFeature{}
		err := json.Unmarshal([]byte(p.FeatureJson), feature)
		if err != nil {
			p.Feature = feature
		}
	}
	return p.Feature
}
