package server

import(
	"github.com/pigfall/tzzGoUtil/log"
	tzNet "github.com/pigfall/tzzGoUtil/net"
	ws "github.com/gorilla/websocket"
	"fmt"
	"context"
		"net"
		"net/http"
		yy "github.com/pigfall/yingying"
		tp "github.com/pigfall/yingying/transport"
)

type tpServerWebSocket struct{
	ipToListen net.IP
	port int
}

func NewTransportServerWebSocket(ipToListen net.IP,port int) yy.TransportServer {
	return &tpServerWebSocket{
		ipToListen:ipToListen,
		port:port,
	}
}


func (this *tpServerWebSocket) Serve(ctx context.Context,logger log.LoggerLite,connCtrl yy.ConnCtrl,tunIfce tzNet.TunIfce,tunIp *tzNet.IpWithMask)(error){
	l,err := net.Listen("tcp",fmt.Sprintf(":%d",this.port))
	if err != nil{
		logger.Error(err)
		return err
	}
	defer l.Close()
	go func(){
		<-ctx.Done()
		l.Close()
	}()
	httpServer := http.NewServeMux()
	httpServer.HandleFunc(
		"/",
		func(res http.ResponseWriter,req *http.Request){
			logger.Debug("rcv request from  ",req.RemoteAddr)
			upgrader := ws.Upgrader{}
			conn,err := upgrader.Upgrade(res,req,nil)
			if err != nil {
				logger.Error(err)
				return
			}
			defer conn.Close()
			err = connCtrl.Serve(ctx,tp.NewTPWebSocket(conn),tunIfce)
			if err != nil{
				logger.Error(err)
			}
		},
	)
	err = http.Serve(l,httpServer)
	if err != nil{
		return err
	}
	return nil
}
