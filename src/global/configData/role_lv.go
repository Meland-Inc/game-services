package configData

import (
	xlsxTable "game-message-core/xlsxTableData"
)

func (mgr *ConfigDataManager) initRoleLv() error {
	mgr.roleLvMax = 1
	mgr.roleLvCnf = make(map[int32]xlsxTable.RoleLvTableRow)

	rows := []xlsxTable.RoleLvTableRow{}
	err := mgr.configDb.Find(&rows).Error
	if err != nil {
		return err
	}

	for _, row := range rows {
		mgr.roleLvCnf[row.Lv] = row
		if mgr.roleLvMax < row.Lv {
			mgr.roleLvMax = row.Lv
		}
	}

	return nil
}

func (mgr *ConfigDataManager) RoleLevelCnf(lv int32) *xlsxTable.RoleLvTableRow {
	cnf, exist := mgr.roleLvCnf[lv]
	if !exist {
		return nil
	}
	return &cnf
}

func (mgr *ConfigDataManager) RoleMaxLevel() int32 {
	return mgr.roleLvMax
}
