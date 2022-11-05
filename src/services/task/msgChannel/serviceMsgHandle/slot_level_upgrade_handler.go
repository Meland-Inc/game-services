package serviceMsgHandle

import (
	"game-message-core/grpc/pubsubEventData"
	"game-message-core/proto"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/services/task/taskModel"
)

func SlotLevelUpgradeHandler(iMsg interface{}) {
	input, ok := iMsg.(*pubsubEventData.SlotLevelUpgradeEvent)
	if !ok {
		serviceLog.Error("iMsg to SlotLevelUpgradeEvent failed")
		return
	}

	taskModel, err := taskModel.GetTaskModel()
	if err != nil {
		serviceLog.Error("SlotLevelUpgradeEvent taskModel not found")
		return
	}

	if err := taskModel.TargetSlotLevelHandler(
		input.UserId,
		proto.TaskListType_TaskListTypeUnknown,
		&proto.TaskOptionTargetSlotLevel{SlotPos: input.SlotPos, Level: input.Level},
	); err != nil {
		serviceLog.Error("task slot level upgrade err:%+v", err)
	}

	if err := taskModel.SlotLevelCountHandler(
		input.UserId, proto.TaskListType_TaskListTypeUnknown,
	); err != nil {
		serviceLog.Error("task slot level count err:%+v", err)
	}

}
