package home_model

import (
	"fmt"
	"game-message-core/proto"
	"time"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/global/gameDB"
	dbData "github.com/Meland-Inc/game-services/src/global/gameDB/data"
	"github.com/Meland-Inc/game-services/src/global/serviceCnf"
	"github.com/Meland-Inc/game-services/src/global/userAgent"
	"gorm.io/gorm"
)

func (p *HomeModel) GetGranaryRows(userId int64) (rows []dbData.HomeGranary, err error) {
	err = gameDB.GetGameDB().Where("user_id = ? AND num > 0", userId).Find(&rows).Error
	return rows, err
}

func (p *HomeModel) GetGranaryRow(userId int64, itemCid int32) (row *dbData.HomeGranary) {
	row = &dbData.HomeGranary{}
	err := gameDB.GetGameDB().Where("user_id = ? AND item_cid = ?", userId, itemCid).First(row).Error
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			serviceLog.Error("GetGranaryRow err:%v", err)
		}
		return nil
	}
	return row
}

func (p *HomeModel) AddGranaryRecord(
	userId int64, itemCid int32, num, quality int32, upTm time.Time, lastPushUser int64, lastPushUserName string,
) error {
	row := dbData.NewHomeGranary(userId, itemCid, num, quality, lastPushUser, lastPushUserName, upTm)
	return gameDB.GetGameDB().Save(row).Error
}

func (p *HomeModel) TryAddGranaryRecord(
	userId int64, itemCid int32, num, quality int32, upTm time.Time, lastPushUser int64, lastPushUserName string,
) error {
	row := p.GetGranaryRow(userId, itemCid)
	if row == nil {
		return p.AddGranaryRecord(userId, itemCid, num, quality, upTm, lastPushUser, lastPushUserName)
	}

	if row != nil && row.UpdateAt.UnixMilli() >= upTm.UnixMilli() {
		return nil
	}
	row.Num += num
	row.UpdateAt = upTm
	row.Quality = quality
	row.LastPushUserId = lastPushUser
	row.LastPushUserName = lastPushUserName
	return gameDB.GetGameDB().Save(row).Error
}

func (p *HomeModel) RemoveGranaryRecord(userId int64, itemCid int32) error {
	row := p.GetGranaryRow(userId, itemCid)
	if row == nil {
		return fmt.Errorf("granary record not found for user[%d],cid[%d]", userId, itemCid)
	}
	row.Num = 0
	row.Quality = 0
	row.LastPushUserId = 0
	row.LastPushUserName = ""
	return gameDB.GetGameDB().Delete(row).Error
}

func (p *HomeModel) ClearGranaryRecord(userId int64) error {
	return gameDB.GetGameDB().Model(&dbData.HomeGranary{}).Where("user_id = ?", userId).Updates(
		map[string]interface{}{
			"num":                 0,
			"quality":             0,
			"last_push_user_id":   0,
			"last_push_user_name": "",
		},
	).Error
}

func (p *HomeModel) BroadCastUpAllGranary(userId int64) {
	rows, err := p.GetGranaryRows(userId)
	if err != nil {
		return
	}

	pbItems := []*proto.ItemBaseInfo{}
	for _, row := range rows {
		pbItems = append(pbItems, row.ToProtoData())
	}

	msg := &proto.Envelope{
		Type: proto.EnvelopeType_BroadCastGranaryUpdate,
		Payload: &proto.Envelope_BroadCastGranaryUpdateResponse{
			BroadCastGranaryUpdateResponse: &proto.BroadCastGranaryUpdateResponse{
				Items: pbItems,
			},
		},
	}
	err = userAgent.SendToPlayer(serviceCnf.GetInstance().AppId, userId, msg)
	if err != nil {
		serviceLog.Error(err.Error())
	}
}

func (p *HomeModel) BroadCastUpGranaryItem(row dbData.HomeGranary) {
	msg := &proto.Envelope{
		Type: proto.EnvelopeType_BroadCastUpGranaryItem,
		Payload: &proto.Envelope_BroadCastUpGranaryItemResponse{
			BroadCastUpGranaryItemResponse: &proto.BroadCastUpGranaryItemResponse{
				Items: row.ToProtoData(),
			},
		},
	}
	err := userAgent.SendToPlayer(serviceCnf.GetInstance().AppId, row.UserId, msg)
	if err != nil {
		serviceLog.Error(err.Error())
	}
}
