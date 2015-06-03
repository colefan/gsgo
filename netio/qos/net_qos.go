package netqos

import (
	"fmt"
	"time"
)

type QosInf interface {
	ShowQos()
	Stat()
	StatAccpetConns()
	StatCloseConns()
	StatReadMsgs()
	StatWriteMsgs()
	IsEnable() bool
	SetEnable(b bool)
}

type ServerQos struct {
	timeSleep       int //每隔多少秒开始统计
	timeStartSecond time.Time
	acceptedConns   int //总共接收的连接数
	closedConns     int //目前还活着的连接数
	readMsgNums     int //总共接收的消息数
	writeMsgNums    int //总共写出去的消息数
	isEnable        bool
}

func NewServerQos() *ServerQos {
	q := &ServerQos{timeSleep: 5}
	return q
}

func (this *ServerQos) ShowQos() {
	t := time.Now().Sub(this.timeStartSecond)
	if this.isEnable {
		fmt.Println("QOS:\tTimes[", t.String(), "],\tConns[", this.acceptedConns, "],\tAlives[", this.acceptedConns-this.closedConns, "],\tReadIO[", this.readMsgNums, "],\tWriteIO[", this.writeMsgNums, "]\t......\n")
	}

}

func (this *ServerQos) Stat() {
	this.timeStartSecond = time.Now()
	go func() {
		timer1 := time.NewTicker(time.Duration(this.timeSleep) * time.Second)
		for {
			select {
			case <-timer1.C:
				go this.ShowQos()
			}
		}
	}()
}

func (this *ServerQos) StatAccpetConns() {
	this.acceptedConns++
}

func (this *ServerQos) StatReadMsgs() {
	this.readMsgNums++
}

func (this *ServerQos) StatWriteMsgs() {
	this.writeMsgNums++
}

func (this *ServerQos) StatCloseConns() {
	this.closedConns++
}

func (this *ServerQos) IsEnable() bool {
	return this.isEnable
}

func (this *ServerQos) SetEnable(b bool) {
	this.isEnable = b
}
