package proxy

import (
	"strconv"
	"sync"
)

//服务器节点
type ServerNode struct {
	NodeType       uint16 //服务器类型
	Ip             string //连接的IP
	Port           uint16 //连接的端口
	key            string //存储关键字
	Onlines        int    //在线用户数
	GameId         uint32 //游戏ID，为非游戏时填写0
	GameAreaId     uint32 //区域ID,为了支持分区，分服，默认可填写为1
	GameServerId   uint32 //默认为1服
	GameCode       string
	GameServerName string
	GameServerDesc string
	index          uint32
}

type ServerNodeManager struct {
	keyNodeMap map[string]*ServerNode
	lsNodeMap  map[uint32]*ServerNode
	hsNodeMap  map[uint32]*ServerNode
	gsNodeMap  map[uint32]*ServerNode
	lsIndexId  uint32
	hsIndexId  uint32
	gsIndexId  uint32
	mutex      sync.Mutex
}

func NewNodeManager() *ServerNodeManager {
	return &ServerNodeManager{keyNodeMap: make(map[string]*ServerNode),
		lsNodeMap: make(map[uint32]*ServerNode),
		hsNodeMap: make(map[uint32]*ServerNode),
		gsNodeMap: make(map[uint32]*ServerNode),
		lsIndexId: NODE_TYPE_LS * 10000,
		hsIndexId: NODE_TYPE_HS * 10000,
		gsIndexId: NODE_TYPE_GS * 10000}
}

var NodeManagerInst *ServerNodeManager = NewNodeManager()

func NewServerNode() *ServerNode {
	return &ServerNode{GameAreaId: 1, GameServerId: 1}
}

func (s *ServerNode) GetKey() string {
	if len(s.key) == 0 {
		s.key = strconv.Itoa(int(s.NodeType)) + "_" + strconv.Itoa(int(s.GameId)) + "_" + strconv.Itoa(int(s.GameAreaId)) + "_" + strconv.Itoa(int(s.GameServerId))
	}
	return s.key
}

func (s *ServerNode) SetIndex(index uint32) {
	s.index = index
}

func (s *ServerNode) GetIndex() uint32 {
	return s.index
}

func (this *ServerNodeManager) RegNodeConnection(node *ServerNode) int {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	tmp := this.keyNodeMap[node.GetKey()]
	if tmp != nil {
		return SERVER_NODE_REG_REPEAT
	}

	if node.NodeType == NODE_TYPE_LS {
		if len(this.lsNodeMap) > LS_MAX_LIMIT {
			return SERVER_NODE_REG_MAX_LIMIT
		}
		this.keyNodeMap[node.GetKey()] = node
		this.lsIndexId++
		node.SetIndex(this.lsIndexId)
		this.lsNodeMap[this.lsIndexId] = node

	} else if node.NodeType == NODE_TYPE_HS {
		if len(this.hsNodeMap) > HS_MAX_LIMIT {
			return SERVER_NODE_REG_MAX_LIMIT
		}
		this.keyNodeMap[node.GetKey()] = node
		this.hsIndexId++
		node.SetIndex(this.hsIndexId)
		this.hsNodeMap[this.hsIndexId] = node

	} else if node.NodeType == NODE_TYPE_GS {
		if len(this.gsNodeMap) > GS_MAX_LIMIT {
			return SERVER_NODE_REG_MAX_LIMIT
		}
		this.keyNodeMap[node.GetKey()] = node
		this.gsIndexId++
		node.SetIndex(this.gsIndexId)
		this.gsNodeMap[this.gsIndexId] = node
	}
	return 0
}

func (this *ServerNodeManager) UnRegNodeConnection(id uint32) {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	nodeType := id % 10000
	switch nodeType {
	case NODE_TYPE_LS:

	case NODE_TYPE_HS:
	case NODE_TYPE_GS:
	default:
		ProxyLog.Error("unknow Node Index : ", id)

	}

}

func (this *ServerNodeManager) FindLsRoute() (bool, *ServerNode) {
	return false, nil
}

func (this *ServerNodeManager) FindHsRoute() (bool, *ServerNode) {
	return false, nil
}

func (this *ServerNodeManager) FindGsRoute() (bool, *ServerNode) {
	return false, nil
}
