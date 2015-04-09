package logs

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	dateFmt     string = "2006-01-02"
	dateTimeFmt string = "2006-01-02-15-04-05"
	//1M的大小
	msize = 1048576
)

type MutexFileWriter struct {
	sync.Mutex
	fd *os.File
}

func (this *MutexFileWriter) SetFd(fd *os.File) {
	if this.fd != nil {
		this.fd.Close()
	}
	this.fd = fd
}

func (this *MutexFileWriter) Write(b []byte) (int, error) {
	this.Lock()
	defer this.Unlock()
	return this.fd.Write(b)
}

//usage: Init(`{"filename":"testfilename.log","maxsize":0,"daily":"true","dailycuttime":"00:00"}`)
type FileLogOutput struct {
	lg          *log.Logger
	mw          *MutexFileWriter
	FileName    string `json:"filename"`
	curFileName string
	//cut max size,单位M
	CutMaxSize         int `json:"maxsize"`
	cutmaxsize_cursize int
	//cut daily,每日切割日志
	Daily bool `json:daily`
	//05:00
	DailyCutTimeFmt    string `json:"dailycuttime"`
	dailyCutTimeSecond int
	dailyOpenDate      int
	//prefix
	Prefix  string `json:"prefix"`
	Level   int    `json:"level"`
	headLen int
	//cut
	cutLock sync.Mutex
}

func NewFileLogOutput() LoggerOutputInf {
	fileLog := &FileLogOutput{FileName: "",
		CutMaxSize:         128,
		Daily:              true,
		DailyCutTimeFmt:    "00:00",
		dailyCutTimeSecond: 0,
		Prefix:             ""}
	mw := &MutexFileWriter{}
	fileLog.mw = mw
	fileLog.lg = log.New(mw, "", log.Ldate|log.Ltime)
	return fileLog
}

func (this *FileLogOutput) Init(config string) error {
	if len(config) <= 0 {
		return nil
	}
	err := json.Unmarshal([]byte(config), this)
	if err != nil {
		return err
	}
	this.headLen = len(this.Prefix) + len("[T]2012-06-06 09:88:88 ")
	return this.startLogger()
}

func (this *FileLogOutput) startLogger() error {
	this.cutmaxsize_cursize = 0
	fd, err := this.createLogFile()
	if err != nil {
		return err
	}
	this.mw.SetFd(fd)
	return this.initFd()
}

func (this *FileLogOutput) initFd() error {
	fd := this.mw.fd
	finfo, err := fd.Stat()
	if err != nil {
		return fmt.Errorf("get stat err: %s\n", err)
	}
	this.cutmaxsize_cursize = int(finfo.Size())
	this.dailyOpenDate = time.Now().Day()
	return nil
}

func (this *FileLogOutput) createLogFile() (*os.File, error) {
	// Open the log file
	filePath := ""
	postfix := "log"
	posIndex := strings.LastIndex(this.FileName, ".")
	if posIndex > 0 && posIndex < len(this.FileName)-1 {
		filePath = this.FileName[0:posIndex]
		postfix = this.FileName[posIndex+1:]
	}

	if len(filePath) == 0 {
		filePath = this.FileName
	}

	filePath = filePath + "_" + time.Now().Format(dateTimeFmt) + "." + postfix
	this.curFileName = filePath
	//如果文件已经存在，则会被覆盖，这个地方可以考虑加文件下标来区别开来
	fd, err := os.OpenFile(filePath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0660)
	return fd, err
}

func (this *FileLogOutput) WriteMsg(level int, msg string) error {
	if level < this.Level {
		return nil
	}
	n := this.headLen + len(msg)
	this.docheck(n)
	this.lg.SetPrefix(this.Prefix + levelPrefix[level])
	this.lg.Output(2, msg)
	return nil
}
func (this *FileLogOutput) docheck(size int) {
	this.cutLock.Lock()
	defer this.cutLock.Unlock()

	if (this.Daily && time.Now().Day() != this.dailyOpenDate) ||
		(this.CutMaxSize > 0 && this.cutmaxsize_cursize >= this.CutMaxSize) {
		if err := this.doCutLog(); err != nil {
			fmt.Fprintf(os.Stderr, "FileLogOutput(%q):%s\n", this.FileName, err)
			return
		}
	}
	this.cutmaxsize_cursize += size
}

func (this *FileLogOutput) doCutLog() error {
	println("docutlog")
	this.mw.Lock()
	defer this.mw.Unlock()
	this.mw.fd.Close()
	// re-start logger
	err := this.startLogger()
	if err != nil {
		return fmt.Errorf("Rotate StartLogger: %s\n", err)
	}
	return nil
}

func (this *FileLogOutput) Flush() {
	this.mw.fd.Sync()
}

func (this *FileLogOutput) Destroy() {
	this.mw.fd.Close()
}

func init() {
	println("register file log output")
	RegisterLoggerAppender("file", NewFileLogOutput)
}
