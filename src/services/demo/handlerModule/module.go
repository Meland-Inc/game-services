package handlerModule

import (
	"github.com/Meland-Inc/game-services/src/global/globalModule"
)

type HandlerModule struct {
	globalModule.ServiceEventBase
}

func NewHandlerModule() *HandlerModule {
	p := &HandlerModule{}
	p.ServiceEventBase.Init(p)
	return p
}
