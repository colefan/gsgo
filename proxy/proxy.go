package proxy

import (
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
	return &Proxy{}
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

	//开始启动代理服务
	ProxyLog.Info(">>begin to load ProxyService...")

	//判断代理服务是否启动成功

	//等待业务处理或者正常退出
	ProxyLog.Info(">>Proxy is running:)")

}
