package tcp

import (
	"bytes"
	"encoding/json"
	//"fmt"
	l4g "github.com/alecthomas/log4go"
	"net"
	"org/shoper/app/server"
	"os"
	"strconv"
)

type Handler interface {
	SessionOpened()
	ExceptionCaught()
	MessageReceived()
	SessionClosed()
}

var (
	handler Handler //
)

const (
	FILE_PATH           string = "config/serverconfig.json"
	DEFAULT_SERVER_IP   string = "127.0.0.1"
	DEFAULT_SERVER_PORT int    = 9999
)

type TCPServer struct {
}
type Session struct {
	//存在用户多终端登录的情况
	clients   map[server.ID]server.User //key:用户id,value User
	Timeout   int                       //超时时间
	KeepAlive bool                      //是否长连接
}
type ServerStruct struct {
	Server_ip   string
	Server_port int
}

func getConnectionInfo() (s ServerStruct) {
	s = ServerStruct{}
	l4g.Info("get server config info")
	file, err := os.Open(FILE_PATH)
	defer file.Close()
	// if the file is not exists then create new file to save config info
	if err != nil {
		l4g.Warn("open file fail, and create new file to write info msg.")
		newFile, oerr := os.Create(FILE_PATH)
		defer newFile.Close()
		if oerr == nil {
			s.Server_ip = DEFAULT_SERVER_IP
			s.Server_port = DEFAULT_SERVER_PORT
			//take ServerStruct convert to json
			v, e := json.Marshal(s)
			if e != nil {
				l4g.Error(e.Error())
			} else {
				newFile.Write(v)
			}
		}
	} else {
		l4g.Info("read config info")
		b := make([]byte, 1024)
		info := bytes.Buffer{}
		for {
			n, _ := file.Read(b)
			if n == 0 {
				break
			}
			info.Write(b[:n])
		}
		json.Unmarshal(info.Bytes(), &s)
	}
	return
}

func (t *TCPServer) StartTCPServer() {
	info := getConnectionInfo()
	var address string
	if info.Server_ip == "" || info.Server_port == 0 {
		l4g.Info("In server config file ,the server_id or server_port is empty, use default info")
		address = DEFAULT_SERVER_IP + ":" + strconv.Itoa(DEFAULT_SERVER_PORT)
	} else {
		address = info.Server_ip + ":" + strconv.Itoa(info.Server_port)
	}
	for i := 0; i < 3; i++ {
		listen, err := net.Listen("tcp", address)
		defer listen.Close()
		if err != nil {
			//start server error
			l4g.Warn("start server error" + err.Error())
			l4g.Info("try to restart server " + strconv.Itoa(i+1) + " of 3")
			continue
		} else {
			for {
				l4g.Info("Listening connections...")
				_, ce := listen.Accept()
				if ce != nil {
					l4g.Error("get connection error:" + ce.Error())
				}
			}
		}
	}
}
