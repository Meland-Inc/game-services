package serviceMsgHandle

import (
	"game-message-core/grpc/pubsubEventData"
	"game-message-core/proto"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/services/task/taskModel"
)

func UserLevelUpgradeHandler(iMsg interface{}) {
	input, ok := iMsg.(*pubsubEventData.UserLevelUpgradeEvent)
	if !ok {
		serviceLog.Error("iMsg to userLevelUpgrade failed")
		return
	}

	taskModel, err := taskModel.GetTaskModel()
	if err != nil {
		serviceLog.Error("userLevelUpgrade taskModel not found")
		return
	}

	if err := taskModel.UserLevelHandler(
		input.UserId,
		proto.TaskListType_TaskListTypeUnknown,
		input.Level,
	); err != nil {
		serviceLog.Error("task user level handler err:%+v", err)
	}
}
