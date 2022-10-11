package configData

import (
	"game-message-core/proto"
	xlsxTable "game-message-core/xlsxTableData"
)

func (mgr *ConfigDataManager) initChatCnf() error {
	mgr.chatCnf = make(map[proto.ChatChannelType]xlsxTable.ChatTableRow)

	rows := []xlsxTable.ChatTableRow{}
	err := mgr.configDb.Find(&rows).Error
	if err != nil {
		return err
	}

	for _, row := range rows {
		chatType := proto.ChatChannelType(row.ChatType)
		mgr.chatCnf[chatType] = row
	}

	return nil
}

func (mgr *ConfigDataManager) ChatCnfByType(chatType proto.ChatChannelType) *xlsxTable.ChatTableRow {
	cnf, exist := mgr.chatCnf[chatType]
	if !exist {
		return nil
	}
	return &cnf
}
