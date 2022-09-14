package service

import (
	"fmt"
	"time"

	"github.com/Meland-Inc/game-services/src/common/net/session"
	"github.com/Meland-Inc/game-services/src/common/net/tcp"
	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/global/serviceRegister"
	"github.com/Meland-Inc/game-services/src/services/agent/userChannel"
)

func (s *Service) onStart() error {
	if err := s.modelMgr.StartModel(); err != nil {
		return err
	}

	if err := s.registerService(); err != nil {
		return err
	}

	if err := s.initTcpServer(); err != nil {
		return err
	}

	return nil
}

func (s *Service) registerService() error {
	time.Sleep(time.Millisecond * 300) // 延时 300Ms 等待dapr init 完成
	err := serviceRegister.RegisterService(*s.serviceCnf, 0)
	serviceLog.Info("registerService ------ end ----------data: %+v, err: %v", *s.serviceCnf, err)
	return err
}

func (s *Service) initTcpServer() (err error) {
	s.tcpServer, err = tcp.NewTcpServer(
		fmt.Sprintf(":%d", s.serviceCnf.Port),
		uint32(s.serviceCnf.MaxOnline),
		180,
		s.OnSessionConnect,
	)
	return err
}

func (s *Service) OnSessionConnect(se *session.Session) {
	fmt.Printf("session [%s][%s] ---- Connect to agent service ", se.SessionId(), se.RemoteAddr())
	channel := userChannel.NewUserChannel(se)
	se.SetCallBack(channel.OnSessionReceivedData, channel.OnSessionClose)
	userChannel.GetInstance().AddUserChannelById(channel)
	channel.Run()
}
