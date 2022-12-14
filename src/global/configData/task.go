package configData

import (
	xlsxTable "game-message-core/xlsxTableData"
)

func (mgr *ConfigDataManager) initTask() error {
	mgr.taskCnf = make(map[int32]xlsxTable.TaskTableRow)

	rows := []xlsxTable.TaskTableRow{}
	err := mgr.configDb.Find(&rows).Error
	if err != nil {
		return err
	}

	for _, row := range rows {
		mgr.taskCnf[row.Id] = row
	}

	return nil
}

func (mgr *ConfigDataManager) DailyTaskRateLimit() int32 { return 10 }

func (mgr *ConfigDataManager) RewardTaskListRateMin() int32 { return 50 }

func (mgr *ConfigDataManager) RewardTaskListRateMax() int32 { return 100 }

func (mgr *ConfigDataManager) AllTaskCnfs() map[int32]xlsxTable.TaskTableRow {
	return configMgr.taskCnf
}

func (mgr *ConfigDataManager) TaskCnfById(id int32) *xlsxTable.TaskTableRow {
	cnf, exist := mgr.taskCnf[id]
	if !exist {
		return nil
	}
	return &cnf
}
