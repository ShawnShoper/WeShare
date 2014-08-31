package server

import (
	"net"
	"time"
)

//user login info
type User struct {
	id         string     //id
	loginTime  time.Time  //login time
	logoutTime time.Time  //logout time
	ip         net.IPAddr //login ip
	via        string     //login via
}
type ID string
