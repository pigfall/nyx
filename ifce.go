package yingying

import(
	"github.com/pigfall/yingying/proto"
)

type TransprtMsgType int

const(
		IpPacket TransprtMsgType  = 1
		Proto TransprtMsgType  = 2
)

type Transport interface{
	WriteIpPacket(ipPacketBytes []byte)(error)
	WriteMsg(msg *proto.Msg,body interface{})(error)
	Read()(msgType TransprtMsgType,data []byte,err error)
}
