package server

import(
	"context"
	ws "github.com/gorilla/websocket"
"github.com/pigfall/tzzGoUtil/net"
"github.com/pigfall/tzzGoUtil/log"
)

type connCtrl struct {
	conns map[string]*ws.Conn
	tunIp *net.IpWithMask
	ipPoolIfce ipPoolIfce
	logger log.LoggerLite
}

func newConnCtrl(ipPoolIfce ipPoolIfce,rawLogger log.Logger_Log)*connCtrl {
	return &connCtrl{
		conns :make(map[string]*ws.Conn),
		ipPoolIfce:ipPoolIfce,
		logger:log.NewHelper("connCtrl",rawLogger,log.LevelDebug),
	}
}



func (this *connCtrl) Serve(ctx context.Context,conn *ws.Conn) (error){
	logger := this.logger
	// tell the client his vpn ip
	clientVPNIpNet,err := this.ipPoolIfce.GetIpNet()
	if err != nil{
		logger.Error(err)
		return err
	}
	conn.WriteJSON()
	// 
}
