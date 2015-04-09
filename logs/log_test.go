package logs

import (
	"testing"
	"time"
)

func TestConsoleLog(t *testing.T) {
	log := NewLogger("test", 100)
	log.SetAppender("console", `{"level":0,"prefix":"[cdc]"}`)
	//time.Sleep(time.Second * 2)
	log.Trace("trace")
	log.Debug("debug")
	log.Info("info")
	//log.Close()
}
func TestFileLog(t *testing.T) {
	log := NewLogger("testfile", 100)
	log.SetAppender("file", `{"level":1,"maxsize":10,"filename":"a.log"}`)
	log.SetAppender("console", `{"level":0}`)
	log.Debug("test debug,haha")
	time.Sleep(time.Second * 1)
	log.Info("test Info haha")
	time.Sleep(time.Second * 1)
	log.Info("test error ,xx")
	time.Sleep(time.Second * 1)
	log.Error("test error1 ,xx")
	time.Sleep(time.Second * 1)
	log.Error("test error2 ,xx")
	time.Sleep(time.Second * 1)
	log.Error("test error3 ,xx")
	log.Close()
}
