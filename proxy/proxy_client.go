package proxy

import (
	"sync"

	"github.com/colefan/gsgo/netio"
)

type ProxyClient struct {
	status       int
	openId       string //账号
	openKey      string //账号的校验码
	lsServerNode *ServerNode
	hsServerNode *ServerNode
	gsServerNode *ServerNode
	physicalLink netio.ConnInf //客户端物理连接
}

func NewProxyClient() *ProxyClient {
	return &ProxyClient{}
}

func (c *ProxyClient) SetStatus(status int) {
	c.status = status
}

func (c *ProxyClient) GetStatus() int {
	return c.status
}

func (c *ProxyClient) SetOpenId(openId string) {
	c.openId = openId
}

func (c *ProxyClient) GetOpenId() string {
	return c.openId
}

func (c *ProxyClient) SetOpenKey(key string) {
	c.openKey = key
}

func (c *ProxyClient) GetOpenKey() string {
	return c.openKey
}

func (c *ProxyClient) SetPhysicalLink(conn netio.ConnInf) {
	c.physicalLink = conn
}

func (c *ProxyClient) GetPhysicalLink() netio.ConnInf {
	return c.physicalLink
}

//设置登录路由
func (c *ProxyClient) SetLsRoute(node *ServerNode) {
	c.lsServerNode = node
}

//获取该客户端的登录服务器路由信息
func (c *ProxyClient) GetLsRoute() *ServerNode {
	return c.lsServerNode
}

//目录服务路由
func (c *ProxyClient) SetHsRoute(node *ServerNode) {
	c.hsServerNode = node
}

//目录服务路由
func (c *ProxyClient) GetHsRoute() *ServerNode {
	return c.hsServerNode
}

//游戏服路由
func (c *ProxyClient) SetGsRoute(node *ServerNode) {
	c.gsServerNode = node
}

//游戏服路由
func (c *ProxyClient) GetGsRoute() *ServerNode {
	return c.gsServerNode
}

type ProxyClientManager struct {
	authedClients   map[uint32]*ProxyClient
	unAuthedClients map[uint32]*ProxyClient
	deleteClients   map[uint32]*ProxyClient
	mutex           sync.RWMutex
}

func newClientManager() *ProxyClientManager {
	return &ProxyClientManager{authedClients: make(map[uint32]*ProxyClient),
		unAuthedClients: make(map[uint32]*ProxyClient),
		deleteClients:   make(map[uint32]*ProxyClient)}
}

var ClientManagerInst = newClientManager()

func (this *ProxyClientManager) FindClientByUserID(uid uint32) *ProxyClient {
	this.mutex.RLock()
	defer this.mutex.RUnlock()
	tmp := this.authedClients[uid]
	return tmp

}

func (this *ProxyClientManager) FindClientByConnID(connID uint32) *ProxyClient {
	this.mutex.RLock()
	defer this.mutex.RUnlock()
	return this.unAuthedClients[connID]
}

func (this *ProxyClientManager) AddClient(client *ProxyClient) {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	this.unAuthedClients[client.GetPhysicalLink().GetConnID()] = client

}

func (this *ProxyClientManager) RemoveClientByUserID(uid uint32) {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	delete(this.authedClients, uid)
}

func (this *ProxyClientManager) RemoveClientByConnID(connID uint32) {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	delete(this.unAuthedClients, connID)
}

//物理链接与具体的用户关联后，需要将新的链接替换老链接。
//同号剔除就需要从此接口处理
func (this *ProxyClientManager) ChangeClientFromConnIDtoUserID(newClient *ProxyClient) (oldClient *ProxyClient) {
	//先从未验证的队列里删除链接
	this.mutex.Lock()
	defer this.mutex.Unlock()
	delete(this.unAuthedClients, newClient.GetPhysicalLink().GetConnID())
	oldClient = this.authedClients[newClient.GetPhysicalLink().GetUID()] //
	if oldClient != nil {
		this.deleteClients[oldClient.GetPhysicalLink().GetConnID()] = oldClient
	}
	this.authedClients[newClient.GetPhysicalLink().GetUID()] = newClient //

	return oldClient
}

func (this *ProxyClientManager) ModifyClient(connID uint32, uid uint32, client *ProxyClient) bool {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	if uid > 0 {
		this.authedClients[uid] = client
		return true
	}

	if connID > 0 {
		this.unAuthedClients[connID] = client
	}
	return false
}

func (this *ProxyClientManager) CloseClient(connID uint32, uid uint32) (bool, *ProxyClient) {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	if connID > 0 {
		deleteTmp := this.deleteClients[connID]
		if deleteTmp != nil {
			delete(this.deleteClients, connID)
			return true, deleteTmp
		}
	}

	if uid > 0 {
		authTmp := this.authedClients[uid]
		if authTmp != nil {
			delete(this.authedClients, uid)
			return true, authTmp
		}

	}

	return false, nil
}
