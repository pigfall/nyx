package transport

import(
	"github.com/pigfall/tzzGoUtil/net"
	// "encoding/json"

  yy "github.com/pigfall/yingying"
)

type tpUDP struct{
	conn *net.UDPSock
}

func NewTransportUDP(udpSock *net.UDPSock) yy.Transport{
	return &tpUDP{
		conn:udpSock,
	}
} 

func (this *tpUDP) Read()(msgType yy.TransprtMsgType,data []byte,err error){
	// TODO
	panic("TODO")
}

func(this *tpUDP)WriteIpPacket(ipPacketBytes []byte)(error){
	panic("TODO")
}


func (this *tpUDP) WriteJSON(msg interface{})(error){
	panic("TODO")

}
