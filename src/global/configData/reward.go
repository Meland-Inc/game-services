package configData

import (
	"fmt"
	"game-message-core/proto"
	xlsxTable "game-message-core/xlsxTableData"

	"github.com/Meland-Inc/game-services/src/common/matrix"
)

func (mgr *ConfigDataManager) initReward() error {
	mgr.rewardCnf = make(map[int32]xlsxTable.RewardTableRow)

	rows := []xlsxTable.RewardTableRow{}
	err := mgr.configDb.Find(&rows).Error
	if err != nil {
		return err
	}

	for _, row := range rows {
		mgr.rewardCnf[row.RewardId] = row
	}

	return nil
}

func (mgr *ConfigDataManager) AllReward() map[int32]xlsxTable.RewardTableRow {
	return mgr.rewardCnf
}

func (mgr *ConfigDataManager) RewardCnfById(id int32) *xlsxTable.RewardTableRow {
	cnf, exist := mgr.rewardCnf[id]
	if !exist {
		return nil
	}
	return &cnf
}

func RandomRewardItems(rewardId int32) (rewards []*proto.ItemBaseInfo, err error) {
	if rewardId == 0 {
		return rewards, nil
	}

	cnf := ConfigMgr().RewardCnfById(rewardId)
	if cnf == nil {
		return rewards, fmt.Errorf("rewardId [%d] not found", rewardId)
	}

	timeList, err := cnf.GetRewardTimeList()
	if err != nil {
		return rewards, err
	}
	rewardItemList, err := cnf.GetRewardList()
	if err != nil {
		return rewards, err
	}

	// 随机奖励次数
	var rewardTimes int32
	rnd := matrix.Random32(0, timeList.TotalWeight)
	for _, ts := range timeList.RewardTimes {
		if rnd <= ts.Weight {
			rewardTimes = ts.Time
			break
		}
		rnd -= ts.Weight
	}

	// 随机奖励物品
	for i := 0; i < int(rewardTimes); i++ {
		rnd := matrix.Random32(0, rewardItemList.TotalWeight)
		for _, rs := range rewardItemList.Rewards {
			if rnd <= rs.Weight {
				rewards = append(rewards, &proto.ItemBaseInfo{
					Cid:     rs.Cid,
					Num:     rs.Num,
					Quality: rs.Quantity,
				})
				break
			}
			rnd -= rs.Weight
		}
	}

	return rewards, nil
}
