package logs

import (
	"encoding/json"
	"log"
	"os"
)

type ConsoleOutput struct {
	lg     *log.Logger
	Level  int    `json:"level"`
	Prefix string `json:"prefix"`
}

func NewConsoleOutput() LoggerOutputInf {
	co := &ConsoleOutput{lg: log.New(os.Stdout, "", log.Ldate|log.Ltime), Level: LevelTrace}
	return co
}

func (this *ConsoleOutput) Init(jsonconfig string) error {
	if len(jsonconfig) == 0 {
		return nil
	}
	//println("json", jsonconfig)
	return json.Unmarshal([]byte(jsonconfig), this)
}

func (this *ConsoleOutput) WriteMsg(level int, msg string) error {
	if level < this.Level {
		return nil
	}
	this.lg.SetPrefix(this.Prefix + levelPrefix[level])
	this.lg.Output(2, msg)
	return nil
}

//console无须destroy
func (this *ConsoleOutput) Destroy() {
}

//console无须flush
func (this *ConsoleOutput) Flush() {
}

func init() {
	println("register console log output")
	RegisterLoggerAppender("console", NewConsoleOutput)
}
