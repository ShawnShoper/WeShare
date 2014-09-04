package bean

import (
	"net"
	"time"
)

//user login info
type User struct {
	id         ID         //id
	loginTime  time.Time  //login time
	logoutTime time.Time  //logout time
	ip         net.IPAddr //login ip
	via        string     //login via
	conn       net.Conn   //connection
}
type ID string
