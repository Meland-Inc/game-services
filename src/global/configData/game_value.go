package configData

import (
	"fmt"
	xlsxTable "game-message-core/xlsxTableData"
)

func (mgr *ConfigDataManager) initGameValue() error {
	mgr.gameValueCnf = make(map[int32]xlsxTable.GameValueTable)

	rows := []xlsxTable.GameValueTable{}
	err := mgr.configDb.Find(&rows).Error
	if err != nil {
		return err
	}

	for _, row := range rows {
		mgr.gameValueCnf[row.Id] = row
	}

	return nil
}

func RoleCurrentExpLimit() int32 { return 2100000000 }

func GameValueById(id int32) (xlsxTable.GameValueTable, error) {
	value := xlsxTable.GameValueTable{}
	if id == 0 {
		return value, fmt.Errorf("invalid game value id[%d]", id)
	}

	value, exist := ConfigMgr().gameValueCnf[id]
	if !exist {
		return value, fmt.Errorf("game value id[%d] not found", id)
	}
	return value, nil
}
