package configData

import (
	xlsxTable "game-message-core/xlsxTableData"
)

func (mgr *ConfigDataManager) initDrop() error {
	mgr.dropCnf = make(map[int32]xlsxTable.DropTableRow)

	rows := []xlsxTable.DropTableRow{}
	err := mgr.configDb.Find(&rows).Error
	if err != nil {
		return err
	}

	for _, row := range rows {
		mgr.dropCnf[row.DropId] = row
	}

	return nil
}

func (mgr *ConfigDataManager) DropCnfById(id int32) *xlsxTable.DropTableRow {
	cnf, exist := mgr.dropCnf[id]
	if !exist {
		return nil
	}
	return &cnf
}
