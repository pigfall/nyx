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
		var tunIfce  net.TunIfce
		switch msgType {
		case yy.IpPacket:
			err = handleIpPacket(data,tunIfce)
			if err != nil{
				logger.Error(err)
			}
		case yy.Proto:
			handleConnProto(data,clientIp,tunIfce)
		default:
			panic(fmt.Errorf("Undefined msgType %v",msgType))
		}
	}
}

func handleIpPacket(data []byte,tunIfce net.TunIfce)error{
	_,err := tunIfce.Write(data)
	return err
}

func handleConnProto(data []byte,clientIp *net.IpWithMask,tunIfce net.TunIfce){

}
