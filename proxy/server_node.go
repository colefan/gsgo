package proxy

import (
	"strconv"
	"sync"

	"github.com/colefan/gsgo/netio"
)

//服务器节点
type ServerNode struct {
	NodeType       uint16 //服务器类型
	Ip             string //连接的IP
	Port           uint16 //连接的端口
	key            string //存储关键字
	onlines        int    //在线用户数
	GameId         uint32 //游戏ID，为非游戏时填写0
	GameAreaId     uint32 //区域ID,为了支持分区，分服，默认可填写为1
	GameServerId   uint32 //默认为1服
	GameCode       string
	GameServerName string
	GameServerDesc string
	index          uint32
	conn           netio.ConnInf
	mutex          sync.RWMutex
}

type ServerNodeManager struct {
	keyNodeMap map[string]*ServerNode
	lsNodeMap  map[uint32]*ServerNode
	lsNodeKeys []uint32
	hsNodeMap  map[uint32]*ServerNode
	hsNodeKeys []uint32
	gsNodeMap  map[uint32]*ServerNode
	gsNodeKeys []uint32
	lsIndexId  uint32
	hsIndexId  uint32
	gsIndexId  uint32
	mutex      sync.RWMutex
}

func NewNodeManager() *ServerNodeManager {
	return &ServerNodeManager{keyNodeMap: make(map[string]*ServerNode),
		lsNodeMap:  make(map[uint32]*ServerNode),
		lsNodeKeys: make([]uint32, 0, 16),
		hsNodeMap:  make(map[uint32]*ServerNode),
		hsNodeKeys: make([]uint32, 0, 16),
		gsNodeMap:  make(map[uint32]*ServerNode),
		gsNodeKeys: make([]uint32, 0, 32),
		lsIndexId:  LS_INDEX_MIN,
		hsIndexId:  HS_INDEX_MIN,
		gsIndexId:  GS_INDEX_MIN}
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

func (s *ServerNode) GetOnlines() int {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.onlines
}

func (s *ServerNode) IncreaseOnlines() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.onlines++
}

func (s *ServerNode) DescreseOnlines() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.onlines--
	if s.onlines < 0 {
		s.onlines = 0
	}
}

func (s *ServerNode) SetPhysicalLink(conn netio.ConnInf) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.conn = conn
}

func (s *ServerNode) GetPhysicalLink() netio.ConnInf {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.conn
}

func (this *ServerNodeManager) getNextLsIndex() uint32 {
	this.lsIndexId++
	if this.lsIndexId <= LS_INDEX_MIN || this.lsIndexId > LS_INDEX_MAX {
		this.lsIndexId = LS_INDEX_MIN + 1
	}
	for {
		tmp := this.lsNodeMap[this.lsIndexId]
		if tmp != nil {
			this.lsIndexId++
			if this.lsIndexId > LS_INDEX_MAX {
				this.lsIndexId = LS_INDEX_MIN + 1
			}
		} else {
			return this.lsIndexId
		}
	}
}

func (this *ServerNodeManager) getNextHsIndex() uint32 {
	this.hsIndexId++
	if this.hsIndexId <= HS_INDEX_MIN || this.hsIndexId > HS_INDEX_MAX {
		this.hsIndexId = HS_INDEX_MIN + 1
	}
	for {
		tmp := this.hsNodeMap[this.hsIndexId]
		if tmp != nil {
			this.hsIndexId++
			if this.hsIndexId > HS_INDEX_MAX {
				this.hsIndexId = HS_INDEX_MIN + 1
			}
		} else {
			return this.hsIndexId
		}
	}
}

func (this *ServerNodeManager) getNextGsIndex() uint32 {
	this.gsIndexId++
	if this.gsIndexId <= GS_INDEX_MIN || this.gsIndexId > GS_INDEX_MAX {
		this.gsIndexId = GS_INDEX_MIN + 1
	}
	for {
		tmp := this.gsNodeMap[this.gsIndexId]
		if tmp != nil {
			this.gsIndexId++
			if this.gsIndexId > GS_INDEX_MAX {
				this.gsIndexId = GS_INDEX_MIN + 1
			}
		} else {
			return this.gsIndexId
		}
	}
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
		index := this.getNextLsIndex()
		node.SetIndex(index)
		this.lsNodeMap[index] = node
		this.lsNodeKeys = append(this.lsNodeKeys, index)

	} else if node.NodeType == NODE_TYPE_HS {
		if len(this.hsNodeMap) > HS_MAX_LIMIT {
			return SERVER_NODE_REG_MAX_LIMIT
		}
		this.keyNodeMap[node.GetKey()] = node
		index := this.getNextHsIndex()
		node.SetIndex(this.hsIndexId)
		this.hsNodeMap[this.hsIndexId] = node
		this.hsNodeKeys = append(this.hsNodeKeys, index)

	} else if node.NodeType == NODE_TYPE_GS {
		if len(this.gsNodeMap) > GS_MAX_LIMIT {
			return SERVER_NODE_REG_MAX_LIMIT
		}
		this.keyNodeMap[node.GetKey()] = node
		index := this.getNextGsIndex()
		node.SetIndex(index)
		this.gsNodeMap[index] = node
		this.gsNodeKeys = append(this.gsNodeKeys, index)
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
	this.mutex.RLock()
	defer this.mutex.RUnlock()
	mapLen := len(this.lsNodeKeys)
	if mapLen <= 0 {
		return false, nil
	} else if mapLen == 1 {
		val, _ := this.lsNodeMap[this.lsNodeKeys[0]]
		if val != nil {
			val.IncreaseOnlines()
			return true, val
		}

	} else {
		for _, tmpKey := range this.lsNodeKeys {
			if val, ok := this.lsNodeMap[tmpKey]; ok {
				if val.GetOnlines() < LS_MAX_ONLINES {
					val.IncreaseOnlines()
					return true, val
				}
			}

		}

		val, _ := this.lsNodeMap[this.lsNodeKeys[0]]
		if val != nil {
			val.IncreaseOnlines()
			return true, val
		}
	}
	return false, nil
}

func (this *ServerNodeManager) FindHsRoute() (bool, *ServerNode) {
	this.mutex.RLock()
	defer this.mutex.RUnlock()
	keyLen := len(this.hsNodeKeys)
	if keyLen <= 0 {
		return false, nil
	} else if keyLen == 1 {
		val, _ := this.hsNodeMap[this.hsNodeKeys[0]]
		if val != nil {
			val.IncreaseOnlines()
			return true, val
		}
	} else {
		for _, tmpKey := range this.hsNodeKeys {
			val, _ := this.hsNodeMap[tmpKey]
			if val != nil {
				if val.GetOnlines() < HS_MAX_ONLINES {
					val.IncreaseOnlines()
					return true, val
				}
			}
		}

		val, _ := this.hsNodeMap[this.hsNodeKeys[0]]
		if val != nil {
			val.IncreaseOnlines()
			return true, val
		}
	}
	return false, nil
}

func (this *ServerNodeManager) FindGsRoute() (bool, *ServerNode) {
	this.mutex.RLock()
	defer this.mutex.Unlock()
	keyLen := len(this.gsNodeKeys)
	if keyLen <= 0 {
		return false, nil
	} else if keyLen == 0 {
		val, _ := this.gsNodeMap[this.gsNodeKeys[0]]
		if val != nil {
			return true, val
		}
	} else {
		for _, tmpKey := range this.gsNodeKeys {
			val, _ := this.gsNodeMap[tmpKey]
			if val != nil {
				if val.GetOnlines() < GS_MAX_ONLINES {
					val.IncreaseOnlines()
					return true, val
				}
			}
		}

		val, _ := this.gsNodeMap[this.gsNodeKeys[0]]
		if val != nil {
			val.IncreaseOnlines()
			return true, val
		}

	}

	return false, nil
}
