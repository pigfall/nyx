package server

import(
	"encoding/json"
		"context"
		"github.com/pigfall/yingying/proto"
		"github.com/pigfall/tzzGoUtil/net"
		"log"
)

func handleAppMsg(ctx context.Context,reqMsg *proto.Msg,clientTunnelIpGetter clientTunnelIpGetter)(*proto.Msg){
	res := &proto.Msg{}
	var body interface{}
	switch res.Id {
	case proto.ID_C2S_QUERY_IP:
		log.Println("rcv query ip")
		res.Id = proto.ID_S2C_QUERY_IP
		clientTunnelIp,err := clientTunnelIpGetter.GetClientTunnelIp()
		if err != nil{
			res.ErrMsg = err.Error()
			res.ErrReason = err.Error()
			return res
		}
		body = &proto.C2S_QueryIp{IpNet:string(clientTunnelIp)}
	default:
		panic("TODO")
	}
	if body != nil{
		bodyBytes,err := json.Marshal(body)
		if err != nil{
			panic(err)
		}
		res.Body = bodyBytes
	}
	return res
}

type clientTunnelIpGetter interface{
	GetClientTunnelIp()(net.IpNetFormat,error)
}

type clientTunnelIpGetterFunc func()(net.IpNetFormat,error)

func (f clientTunnelIpGetterFunc) GetClientTunnelIp()(net.IpNetFormat,error){
	return f()
}
