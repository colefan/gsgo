//Copyright 2015 gsgo Author. All Rights Reserved.
package logs

import (
	"fmt"
	"strings"
	"sync"
)

//日志等级常量
const (
	LevelTrace = iota
	LevelDebug
	LevelInfo
	LevelWarn
	LevelError
)

var levelPrefix = []string{"[T]", "[D]", "[I]", "[W]", "[E]"}

func getLoggerLevel(level string) int {
	level = strings.ToLower(level)
	switch level {
	case "trace":
		return LevelTrace
	case "debug":
		return LevelDebug
	case "info":
		return LevelInfo
	case "warn":
		return LevelWarn
	case "error":
		return LevelError
	default:
		return LevelTrace
	}
}

//日志输出终端接口
type LoggerOutputInf interface {
	Init(config string) error
	WriteMsg(level int, msg string) error
	Flush()
	Destroy()
}

//日志输出终端函数定义
type loggerType func() LoggerOutputInf

//支持的日志输出终端
var appenderMap = make(map[string]loggerType)

func RegisterLoggerAppender(name string, log loggerType) {
	if log == nil {
		panic("logs : Register provider is nil")
	}

	if _, ok := appenderMap[name]; ok {
		panic("logs: Duplicate Register provider " + name)
	}
	appenderMap[name] = log
}

type logMsg struct {
	Level int
	Msg   string
}

//logger class
type Logger4g struct {
	lock    sync.Mutex
	level   int
	outputs map[string]LoggerOutputInf
	msg     chan *logMsg
	name    string
}

//创建一个新的logger
//usage: NewLogger("defaultlogger",100)
//chanlen:可缓存的日志条数
func NewLogger(name string, chanlen int) *Logger4g {
	log1 := &Logger4g{}
	log1.level = LevelTrace
	log1.outputs = make(map[string]LoggerOutputInf)
	if chanlen < 10 {
		chanlen = 10
	}
	log1.msg = make(chan *logMsg, chanlen)
	log1.name = name
	go log1.startLogger()
	return log1
}

func (this *Logger4g) startLogger() {
	for {
		select {
		case tmpMsg := <-this.msg:
			for _, tmpOutput := range this.outputs {
				err := tmpOutput.WriteMsg(tmpMsg.Level, tmpMsg.Msg)
				if err != nil {
					fmt.Println("Error,unable to WriteMsg:", err)
				}
			}
		}
	}

}

func (this *Logger4g) SetAppender(appenderName string, config string) error {
	this.lock.Lock()
	defer this.lock.Unlock()
	if log, ok := appenderMap[appenderName]; ok {
		lg := log()
		err := lg.Init(config)
		if err != nil {
			fmt.Println("Logs.Logger4g.SetAppender: " + err.Error())
		}
		this.outputs[appenderName] = lg
	} else {
		fmt.Println("logs: unkown appendername %s \n", appenderName)
		return fmt.Errorf("logs: unkown appendername %s ", appenderName)
	}
	return nil
}

func (this *Logger4g) writeMsg(level int, msg string) {
	if level < this.level {
		return
	}

	tmpMsg := &logMsg{Level: level, Msg: msg}
	this.msg <- tmpMsg
	return
}

func (this *Logger4g) Flush() {
	for _, l := range this.outputs {
		l.Flush()
	}
}

func (this *Logger4g) Close() {
	for {
		if len(this.msg) > 0 {
			tmpMsg := <-this.msg
			for _, tmpOutput := range this.outputs {
				err := tmpOutput.WriteMsg(tmpMsg.Level, tmpMsg.Msg)
				if err != nil {
					fmt.Printf("Error unable to WriteMsg while closing logger :", err)
				}
			}
			continue
		}
		break
	}

	for _, o := range this.outputs {
		o.Flush()
		o.Destroy()
	}
}

func (this *Logger4g) SetLevel(level string) {
	this.level = getLoggerLevel(level)
}

func (this *Logger4g) GetLevel() int {
	return this.level
}

func (this *Logger4g) Trace(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	this.writeMsg(LevelTrace, msg)
}

func (this *Logger4g) Debug(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	this.writeMsg(LevelDebug, msg)
}

func (this *Logger4g) Info(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	this.writeMsg(LevelInfo, msg)
}

func (this *Logger4g) Warn(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	this.writeMsg(LevelWarn, msg)
}

func (this *Logger4g) Error(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	this.writeMsg(LevelError, msg)
}

var DefaultLogger *Logger4g = NewLogger("default", 100)

func init() {
	DefaultLogger.SetAppender("console", `{"level":0}`)
}
