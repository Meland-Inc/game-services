package configData

import (
	"game-message-core/proto"
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

func (mgr *ConfigDataManager) GetSceneArea(id int32) (xlsxTable.SceneAreaRow, proto.SceneServiceSubType, bool) {
	row, exist := mgr.sceneAreaCnf[id]
	if !exist {
		return row, proto.SceneServiceSubType_UnknownSubType, false
	}
	return row, ToServiceSubType(row.SceneType), exist
}

func (mgr *ConfigDataManager) AllSceneArea() []xlsxTable.SceneAreaRow {
	rows := make([]xlsxTable.SceneAreaRow, 0, 0)
	for _, row := range mgr.sceneAreaCnf {
		rows = append(rows, row)
	}
	return rows
}

func ToServiceSubType(subTypeStr string) proto.SceneServiceSubType {
	t := proto.SceneServiceSubType_UnknownSubType
	switch subTypeStr {
	case "world":
		t = proto.SceneServiceSubType_World
	case "home":
		t = proto.SceneServiceSubType_Home
	case "dungeon":
		t = proto.SceneServiceSubType_Dungeon
	}
	return t
}
