package taskModel

import (
	"game-message-core/proto"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	dbData "github.com/Meland-Inc/game-services/src/global/gameDB/data"
	"github.com/Meland-Inc/game-services/src/global/serviceCnf"
	"github.com/Meland-Inc/game-services/src/global/userAgent"
)

func (p *TaskModel) broadCastUpdateTaskListInfo(userId int64, tlType proto.TaskListType, tl *dbData.TaskList) {
	var tlPbData *proto.TaskList
	if tl != nil {
		tlPbData = tl.ToPbData()
	}
	msg := &proto.Envelope{
		Type: proto.EnvelopeType_BroadCastUpdateTaskList,
		Payload: &proto.Envelope_BroadCastUpdateTaskListResponse{
			BroadCastUpdateTaskListResponse: &proto.BroadCastUpdateTaskListResponse{
				Kind:         tlType,
				TaskListInfo: tlPbData,
			},
		},
	}

	agentModel := userAgent.GetUserAgentModel()
	agent, exist := agentModel.GetUserAgent(userId)
	if !exist {
		serviceLog.Warning("user [%d] agent data not found", userId)
		return
	}
	agent.SendToPlayer(serviceCnf.GetInstance().ServerName, msg)
}

func (p *TaskModel) broadCastReceiveRewardInfo(
	userId int64, tl *dbData.TaskList,
	rewardExp int32, rewardItemCid, rewardItemNum, rewardItemQuality int32,
	isTaskListReward bool,
) {
	if tl == nil {
		return
	}

	msg := &proto.Envelope{
		Type: proto.EnvelopeType_BroadCastTaskReward,
		Payload: &proto.Envelope_BroadCastTaskRewardResponse{
			BroadCastTaskRewardResponse: &proto.BroadCastTaskRewardResponse{
				IsTaskListReward: isTaskListReward,
				TaskListKind:     proto.TaskListType(tl.TaskListType),
				RewardExp:        rewardExp,
				RewardItem: []*proto.ItemBaseInfo{&proto.ItemBaseInfo{
					Cid:     rewardItemCid,
					Num:     rewardItemNum,
					Quality: rewardItemQuality,
				}},
			},
		},
	}

	agentModel := userAgent.GetUserAgentModel()
	agent, exist := agentModel.GetUserAgent(userId)
	if !exist {
		serviceLog.Warning("user [%d] agent data not found", userId)
		return
	}
	agent.SendToPlayer(serviceCnf.GetInstance().ServerName, msg)
}
