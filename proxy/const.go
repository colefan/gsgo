package proxy

const (
	BS_STATUS_OPENED = 1
	BS_STATUS_AUTHED = 2
	BS_STATUS_CLOSED = 3
)

const (
	NODE_TYPE_LS    = 1     //登录服务器，主要处理用户登录的合法验证
	LS_INDEX_MIN    = 10000 //最小索引
	LS_INDEX_MAX    = 19999 //最大索引
	LS_MAX_ONLINES  = 10000
	NODE_TYPE_HS    = 2     //目录服务器或者说是大厅服务器，主要提供游戏目录服务
	HS_INDEX_MIN    = 20000 //最小索引
	HS_INDEX_MAX    = 29999 //最大索引
	HS_MAX_ONLINES  = 10000
	NODE_TYPE_GS    = 3     //游戏逻辑服务器
	GS_INDEX_MIN    = 30000 //最小索引
	GS_INDEX_MAX    = 39999 //最大索引
	GS_MAX_ONLINES  = 5000
	NODE_TYPE_PROXY = 99
)

const (
	SERVER_NODE_REG_REPEAT    = 20001
	SERVER_NODE_REG_MAX_LIMIT = 20002
	PROXY_FORWARD_DECODE_ERR  = 20003
)

//管理接入服的容量
const (
	LS_MAX_LIMIT = 100 //一个代理最多可以接入100个登录服
	HS_MAX_LIMIT = 10  //一个代理最多可以接入10个目录服
	GS_MAX_LIMIT = 500 //一个代理最多可以接入500个游戏逻辑服
)
