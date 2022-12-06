package configData

import (
	xlsxTable "game-message-core/xlsxTableData"
)

func (mgr *ConfigDataManager) initSceneArea() error {
	mgr.sceneAreaCnf = make(map[int32]xlsxTable.SceneAreaRow)

	rows := []xlsxTable.SceneAreaRow{}
	err := mgr.configDb.Find(&rows).Error
	if err != nil {
		return err
	}

	for _, row := range rows {
		mgr.sceneAreaCnf[row.Id] = row
	}
	return nil
}

func (mgr *ConfigDataManager) GetSceneArea(id int32) (xlsxTable.SceneAreaRow, bool) {
	row, exist := mgr.sceneAreaCnf[id]
	return row, exist
}

func (mgr *ConfigDataManager) AllSceneArea() []xlsxTable.SceneAreaRow {
	rows := make([]xlsxTable.SceneAreaRow, 0, 0)
	for _, row := range mgr.sceneAreaCnf {
		rows = append(rows, row)
	}
	return rows
}
