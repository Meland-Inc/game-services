package globalModule

import (
	"fmt"
	"game-message-core/grpc"
	"game-message-core/grpc/methodData"
	"game-message-core/proto"
	"game-message-core/protoTool"
	"time"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/global/contract"
	"github.com/Meland-Inc/game-services/src/global/grpcAPI/grpcNetTool"
	"github.com/Meland-Inc/game-services/src/global/module"
	"github.com/Meland-Inc/game-services/src/global/userAgent"
)

type ServiceEventBase struct {
	module.ModuleBase
	modelEvent *module.ModuleEvent
	subModule  contract.IServiceEvent

	clientEvents      map[proto.EnvelopeType]func(*userAgent.UserAgentData, *methodData.PullClientMessageInput, *proto.Envelope)
	gameSerDaprEvents map[string]contract.ServiceEventFunc
	gameSerDaprCalls  map[string]contract.ServiceEventFunc
	web3DaprEvents    map[string]contract.ServiceEventFunc
	web3DaprCalls     map[string]contract.ServiceEventFunc
}

func GetServiceEventModel() (contract.IServiceEvent, error) {
	iModel, exist := module.GetModel(module.MODULE_NAME_SERVICE_EVENT)
	if !exist {
		return nil, fmt.Errorf("service event model not found")
	}
	iEventModel, ok := iModel.(contract.IServiceEvent)
	if !ok {
		return nil, fmt.Errorf("service event iModule to IServiceEvent failed")
	}
	return iEventModel, nil
}

func (p *ServiceEventBase) Init(subModel contract.IServiceEvent) {
	p.clientEvents = make(map[proto.EnvelopeType]func(*userAgent.UserAgentData, *methodData.PullClientMessageInput, *proto.Envelope))
	p.gameSerDaprEvents = make(map[string]contract.ServiceEventFunc)
	p.gameSerDaprCalls = make(map[string]contract.ServiceEventFunc)
	p.web3DaprEvents = make(map[string]contract.ServiceEventFunc)
	p.web3DaprCalls = make(map[string]contract.ServiceEventFunc)

	p.subModule = subModel
	p.modelEvent = module.NewModelEvent()
	p.ModuleBase.InitBaseModel(subModel, module.MODULE_NAME_SERVICE_EVENT)

	p.registerEvent()
}

func (p *ServiceEventBase) OnInit() error {
	p.ModuleBase.OnInit()
	return nil
}

func (p *ServiceEventBase) OnTick(utc time.Time) {
	p.ModuleBase.OnTick(utc)
	if env := p.ReadEvent(); env != nil {
		p.OnEvent(env, utc.UnixMilli())
	}
}

func (p *ServiceEventBase) GetGameServiceDaprCallTypes() []string {
	keys := []string{}
	for key, _ := range p.gameSerDaprCalls {
		keys = append(keys, key)
	}
	return keys
}

func (p *ServiceEventBase) GetGameServiceDaprEventTypes() []string {
	keys := []string{}
	for key, _ := range p.gameSerDaprEvents {
		keys = append(keys, key)
	}
	return keys
}

func (p *ServiceEventBase) GetWeb3DaprCallTypes() []string {
	keys := []string{}
	for key, _ := range p.web3DaprCalls {
		keys = append(keys, key)
	}
	return keys
}

func (p *ServiceEventBase) GetWeb3DaprEventTypes() []string {
	keys := []string{}
	for key, _ := range p.web3DaprEvents {
		keys = append(keys, key)
	}
	return keys
}

func (p *ServiceEventBase) registerEvent() {
	p.subModule.RegisterClientEvent()
	p.subModule.RegisterGameServiceDaprCall()
	p.subModule.RegisterGameServiceDaprEvent()
	p.subModule.RegisterWeb3DaprCall()
	p.subModule.RegisterWeb3DaprEvent()
}

func (p *ServiceEventBase) RegisterClientEvent()          {}
func (p *ServiceEventBase) RegisterGameServiceDaprCall()  {}
func (p *ServiceEventBase) RegisterGameServiceDaprEvent() {}
func (p *ServiceEventBase) RegisterWeb3DaprCall()         {}
func (p *ServiceEventBase) RegisterWeb3DaprEvent()        {}

func (p *ServiceEventBase) AddClientEvent(
	msgType proto.EnvelopeType,
	f func(*userAgent.UserAgentData, *methodData.PullClientMessageInput, *proto.Envelope),
) {
	if msgType <= proto.EnvelopeType_Unknown || f == nil {
		return
	}
	p.clientEvents[msgType] = f
}

func (p *ServiceEventBase) AddGameServiceDaprCall(name string, f contract.ServiceEventFunc) {
	if name == "" || f == nil {
		return
	}
	p.gameSerDaprCalls[name] = f
}

func (p *ServiceEventBase) AddGameServiceDaprEvent(name string, f contract.ServiceEventFunc) {
	if name == "" || f == nil {
		return
	}
	p.gameSerDaprEvents[name] = f
}

func (p *ServiceEventBase) AddWeb3DaprCall(name string, f contract.ServiceEventFunc) {
	if name == "" || f == nil {
		return
	}
	p.web3DaprCalls[name] = f
}

func (p *ServiceEventBase) AddWeb3DaprEvent(name string, f contract.ServiceEventFunc) {
	if name == "" || f == nil {
		return
	}
	p.web3DaprEvents[name] = f
}

func (p *ServiceEventBase) EventCall(env contract.IModuleEventReq) contract.IModuleEventResult {
	return p.modelEvent.EventCall(env)
}

func (p *ServiceEventBase) EventCallNoReturn(env contract.IModuleEventReq) {
	p.modelEvent.EventCallNoReturn(env)
}

func (p *ServiceEventBase) ReadEvent() contract.IModuleEventReq {
	return p.modelEvent.ReadEvent()
}

func (p *ServiceEventBase) OnEvent(env contract.IModuleEventReq, curMs int64) {
	defer func() {
		err := recover()
		if err != nil {
			serviceLog.StackError("ServiceEventBase.onEvent err: %v", err)
		}
	}()

	if env.GetEventType() == string(grpc.ProtoMessageActionPullClientMessage) {
		p.callClientMsg(env, curMs)
	} else {
		p.callServiceEvent(env, curMs)
	}
}

func (p *ServiceEventBase) callClientMsg(env contract.IModuleEventReq, curMs int64) {
	bs, ok := env.GetMsg().([]byte)
	// serviceLog.Info("client msg: %s, [%v]", bs, ok)
	if !ok {
		serviceLog.Error("client msg to string failed: %v", bs)
		return
	}

	// serviceLog.Info("service event received clientPbMsg data: %v", string(bs))

	input := &methodData.PullClientMessageInput{}
	err := grpcNetTool.UnmarshalGrpcData(bs, input)
	if err != nil {
		serviceLog.Error("client msg input Unmarshal error: %v", err)
		return
	}

	agent := userAgent.GetOrStoreUserAgent(input)
	msg, err := protoTool.UnMarshalToEnvelope(input.MsgBody)
	if err != nil {
		serviceLog.Error("Unmarshal Envelope fail err: %+v", err)
		return
	}

	pbMsgType := proto.EnvelopeType(input.MsgId)
	if handler, exist := p.clientEvents[pbMsgType]; exist {
		handler(agent, input, msg)
	}
}

func (p *ServiceEventBase) callServiceEvent(env contract.IModuleEventReq, curMs int64) {
	eventName := env.GetEventType()

	if handler, exist := p.gameSerDaprCalls[eventName]; exist {
		handler(env, curMs)
		return
	}
	if handler, exist := p.gameSerDaprEvents[eventName]; exist {
		handler(env, curMs)
		return
	}
	if handler, exist := p.web3DaprCalls[eventName]; exist {
		handler(env, curMs)
		return
	}
	if handler, exist := p.web3DaprEvents[eventName]; exist {
		handler(env, curMs)
		return
	}
}
