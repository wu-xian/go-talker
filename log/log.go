package log

import "github.com/astaxie/beego/logs"

var Logger *logs.BeeLogger

func init() {
	Logger = logs.NewLogger(0)
	//Logger.SetLogger("console", "")
	//Logger.SetLogger("file", `{"filename":"test.log"}`)
}
