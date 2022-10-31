package service

import (
	"os"

	"github.com/Meland-Inc/game-services/src/common/net/tcp"
	"github.com/Meland-Inc/game-services/src/global/component"
	"github.com/Meland-Inc/game-services/src/global/serviceCnf"
)

type Service struct {
	serviceCnf *serviceCnf.ServiceConfig
	modelMgr   *component.ModelManager
	tcpServer  *tcp.Server
	osSignal   chan os.Signal
	stopChan   chan chan struct{}
	closed     bool
}

func NewAgentService() *Service {
	s := &Service{
		osSignal: make(chan os.Signal, 1),
		stopChan: make(chan chan struct{}, 1),
		closed:   false,
	}
	s.modelMgr = component.InitModelManager(s)
	return s
}

func (s *Service) OnInit() error {
	return s.init()
}

func (s *Service) OnStart() error {
	return s.onStart()
}

func (s *Service) Run() {
	s.run()
}

func (s *Service) OnExit() {
	s.onStop()
	s.onExit()
}
