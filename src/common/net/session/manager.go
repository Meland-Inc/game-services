package session

import (
	"sync"
	"time"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
)

type SessionManager struct {
	sync.RWMutex
	maxConn    uint32
	sessions   map[string]*Session
	count      uint32
	timeoutSec int64
	closed     bool
}

func NewSessionMgr(maxConnNum uint32, timeoutSec int64) *SessionManager {
	mgr := &SessionManager{
		maxConn:    maxConnNum,
		timeoutSec: timeoutSec,
		count:      0,
		sessions:   make(map[string]*Session),
		closed:     false,
	}
	go func() {
		mgr.checkTimeout()
	}()
	return mgr
}

func (mgr *SessionManager) Count() uint32 { return mgr.count }

func (mgr *SessionManager) Stop() {
	mgr.RLock()
	defer mgr.RUnlock()
	mgr.closed = true
	for _, s := range mgr.sessions {
		s.Stop()
	}
	mgr.count = 0
}

func (mgr *SessionManager) AddSession(s *Session) {
	if s == nil {
		return
	}

	mgr.RLock()
	defer mgr.RUnlock()

	_, exist := mgr.sessions[s.SessionId()]
	if !exist {
		mgr.count++
		if mgr.count > mgr.maxConn {
			serviceLog.StackError("cur connect > max connect num")
		}
	}
	mgr.sessions[s.SessionId()] = s
}

func (mgr *SessionManager) RemoveSession(s *Session) {
	if s == nil {
		return
	}
	mgr.RLock()
	defer mgr.RUnlock()

	_, exist := mgr.sessions[s.SessionId()]
	if !exist {
		return
	}

	delete(mgr.sessions, s.SessionId())
	mgr.count--
}

func (mgr *SessionManager) SessionById(id string) (s *Session) {
	mgr.RLock()
	defer mgr.RUnlock()
	s, _ = mgr.sessions[id]
	return s
}

func (mgr *SessionManager) SessionByOwner(ownerId int32) (s *Session) {
	mgr.RLock()
	defer mgr.RUnlock()
	for _, s := range mgr.sessions {
		if s.GetOwner() == ownerId {
			return s
		}
	}
	return nil
}

func (mgr *SessionManager) RangeSessions(f func(s *Session) bool) {
	mgr.RLock()
	defer mgr.RUnlock()
	for _, s := range mgr.sessions {
		if !f(s) {
			return
		}
	}
}

func (mgr *SessionManager) checkTimeout() {
	defer func() {
		err := recover()
		if err != nil {
			go mgr.checkTimeout()
		}
	}()

	for {
		time.Sleep(3 * time.Second)
		if mgr.closed {
			return
		}

		timeoutSessions := make([]*Session, 0, 5)
		nowSec := time.Now().UTC().Unix()
		mgr.RangeSessions(func(s *Session) bool {
			if s != nil {
				// ???????????????????????????????
				if (nowSec - s.GetActiveTime()) > mgr.timeoutSec {
					serviceLog.Info("session remote addr[%v][%s] timeout and close", s.RemoteAddr(), s.SessionId())
					timeoutSessions = append(timeoutSessions, s)
				}
			}
			return true
		})

		for _, s := range timeoutSessions {
			s.Stop()
			mgr.RemoveSession(s)
		}

	}
}
