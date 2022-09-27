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

	taskModel.UpGradeTaskProgress(
		input.UserId, proto.TaskListType_TaskListTypeDaily,
		nil, nil, nil, 0, input.Cid, input.Num, 0,
	)
	taskModel.UpGradeTaskProgress(
		input.UserId, proto.TaskListType_TaskListTypeRewarded,
		nil, nil, nil, 0, input.Cid, input.Num, 0,
	)
}
