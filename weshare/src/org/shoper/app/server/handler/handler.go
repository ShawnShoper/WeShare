package handler

import (
	"org/shoper/app/server/session"
)

type Handler interface {
	SessionOpened(session session.Session, message interface{}) error
	ExceptionCaught(session session.Session, err error)
	MessageReceived(session session.Session, message interface{}) error
	SessionClosed(session session.Session, message interface{}) error
}
