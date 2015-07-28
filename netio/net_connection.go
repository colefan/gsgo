package netio

import (
	"fmt"
	"net"
	"sync"

	"github.com/colefan/gsgo/logs"
	"github.com/colefan/gsgo/netio/iobuffer"
	"github.com/colefan/gsgo/netio/qos"
)

type SessionHandlerInf interface {
	GetPackParser() PackParser
	GetPackDispatcher() PackDispatcher
}

//客户端连接接口
type ConnInf interface {
	Start()
	read()
	Write([]byte)
	write()
	Close()
	handshake()
	SetQos(qos netqos.QosInf)
	GetQos() netqos.QosInf
	SetConnID(id uint32)
	GetConnID() uint32
	SetBsStatus(s int)
	GetBsStatus() int
	GetRemoteIp() string
	BindObj(interface{})
	GetBindObj() interface{}
	SetUID(id uint32)
	GetUID() uint32
	SetOpenID(openId string)
	GetOpenID() string
}

type Connection struct {
	c              net.Conn          //物理连接
	s              SessionHandlerInf //所属服务器,或客户端
	handshaked     bool              //握手是否完成
	handshakecount int               //握手次数
	address        string
	readBuff       *iobuffer.InBuffer //读取缓存区
	writeBuff      *iobuffer.InBuffer //写入缓冲区
	writeMux       sync.Mutex
	statusMux      sync.Mutex //状态值的mutex
	status         int
	id             uint32
	readMsgs       int
	qos            netqos.QosInf
	bsStatus       int //业务状态
	remoteIp       string
	bindObj        interface{}
	uid            uint32
	openId         string
}

func NewConnection(c net.Conn, s SessionHandlerInf) *Connection {
	return &Connection{c: c,
		s:         s,
		status:    SESSION_STATUS_INIT,
		readBuff:  iobuffer.NewInBuffer(iobuffer.SIZE_1_K, iobuffer.SIZE_64_K),
		writeBuff: iobuffer.NewInBuffer(iobuffer.SIZE_1_K, iobuffer.SIZE_1_M)}
}

func (this *Connection) GetConnID() uint32 {
	return this.id
}

func (this *Connection) IncreaseNums() {
	this.readMsgs++
}

func (this *Connection) ReadNums() int {
	return this.readMsgs
}

func (this *Connection) SetConnID(id uint32) {
	this.id = id
}

func (this *Connection) setStatus(s int) {
	this.statusMux.Lock()
	defer this.statusMux.Unlock()
	this.status = s
}

func (this *Connection) getStatus() int {
	this.statusMux.Lock()
	defer this.statusMux.Unlock()
	return this.status
}
func (this *Connection) Start() {
	//read first files

	this.address = this.c.RemoteAddr().String()
	this.setStatus(SESSION_STATUS_OPEN)
	//handshark
	if !this.handshaked {
		this.handshake()
	}
	//parsepack
	go this.read()
	//go this.write()
}

func (this *Connection) read() {
	//this.c.SetReadDeadline(time.Millisecond * 5)
	//this.c.SetReadDeadline(time.Now().Add(time.Millisecond * 10))

	defer this.Close()
	for {
		buf := make([]byte, 1024)
		n, err := this.c.Read(buf)

		if err != nil {
			logs.DefaultLogger.Error("Read data from conn err, err = ", err)
			break
		}

		data := this.s.GetPackParser().ParseMsg(buf, n, this.readBuff)
		if len(data) > 0 {
			if this.qos != nil {
				this.qos.StatReadMsgs()
			}
			go this.s.GetPackDispatcher().DispatchMsg(data, this)
		}

	}
}

func (this *Connection) write() {

	this.writeMux.Lock()
	defer this.writeMux.Unlock()

	for {

		if this.writeBuff.GetBuffLen() > 0 {
			tmp := this.writeBuff.CutPackData(this.writeBuff.GetBuffLen())
			if len(tmp) > 0 {
				_, err := this.c.Write(tmp)
				if err != nil {
					logs.DefaultLogger.Error("connection write error,", err)
					this.Close()
					return
				}
			}
		} else {
			return
		}
		if this.getStatus() == SESSION_STATUS_CLOSED {
			break
		}
	}
}

func (this *Connection) Write(data []byte) {
	this.writeMux.Lock()
	this.writeBuff.AppendData(data)
	if this.qos != nil {
		this.qos.StatWriteMsgs()
	}
	this.writeMux.Unlock()
	go this.write()
}

func (this *Connection) Close() {
	this.statusMux.Lock()
	defer this.statusMux.Unlock()
	fmt.Println("connetion close()")
	if this.status != SESSION_STATUS_CLOSED {
		err := this.c.Close()
		if err != nil {
			logs.DefaultLogger.Error("connection close error, ", err)
		}

		this.status = SESSION_STATUS_CLOSED
		fmt.Println("sessionclose()")
		this.s.GetPackDispatcher().SessionClose(this)
	}
}

func (this *Connection) handshake() {
	this.handshaked = true
	this.handshakecount++
	//过滤掉前面的128个字节,根据实际的需要进行修改
	//	filterBytes := make([]byte, 128)
	//	if !this.handshaked {
	//		n, err := io.ReadFull(this.c, filterBytes)
	//		if err != nil {
	//			logs.DefaultLogger.Error("handshake error")
	//		}
	//		println("n = ", n)
	//		this.handshaked = true
	//		this.handshakecount++
	//	}

}

func (this *Connection) SetQos(qos netqos.QosInf) {
	this.qos = qos
}

func (this *Connection) GetQos() netqos.QosInf {
	return this.qos
}

func (this *Connection) SetBsStatus(s int) {
	this.statusMux.Lock()
	defer this.statusMux.Unlock()
	this.bsStatus = s
}

func (this *Connection) GetBsStatus() int {
	this.statusMux.Lock()
	defer this.statusMux.Unlock()
	return this.bsStatus
}

func (this *Connection) GetRemoteIp() string {
	if len(this.remoteIp) == 0 {
		this.remoteIp = this.c.RemoteAddr().String()
	}
	return this.remoteIp
}

func (this *Connection) BindObj(obj interface{}) {
	this.bindObj = obj
}

func (this *Connection) GetBindObj() interface{} {
	return this.bindObj
}

func (this *Connection) SetUID(id uint32) {
	this.uid = id
}

func (this *Connection) GetUID() uint32 {
	return this.uid
}

func (this *Connection) SetOpenID(openId string) {
	this.openId = openId
}

func (this *Connection) GetOpenID() string {
	return this.openId
}
