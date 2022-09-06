package configModule

import (
	xlsxTable "game-message-core/xlsxTableData"
)

func (mgr *ConfigDataManager) initRoleLv() error {
	mgr.roleLvCnf = make(map[int32]*xlsxTable.RoleLvTableRow)

	rows := []xlsxTable.RoleLvTableRow{}
	err := mgr.configDb.Find(&rows).Error
	if err != nil {
		return err
	}

	for _, row := range rows {
		mgr.roleLvCnf[row.Lv] = &row
	}

	return nil
}

func (mgr *ConfigDataManager) LevelCnfById(lv int32) *xlsxTable.RoleLvTableRow {
	cnf, exist := mgr.roleLvCnf[lv]
	if !exist {
		return nil
	}
	return cnf
}
