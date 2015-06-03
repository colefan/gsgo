//代理服务器的配置
package proxy

import (
	"github.com/colefan/gsgo/config"
)

type ProxyConfig struct {
	ListenIp    string
	ListenPort  int
	ProxyName   string
	ForwardIp   string
	ForwardPort int
	EnableHs    bool
}

var ProxyConf *ProxyConfig

func init() {
	ProxyConf = &ProxyConfig{}
}

func (c *ProxyConfig) Init() error {
	cfg, err := config.NewConfig("ini", "proxyconfig.ini")
	if err != nil {
		return err
	}
	c.ProxyName = cfg.String("ProxyName")
	c.ListenIp = cfg.String("ProxyListenIp")
	c.ListenPort, err = cfg.Int("ProxyPort")
	if err != nil {
		return err
	}
	c.EnableHs, err = cfg.Bool("EnableHs")
	if err != nil {
		return err
	}

	c.ForwardIp = cfg.String("Forward::IP")
	c.ForwardPort, err = cfg.Int("Forward::PORT")
	if err != nil {
		return err
	}

	return nil
}
