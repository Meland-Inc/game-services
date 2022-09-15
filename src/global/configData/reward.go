package configData

import (
	xlsxTable "game-message-core/xlsxTableData"
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
