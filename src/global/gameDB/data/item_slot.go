package dbData

import (
	"encoding/json"
	"game-message-core/proto"
	"time"

	message "github.com/Meland-Inc/game-services/src/global/web3Message"
)

type SlotList struct {
	SlotList []*message.PlayerItemSlot
}

type ItemSlot struct {
	UId       uint      `gorm:"primaryKey;autoIncrement" json:"uid,string"`
	UserId    int64     `json:"userId"`
	SlotsJson string    `db:"slots"`
	CreatedAt time.Time `json:"createdAt"`
	UpdateAt  time.Time `json:"updateAt"`

	Slots *SlotList `gorm:"-" json:"-"`
}

func (this *ItemSlot) GetSlotList() *SlotList {
	if this.Slots == nil && len(this.SlotsJson) >= 2 {
		slots := &SlotList{}
		json.Unmarshal([]byte(this.SlotsJson), slots)
		this.Slots = slots
	}
	return this.Slots
}

func (this *ItemSlot) InitSlotList() error {
	slotList := &SlotList{}
	posBegin := int(proto.AvatarPosition_AvatarPositionHead)
	posEnd := int(proto.AvatarPosition_AvatarPositionWeapon)
	for i := posBegin; i <= posEnd; i++ {
		slotList.SlotList = append(
			slotList.SlotList,
			&message.PlayerItemSlot{Position: i, Level: 1})
	}
	return this.setSlots(slotList)
}

func (this *ItemSlot) setSlots(sockets *SlotList) error {
	bs, err := json.Marshal(sockets)
	if err != nil {
		return err
	}
	this.Slots = sockets
	this.SlotsJson = string(bs)
	return nil
}

func (this *ItemSlot) SetSlotLevel(pos proto.AvatarPosition, lv int32) {
	list := this.GetSlotList()
	for idx, slot := range list.SlotList {
		if slot.Position == int(pos) {
			list.SlotList[idx].Level = int(lv)
			break
		}
	}
	this.setSlots(list)
}
