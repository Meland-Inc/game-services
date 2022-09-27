package serviceMsgHandle

import (
	"game-message-core/grpc/pubsubEventData"
	"game-message-core/proto"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/services/task/taskModel"
)

func KillMonsterHandler(iMsg interface{}) {
	input, ok := iMsg.(*pubsubEventData.KillMonsterEventData)
	if !ok {
		serviceLog.Error("iMsg to KillMonsterEvent failed")
		return
	}

	taskModel, err := taskModel.GetTaskModel()
	if err != nil {
		serviceLog.Error("KillMonsterEvent taskModel not found")
		return
	}

	taskModel.UpGradeTaskProgress(
		input.UserId, proto.TaskListType_TaskListTypeDaily,
		nil, nil, nil, input.MonsterCid, 0, 0, 0,
	)
	taskModel.UpGradeTaskProgress(
		input.UserId, proto.TaskListType_TaskListTypeRewarded,
		nil, nil, nil, input.MonsterCid, 0, 0, 0,
	)
}
