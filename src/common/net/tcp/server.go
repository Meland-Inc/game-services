package tcp

import (
	"log"
	"net"
	"time"

	"github.com/Meland-Inc/game-services/src/common/net/session"
	"github.com/Meland-Inc/game-services/src/common/serviceLog"
)

const (
	CONNECT_LIMIT = 50000
)

type Server struct {
	listener          net.Listener
	addr              string
	maxConNum         uint32
	sessionMgr        *session.SessionManager
	onConnectCallback func(s *session.Session)
}

func NewTcpServer(
	addr string, maxConnNum uint32, timeoutSec int64,
	connectCallback func(s *session.Session),
) (*Server, error) {
	s := &Server{
		addr:              addr,
		maxConNum:         maxConnNum,
		sessionMgr:        session.NewSessionMgr(timeoutSec),
		onConnectCallback: connectCallback,
	}

	if s.maxConNum < 1 || s.maxConNum > CONNECT_LIMIT {
		s.maxConNum = CONNECT_LIMIT
	}

	// s.setLimit()

	err := s.ListenAndServe()
	return s, err
}

func (s *Server) SessionMgr() *session.SessionManager {
	return s.sessionMgr
}
func (s *Server) setLimit() {
	// var rLimit syscall.Rlimit
	// if err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit); err != nil {
	// 	panic(err)
	// }
	// rLimit.Cur = rLimit.Max
	// if err := syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit); err != nil {
	// 	panic(err)
	// }

	// log.Printf("set cur limit: %d", rLimit.Cur)
}

func (s *Server) Stop() error {
	return s.listener.Close()
}

func (s *Server) ListenAndServe() (err error) {
	log.Println("socket listen ", s.addr)
	s.listener, err = net.Listen("tcp", s.addr)
	if err != nil {
		log.Println("listen ", err.Error())
		return err
	}

	go func() { s.listen() }()
	return nil
}

func (s *Server) listen() {
	defer s.listener.Close()
	var tempDelay time.Duration
	for {
		connect, err := s.listener.Accept()
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				tempDelay = 5 * time.Millisecond
				log.Println("accept error:", err, "retrying in :", tempDelay)
				time.Sleep(tempDelay)
				continue
			} else {
				return
			}
		}
		tempDelay = 0

		s.onConnect(connect)
	}
}

func (s *Server) onConnect(connect net.Conn) {
	count := s.sessionMgr.Count()
	serviceLog.Debug("%s: connected to agent, current count: %d", connect.RemoteAddr(), count)

	if count >= s.maxConNum {
		connect.Close()
		serviceLog.Warning("too many connections, connCount(%d) >= maxConnNum(%d)", count, s.maxConNum)
		return
	}

	session := session.NewSession(connect)
	if session != nil {
		serviceLog.Info("tcp socket connected remoteAddr[%v], localAddress[%v], socketId[%s]",
			connect.RemoteAddr(), connect.LocalAddr(), session.SessionId())
		s.sessionMgr.AddSession(session)
		s.onConnectCallback(session)
		go func() {
			session.Run()
			session.Stop()
			s.sessionMgr.RemoveSession(session)
		}()
	} else {
		connect.Close()
	}
}
