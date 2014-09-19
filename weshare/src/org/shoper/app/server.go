package main

import (
	l4g "github.com/alecthomas/log4go"
	handler "org/shoper/app/server/handler/impl"
	tcp "org/shoper/app/server/tcp"
)

func main() {
	l4g.AddFilter("stdout", l4g.DEBUG, l4g.NewConsoleLogWriter())             //输出到控制台,级别为DEBUG
	l4g.AddFilter("file", l4g.DEBUG, l4g.NewFileLogWriter("test.log", false)) //输出到文件,级别为DEBUG,文件名为test.log,每次追加该原文件
	defer l4g.Close()
	tcp := new(tcp.TCPServer)
	tcp.SetHandler(new(handler.SessionHanler))
	/**
	 *module start
	 */
	go tcp.StartTCPServer()
	/**
	 * module end
	 */
	//make一个chan用于阻塞主线程,避免程序退出
	blockMainRoutine := make(chan bool)
	<-blockMainRoutine
	//test end for git
}
