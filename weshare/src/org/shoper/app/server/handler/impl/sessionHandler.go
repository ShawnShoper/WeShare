package impl

import (
	l4g "github.com/alecthomas/log4go"
	"org/shoper/app/server/session"
)

type SessionHanler struct {
}

func (s *SessionHanler) SessionOpened(session session.Session, message interface{}) error {
	l4g.Info("Session ip %s come in.", session.IP)
	return nil
}
func (s *SessionHanler) ExceptionCaught(session session.Session, err error) {
	l4g.Error("Session has caught a exceptions : %s", err.Error())
}
func (s *SessionHanler) MessageReceived(session session.Session, message interface{}) error {
	l4g.Error("Session received a message")
	return nil
}
func (s *SessionHanler) SessionClosed(session session.Session, message interface{}) error {
	l4g.Error("Session closed")
	return nil
}
