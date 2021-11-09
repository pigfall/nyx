package yingying

import(
	"context"
	"github.com/pigfall/tzzGoUtil/log"
	"github.com/pigfall/tzzGoUtil/net"
)

type TransprtMsgType int

const(
		IpPacket TransprtMsgType  = 1
		Proto TransprtMsgType  = 2
)

type Transport interface{
	WriteIpPacket(ipPacketBytes []byte)(error)
	// WriteMsg(msg *proto.Msg,body interface{})(error)
	Read()(msgType TransprtMsgType,data []byte,err error)
	WriteJSON(msg interface{})(error)
}

type TransportServer interface{
	Serve(ctx context.Context,logger log.LoggerLite,connCtrl ConnCtrl,tunIfce net.TunIfce,tunIp *net.IpWithMask)error
}

type ConnCtrl interface{
	Serve(ctx context.Context, conn Transport,tunIfce net.TunIfce ) error
	GetConns() map[string]Transport
}
