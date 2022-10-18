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
}

func NewSessionMgr(maxConnNum uint32, timeoutSec int64) *SessionManager {
	mgr := &SessionManager{
		maxConn:    maxConnNum,
		timeoutSec: timeoutSec,
		count:      0,
		sessions:   make(map[string]*Session),
	}
	go func() {
		mgr.checkTimeout()
	}()
	return mgr
}

func (mgr *SessionManager) Count() uint32 { return mgr.count }

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
		time.Sleep(5 * time.Second)

		timeoutSessions := make([]*Session, 0, 5)
		nowSec := time.Now().UTC().Unix()
		mgr.RangeSessions(func(s *Session) bool {
			if s != nil {
				// 间隔大于?秒客户端超时
				if (nowSec - s.GetActiveTime()) > mgr.timeoutSec {
					serviceLog.Info("session remote addr[%v][%s] timeout and close", s.RemoteAddr(),s.SessionId())
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
