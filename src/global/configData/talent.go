package configData

import (
	xlsxTable "game-message-core/xlsxTableData"
)

func (mgr *ConfigDataManager) initTalentTree() error {
	mgr.talentTreeCnf = make(map[int32]xlsxTable.TalentTreeRow)

	rows := []xlsxTable.TalentTreeRow{}
	err := mgr.configDb.Find(&rows).Error
	if err != nil {
		return err
	}

	for _, row := range rows {
		mgr.talentTreeCnf[row.NodeId] = row
	}

	return nil
}

func (mgr *ConfigDataManager) TalentNodeById(id int32) *xlsxTable.TalentTreeRow {
	cnf, exist := mgr.talentTreeCnf[id]
	if !exist {
		return nil
	}
	return &cnf
}
