package session

import (
	"org/shoper/app/server/bean"
)

type Session struct {
	//存在用户多终端登录的情况
	clients   map[bean.ID]bean.User //key:用户id,value User
	timeout   int                   //超时时间
	keepAlive bool                  //是否长连接
	currConn  bean.User             //当前连接
}
