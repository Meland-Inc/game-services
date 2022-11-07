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

	if err := taskModel.KillMonsterHandler(
		input.UserId,
		proto.TaskListType_TaskListTypeUnknown,
		&proto.TaskOptionKillMonster{MonCid: input.MonsterCid, Num: 1},
	); err != nil {
		serviceLog.Error("task killMon handler err:%+v", err)
	}

	if len(input.DropList) > 0 {
		pickItems := []*proto.TaskOptionItem{}
		for _, drop := range input.DropList {
			pickItems = append(pickItems, &proto.TaskOptionItem{ItemCid: drop.Cid, Num: drop.Num})
		}
		if err := taskModel.GetItemHandler(
			input.UserId, proto.TaskListType_TaskListTypeUnknown, pickItems,
		); err != nil {
			serviceLog.Error("task get item handler err:%+v", err)
		}
	}
}
