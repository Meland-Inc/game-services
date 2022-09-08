package configData

import (
	xlsxTable "game-message-core/xlsxTableData"
)

func (mgr *ConfigDataManager) initItem() error {
	mgr.itemCnf = make(map[int32]*xlsxTable.ItemTable)

	rows := []xlsxTable.ItemTable{}
	err := mgr.configDb.Find(&rows).Error
	if err != nil {
		return err
	}

	for _, row := range rows {
		mgr.itemCnf[row.ItemCid] = &row
	}

	return nil
}

func (mgr *ConfigDataManager) ItemCnfById(cid int32) *xlsxTable.ItemTable {
	cnf, exist := mgr.itemCnf[cid]
	if !exist {
		return nil
	}
	return cnf
}
