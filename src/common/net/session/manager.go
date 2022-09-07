package session

import (
	"sync"
	"time"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
)

type SessionManager struct {
	count      uint32
	timeoutSec int64
	sessions   sync.Map
}

func NewSessionMgr(timeoutSec int64) *SessionManager {
	mgr := &SessionManager{timeoutSec: timeoutSec}
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
	mgr.sessions.Store(s.SessionId(), s)
	mgr.count++
}

func (mgr *SessionManager) RemoveSession(s *Session) {
	if s == nil {
		return
	}
	mgr.sessions.Delete(s.SessionId())
	mgr.count--
}

func (mgr *SessionManager) SessionById(id string) (s *Session) {
	mgr.sessions.Range(func(key interface{}, value interface{}) bool {
		session, ok := value.(*Session)
		if ok && session.SessionId() == id {
			s = session
			return false
		}
		return true
	})
	return s
}

func (mgr *SessionManager) SessionByOwner(ownerId int32) (s *Session) {
	mgr.sessions.Range(func(key interface{}, value interface{}) bool {
		session, ok := value.(*Session)
		if ok && session.GetOwner() == ownerId {
			s = session
			return false
		}
		return true
	})
	return s
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

		nowSec := time.Now().UTC().Unix()
		mgr.sessions.Range(func(key interface{}, value interface{}) bool {
			session, _ := value.(*Session)
			if session != nil {
				// 间隔大于?秒客户端超时
				if (nowSec - session.GetActiveTime()) > mgr.timeoutSec {
					serviceLog.Info("session remote addr[%v] timeout and close", session.RemoteAddr())
					session.Stop()
				}
			}
			return true
		})
	}
}
