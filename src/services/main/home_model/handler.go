package home_model

import (
	"game-message-core/grpc"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/global/component"
)

func (p *HomeModel) OnEvent(env *component.ModelEventReq, curMs int64) {
	defer func() {
		err := recover()
		if err != nil {
			serviceLog.StackError("HomeModel.onEvent err: %v", err)
		}
	}()

	switch env.EventType {
	case string(grpc.MainServiceActionGetHomeData):
		p.GRPCGetHomeDataHandler(env, curMs)
	case string(grpc.SubscriptionEventSaveHomeData):
		p.GRPCSaveHomeDataEvent(env, curMs)
	}
}
