package serviceMsgHandle

import (
	"game-message-core/grpc/pubsubEventData"
	"game-message-core/proto"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/services/task/taskModel"
)

func PlayerUseItemHandler(iMsg interface{}) {
	input, ok := iMsg.(*pubsubEventData.UserUseNFTEvent)
	if !ok {
		serviceLog.Error("iMsg to UserUseNFTEvent failed")
		return
	}

	taskModel, err := taskModel.GetTaskModel()
	if err != nil {
		serviceLog.Error("UserUseNFTEvent taskModel not found")
		return
	}

	if err := taskModel.UseItemHandler(
		input.UserId,
		proto.TaskListType_TaskListTypeUnknown,
		&proto.TaskOptionItem{ItemCid: input.Cid, Num: input.Num},
	); err != nil {
		serviceLog.Error("task use item handler err:%+v", err)
	}
}
