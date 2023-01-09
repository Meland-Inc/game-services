package contract

type ServiceEventFunc func(env IModuleEventReq, curMs int64)

type IServiceEvent interface {
	IModuleInterface
	IModuleEvent
	GetGameServiceDaprCallTypes() []string
	GetGameServiceDaprEventTypes() []string
	GetWeb3DaprCallTypes() []string
	GetWeb3DaprEventTypes() []string
	RegisterClientEvent()
	RegisterGameServiceDaprCall()
	RegisterGameServiceDaprEvent()
	RegisterWeb3DaprCall()
	RegisterWeb3DaprEvent()
}
