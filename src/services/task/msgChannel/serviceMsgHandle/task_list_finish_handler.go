package serviceMsgHandle

import (
	"game-message-core/grpc/pubsubEventData"
	"game-message-core/proto"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/services/task/taskModel"
)

func TaskListFinishHandler(iMsg interface{}) {
	input, ok := iMsg.(*pubsubEventData.TaskListFinishEvent)
	if !ok {
		serviceLog.Error("iMsg to TaskListFinishEvent failed")
		return
	}

	taskModel, err := taskModel.GetTaskModel()
	if err != nil {
		serviceLog.Error("TaskListFinishEvent taskModel not found")
		return
	}

	if err := taskModel.TaskListTypeCountHandler(
		input.UserId,
		proto.TaskListType_TaskListTypeUnknown,
		input.TaskListType,
	); err != nil {
		serviceLog.Error("task TaskListFinishEvent err:%+v", err)
	}

	getItems := []*proto.TaskOptionItem{}
	for _, item := range input.RewardItem {
		getItems = append(getItems, &proto.TaskOptionItem{
			ItemCid: item.Cid,
			Num:     item.Num,
		})
	}
	if err := taskModel.GetItemHandler(
		input.UserId, proto.TaskListType_TaskListTypeUnknown, getItems,
	); err != nil {
		serviceLog.Error("task TaskFinishEvent get item err:%+v", err)
	}
}
