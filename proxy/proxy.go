package proxy

import (
	"github.com/colefan/gsgo/console"
	"github.com/colefan/gsgo/logs"
	"github.com/colefan/gsgo/netio/qos"
)

//代理服务器，管理整个代理服务器的运行
var ProxyLog = logs.NewLogger("proxylog", 100)
var ProxyQos = netqos.NewServerQos()

type Proxy struct {
	p *ProxyService
	n *NodeService
}

func NewProxy() *Proxy {
	p := &Proxy{}
	p.n = NewNodeService()
	p.p = NewProxyService()
	return p
}

func (this *Proxy) Run() {
	ProxyLog.Info(">>begin to load proxyconfig.ini...")
	err := ProxyConf.Init()
	if err != nil {
		ProxyLog.Error(">>load proxyconfig.ini error:(", err)
		return
	}
	ProxyLog.Info(">>proxyconfig.ini load successful :)")
	//开始启动节点服务
	ProxyLog.Info(">>begin to load NodeService...")
	//判断节点服务是否启动成功
	err = this.n.InitService()
	if err != nil {
		ProxyLog.Error(">>load NodeService error,", err)
		return
	}

	//开始启动代理服务
	ProxyLog.Info(">>begin to load ProxyService...")
	err = this.p.InitService()
	//判断代理服务是否启动成功
	if err != nil {
		ProxyLog.Info(">>load ProxyService error,", err)
		return
	}
	//等待业务处理或者正常退出
	ProxyLog.Info(">>Proxy is running:)")
	console.CheckInput()
	this.ShutDown()

}

func (this *Proxy) ShutDown() {
	if this.p != nil {
		this.p.Close()
	}

	if this.n != nil {
		this.n.Close()
	}
}

func init() {
	ProxyLog.SetAppender("console", `{"level":0}`)
}
