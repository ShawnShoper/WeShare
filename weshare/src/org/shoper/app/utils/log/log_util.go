package utils

import (
	l4g "github.com/alecthomas/log4go"
)

func init() {
	l4g.AddFilter("stdout", l4g.DEBUG, l4g.NewConsoleLogWriter())             //输出到控制台,级别为DEBUG
	l4g.AddFilter("file", l4g.DEBUG, l4g.NewFileLogWriter("test.log", false)) //输出到文件,级别为DEBUG,文件名为test.log,每次追加该原文件

}
func Debug(v interface{}) {
	//l4g.LoadConfiguration("log.xml")//使用加载配置文件,类似与java的log4j.propertites
	l4g.Debug(v)
	defer l4g.Close()
}
func Info(v interface{}) {
	l4g.Info(v)
	defer l4g.Close()
}
func Error(v interface{}) {
	l4g.Error(v)
	defer l4g.Close()
}
func Warn(v interface{}) {
	l4g.Warn(v)
	defer l4g.Close()
}
