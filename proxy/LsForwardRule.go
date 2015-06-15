package proxy

import (
	"github.com/colefan/gsgo/gameprotocol/protocol_login"
	"github.com/colefan/gsgo/gameprotocol/protocol_proxy"
	"github.com/colefan/gsgo/netio"
	"github.com/colefan/gsgo/netio/packet"
)

type LsForwardRule struct {
	ForwardRule
}

func newLsForwardRule() *LsForwardRule {
	return &LsForwardRule{}
}

var lsForwarderInst = newLsForwardRule()

func (rule *LsForwardRule) FowardToServer(client *ProxyClient, cmdID uint16, pack *packet.Packet, conn netio.ConnInf) {
	if code := rule.preRule(pack); code != 0 {
		errorResponse(cmdID, uint16(code), pack.ID, conn)
		return
	}
	server := client.GetLsRoute()
	if server == nil {
		errorResponse(cmdID, EC_LSROUTE_IS_NIL, pack.ID, conn)
		return
	}
	serverLink := server.GetPhysicalLink()
	if serverLink == nil {
		errorResponse(cmdID, EC_LSROUTE_IS_NIL, pack.ID, conn)
		return
	}

	if serverLink.GetBsStatus() != BS_STATUS_AUTHED {
		errorResponse(cmdID, EC_LSROUTE_IS_UNAUTHED, pack.ID, conn)
		return
	}
	serverLink.Write(pack.GetClientFromRawData())
	//rule.postRule(pack)//TODO 没有必要，但留有接口

}

func (rule *LsForwardRule) FowardToClient(cmdID uint16, pack *packet.Packet, conn netio.ConnInf) {
	clientID := pack.ID
	var client *ProxyClient
	if cmdID <= protocol_login.CMD_C_LOGIN_VALID_RESP {
		client = ClientManagerInst.FindClientByConnID(clientID)
	} else {
		client = ClientManagerInst.FindClientByUserID(clientID)
	}

	if client == nil {
		ProxyLog.Error("LS msg can't find physical conn,conn_id = ", clientID, " cmd_id = ", cmdID)
		return
	}

	switch cmdID {
	case protocol_proxy.CMD_C_P_USER_OFFLINE_REQ:
		client.GetPhysicalLink().Close()
	case protocol_login.CMD_C_LOGIN_VALID_RESP:
		if pack.ErrCode == 0 {
			loginValid := &protocol_login.LoginValidResp{}
			loginValid.Packet = pack
			if loginValid.DecodePacket() {
				client.GetPhysicalLink().SetBsStatus(BS_STATUS_AUTHED)
				client.GetPhysicalLink().SetUID(loginValid.UserId)
				oldClient := ClientManagerInst.ChangeClientFromConnIDtoUserID(client)
				if oldClient != nil && oldClient.GetPhysicalLink() != nil {
					oldClient.GetPhysicalLink().Close()
				}

			} else {
				loginValid.Packet = packet.NewEmptyPacket()
				loginValid.ErrCode = PROXY_FORWARD_DECODE_ERR
				loginValid.CmdID = protocol_login.CMD_C_LOGIN_VALID_RESP
				buf := loginValid.EncodePacket(128)
				client.GetPhysicalLink().Write(buf.GetData())
			}

		} else {
			client.GetPhysicalLink().Write(pack.GetClientFromRawData())
		}

	default:
		client.GetPhysicalLink().Write(pack.GetClientFromRawData())
	}

}
