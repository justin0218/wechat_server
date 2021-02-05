package store

import (
	"github.com/astaxie/beego/logs"
)

type Log struct {
}

func (s *Log) Get() *logs.BeeLogger {
	logOnce.Do(func() {
		logs.SetLevel(0)
		logClient = logs.NewLogger()
		err := logClient.SetLogger(logs.AdapterFile, `{"filename":"./logs/log.log"}`)
		if err != nil {
			panic(err)
		}
		logClient.EnableFuncCallDepth(true)
	})
	return logClient
}
