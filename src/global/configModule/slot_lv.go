package configModule

import (
	xlsxTable "game-message-core/xlsxTableData"
)

func (mgr *ConfigDataManager) initSlotLv() error {
	mgr.slotLvCnf = make(map[int32][]*xlsxTable.SlotLvTableRow)

	rows := []xlsxTable.SlotLvTableRow{}
	err := mgr.configDb.Find(&rows).Error
	if err != nil {
		return err
	}

	for _, row := range rows {
		_, exist := mgr.slotLvCnf[row.Position]
		if exist {
			mgr.slotLvCnf[row.Position] = append(mgr.slotLvCnf[row.Position], &row)
		} else {
			mgr.slotLvCnf[row.Position] = []*xlsxTable.SlotLvTableRow{&row}
		}
	}

	return nil
}

func (mgr *ConfigDataManager) GetSlotCnf(position, lv int32) *xlsxTable.SlotLvTableRow {
	if settings, exist := mgr.slotLvCnf[position]; exist {
		for _, setting := range settings {
			if setting.Lv == lv {
				return setting
			}
		}
	}
	return nil
}
