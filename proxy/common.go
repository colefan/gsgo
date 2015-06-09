package proxy

import (
	"github.com/colefan/gsgo/gameprotocol/protocol_proxy"
	"github.com/colefan/gsgo/netio"
	"github.com/colefan/gsgo/netio/packet"
)

//包内通用函数
func errorResponse(cmdId uint16, errCode uint16, userId uint32, conn netio.ConnInf) {
	if conn == nil {
		ProxyLog.Error("send error msg error,physical link is nil,req-cmd-id = ", cmdId)
		return
	}
	resp := &protocol_proxy.ProxyErrorNt{}
	resp.Packet = packet.NewEmptyPacket()
	resp.CmdID = protocol_proxy.CMD_P_C_PROXY_ERROR_NT
	resp.ID = userId
	resp.FSID = NODE_TYPE_PROXY
	resp.ReqCmdID = cmdId
	resp.ErrCode = errCode

	buf := resp.EncodePacket(64)
	conn.Write(buf.GetData())

}
