package handler

import (
	"org/shoper/app/server/session"
)

type Handler interface {
	SessionOpened(session session.Session) error
	ExceptionCaught(session session.Session) error
	MessageReceived(session session.Session) error
	SessionClosed(session session.Session) error
}
