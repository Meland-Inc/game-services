package configData

import (
	"game-message-core/proto"
	xlsxTable "game-message-core/xlsxTableData"
)

func (mgr *ConfigDataManager) initTaskList() error {
	mgr.taskListCnf = make(map[int32]xlsxTable.TaskListTableRow)

	rows := []xlsxTable.TaskListTableRow{}
	err := mgr.configDb.Find(&rows).Error
	if err != nil {
		return err
	}

	for _, row := range rows {
		mgr.taskListCnf[row.Id] = row
	}
	return nil
}

func (mgr *ConfigDataManager) AllTaskListCnfs() map[int32]xlsxTable.TaskListTableRow {
	return configMgr.taskListCnf
}

func (mgr *ConfigDataManager) TaskListCnfById(id int32) *xlsxTable.TaskListTableRow {
	cnf, exist := mgr.taskListCnf[id]
	if !exist {
		return nil
	}
	return &cnf
}

func (mgr *ConfigDataManager) TaskListCnfByLevel(taskListType, lv int32) *xlsxTable.TaskListTableRow {
	for _, cnf := range configMgr.taskListCnf {
		if cnf.System != taskListType {
			continue
		}
		if cnf.System == int32(proto.TaskListType_TaskListTypeGuide) || cnf.Level == lv {
			return &cnf
		}
	}
	return nil
}
