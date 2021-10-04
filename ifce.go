package yingying

import(
"github.com/pigfall/yingying/proto"
)


type Transport interface{
	WriteIpPacket(ipPacketBytes []byte)(error)
	WriteMsg(msg *proto.Msg,body interface{})(error)
}
