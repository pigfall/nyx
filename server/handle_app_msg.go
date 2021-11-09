package server

import(
		"context"
		"github.com/pigfall/yingying/proto"
		"github.com/pigfall/tzzGoUtil/net"
)

func handleAppMsg(ctx context.Context,reqMsg *proto.Msg,clientTunnelIpGetter clientTunnelIpGetter)(*proto.Msg){
	res := &proto.Msg{}
	switch res.Id {
	default:
		panic("TODO")
	}
}

type clientTunnelIpGetter interface{
	GetClientTunnelIp()(net.IpNetFormat,error)
}

type clientTunnelIpGetterFunc func()(net.IpNetFormat,error)

func (f clientTunnelIpGetterFunc) GetClientTunnelIp()(net.IpNetFormat,error){
	return f()
}
