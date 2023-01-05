package home_model

import (
	"game-message-core/grpc"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/global/contract"
)

func (p *HomeModel) OnEvent(env contract.IModuleEventReq, curMs int64) {
	defer func() {
		err := recover()
		if err != nil {
			serviceLog.StackError("HomeModel.onEvent err: %v", err)
		}
	}()

	switch env.GetEventType() {
	case string(grpc.ProtoMessageActionPullClientMessage):
		p.clientMsgHandler(env, curMs)

	case string(grpc.MainServiceActionGetHomeData):
		p.GRPCGetHomeDataHandler(env, curMs)
	case string(grpc.SubscriptionEventSaveHomeData):
		p.GRPCSaveHomeDataEvent(env, curMs)
	case string(grpc.SubscriptionEventGranaryStockpile):
		p.GRPCGranaryStockpileEvent(env, curMs)

	}
}
