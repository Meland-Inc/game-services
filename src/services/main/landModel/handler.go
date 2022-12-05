package land_model

import (
	"game-message-core/grpc"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/global/component"
	message "github.com/Meland-Inc/game-services/src/global/web3Message"
)

func (p *LandModel) OnEvent(env *component.ModelEventReq, curMs int64) {
	defer func() {
		err := recover()
		if err != nil {
			serviceLog.StackError("ControllerModel.onEvent err: %v", err)
		}
	}()

	switch env.EventType {
	case string(message.SubscriptionEventMultiLandDataUpdateEvent):
		p.Web3MultiLandDataUpdateEvent(env, curMs)
	case string(message.SubscriptionEventMultiRecyclingEvent):
		p.Web3MultiRecyclingEvent(env, curMs)
	case string(message.SubscriptionEventMultiBuildUpdateEvent):
		p.Web3MultiBuildUpdateEvent(env, curMs)

	case string(grpc.ProtoMessageActionPullClientMessage):
		p.clientMsgHandler(env, curMs)
	case string(grpc.MainServiceActionGetAllBuild):
		p.GRPCGetAllBuildHandler(env, curMs)

	}

}
