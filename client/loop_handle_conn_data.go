package client

import(
	"fmt"
		"context"
ws "github.com/gorilla/websocket"
	"github.com/pigfall/tzzGoUtil/log"
	"github.com/pigfall/tzzGoUtil/async"
	"github.com/pigfall/tzzGoUtil/net"
 yy "github.com/pigfall/yingying"
  "github.com/pigfall/yingying/transport"
)


func handleConnData(
	ctx context.Context,
	rawLogger log.Logger_Log,
	conn *ws.Conn,
	asyncCtrl *async.Ctrl,
){
	tp := transport.NewTPWebSocket(conn)
	logger := log.NewHelper("handleConnData",rawLogger,log.LevelDebug)
	for {
		msgType,data,err := tp.Read()
		if err != nil{
			logger.Error(err)
			return 
		}
		var clientIp *net.IpWithMask
		switch msgType {
		case yy.IpPacket:
			handleIpPacket(data)
		case yy.Proto:
			handleConnProto(data,clientIp)
		default:
			panic(fmt.Errorf("Undefined msgType %v",msgType))
		}
	}
}

func handleIpPacket(data []byte){

}

func handleConnProto(data []byte,clientIp *net.IpWithMask){

}
