package netio

import (
	"github.com/colefan/gsgo/netio/packet"
)

type PackDispatcher interface {
	DispatchMsg(data []byte, conn ConnInf)
	SessionOpen(conn ConnInf)
	SessionClose(conn ConnInf)
	SessionIdle(conn ConnInf)
	AddPackEventListener(name string, listener PackListener)
	RemovePackEventListener(name string)
}

type PackListener interface {
	HandleMsg(cmdid uint16, pack *packet.Packet, conn ConnInf)
}

type DefaultPackDispatcher struct {
	msgListeners map[string]PackListener
}

func NewDefaultPackDispatcher() *DefaultPackDispatcher {
	return &DefaultPackDispatcher{msgListeners: make(map[string]PackListener)}
}

func (this *DefaultPackDispatcher) DispatchMsg(data []byte, conn ConnInf) {
	nLen := len(data)
	if nLen <= 0 {
		return
	}

	pack := packet.Packing(data)
	if pack == nil {
		return
	}

	//fmt.Println("read msg,listeners = ",len(this.msgListeners))

	for _, v := range this.msgListeners {
		v.HandleMsg(pack.CmdID, pack, conn)
	}

}

func (this *DefaultPackDispatcher) SessionOpen(conn ConnInf) {

}

func (this *DefaultPackDispatcher) SessionClose(conn ConnInf) {

}

func (this *DefaultPackDispatcher) SessionIdle(conn ConnInf) {

}
func (this *DefaultPackDispatcher) AddPackEventListener(name string, listener PackListener) {

	if this.msgListeners == nil {
		this.msgListeners = make(map[string]PackListener)
	}
	if len(name) == 0 {
		return
	}
	if this.msgListeners != nil {
		this.msgListeners[name] = listener
	}

}

func (this *DefaultPackDispatcher) RemovePackEventListener(name string) {
	if this.msgListeners != nil {
		delete(this.msgListeners, name)
	}
}
