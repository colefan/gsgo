package netio

const (
	//服务器状态 初始化
	SERVER_STATUS_INIT int = 0
	//服务器状态 侦听状态
	SERVER_STATUS_LISTENING int = 1
	//服务器状态 已经关闭
	SERVER_STATUS_CLOSED int = 2
)

const (
	//会话状态 初始化
	SESSION_STATUS_INIT int = 0
	//会话状态 已经连接
	SESSION_STATUS_OPEN int = 1
	//会话状态 已关闭
	SESSION_STATUS_CLOSED int = 2
)
