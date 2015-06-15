package proxy

import (
	"github.com/colefan/gsgo/gameprotocol/protocol_comm"
	"github.com/colefan/gsgo/netio"
	"github.com/colefan/gsgo/netio/packet"
)

//包内通用函数
func errorResponse(cmdId uint16, errCode uint16, userId uint32, conn netio.ConnInf) {
	if conn == nil {
		ProxyLog.Error("send error msg error,physical link is nil,req-cmd-id = ", cmdId)
		return
	}
	resp := &protocol_comm.ServerErrorNt{}
	resp.Packet = packet.NewEmptyPacket()
	resp.CmdID = protocol_comm.CMD_S_C_ERROR_NT
	resp.ID = userId
	resp.FSID = NODE_TYPE_PROXY
	resp.ReqCmdID = cmdId
	resp.ErrCode = errCode

	buf := resp.EncodePacket(64)
	conn.Write(buf.GetData())

}
