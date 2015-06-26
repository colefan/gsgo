//代理服务器的配置
package proxy

import (
	"fmt"

	"github.com/colefan/gsgo/config"
)

type ProxyConfig struct {
	ListenIp           string
	ListenPort         int
	ProxyName          string
	ForwardIp          string
	ForwardPort        int
	EnableHs           bool
	LoadLimit          int
	EnableReserveProxy bool   //是否支持后备代理
	ReserveProxyIp     string //后备代理IP
	ReserveProxyPort   int    //后备代理PORT
}

var ProxyConf *ProxyConfig

func init() {
	ProxyConf = &ProxyConfig{}
}

func (c *ProxyConfig) Init() error {
	cfg, err := config.NewConfig("ini", "proxyconfig.ini")
	if err != nil {
		return fmt.Errorf("load file.error->", err.Error())
	}
	c.ProxyName = cfg.String("ProxyName")
	c.ListenIp = cfg.String("ProxyListenIp")
	c.ListenPort, err = cfg.Int("ProxyPort")
	if err != nil {
		return fmt.Errorf("ProxyPort.error->", err.Error())
	}
	c.EnableHs, err = cfg.Bool("EnableHs")
	if err != nil {
		return fmt.Errorf("EnableHs.error->", err.Error())
	}
	c.LoadLimit, err = cfg.Int("LoadLimit")
	if err != nil {
		return fmt.Errorf("LoadLimit.error->", err.Error())
	}

	c.ForwardIp = cfg.String("Forward::IP")
	c.ForwardPort, err = cfg.Int("Forward::PORT")
	if err != nil {
		return fmt.Errorf("Forward::PORT.error->", err.Error())
	}

	c.EnableReserveProxy, err = cfg.Bool("ReserveProxy::Enable")
	if err != nil {
		return fmt.Errorf("ReserveProxy::Enable.error->", err.Error())
	}
	c.ReserveProxyIp = cfg.String("ReserveProxy::IP")
	c.ReserveProxyPort, err = cfg.Int("ReserveProxy::PORT")
	if err != nil {
		return fmt.Errorf("ReserveProxy::PORT.error->", err.Error())
	}

	if c.EnableReserveProxy {
		if len(c.ReserveProxyIp) < 3 || c.ReserveProxyPort == 0 {
			return fmt.Errorf("reserve proxy ip or port invalid")
		} else if c.ReserveProxyIp == c.ListenIp && c.ReserveProxyPort == c.ListenPort {
			return fmt.Errorf("reserve proxy ip and port is same as the ace")
		}

	}

	return nil
}
