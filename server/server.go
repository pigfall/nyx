package server

import(
	"fmt"
	"net/http"
	"context"
	"github.com/pigfall/tzzGoUtil/log"
	stdnet "net"
	"github.com/pigfall/tzzGoUtil/net"
	ws "github.com/gorilla/websocket"
)



// serve once
// ready tun interface
// ready connection
func Serve(
	ctx context.Context,
	rawLogger log.Logger_Log,
	cfg *ServeCfg,
)error{
	logger := log.NewHelper("Serve",rawLogger,log.LevelDebug)
	// { ready tun interface
	tunIfce,tunIp, err := tunReady(ctx,logger)
	if err != nil{
		err = fmt.Errorf("准备 tun ifce 环境失败: %w",err)
		logger.Error(err)
		return err
	}
	defer tunIfce.Close()
	// }
	ipPool,err :=net.NewIpPool(
		tunIp.BaseIpNet(),
		[]*net.IpWithMask{
			tunIp,
		},
	)
	if err != nil{
		panic(err)
	}
	connCtrl := newConnCtrl(ipPool,rawLogger)

	// { socket listen
	loggerHttpSvr := log.NewHelper("httpServer",rawLogger,log.LevelDebug)
	httpServer := http.NewServeMux()
	httpServer.HandleFunc(
		"/",
		func(res http.ResponseWriter,req *http.Request){
			loggerHttpSvr.Debug("rcv request from  ",req.RemoteAddr)
			upgrader := ws.Upgrader{}
			conn,err := upgrader.Upgrade(res,req,nil)
			if err != nil {
				loggerHttpSvr.Error(err)
				return
			}
			err = connCtrl.Serve(ctx,conn,tunIfce)
			if err != nil{
				logger.Error(err)
			}
		},
	)
	l,err := stdnet.Listen("tcp",fmt.Sprintf(":%d",cfg.Port))
	if err != nil{
		return err
	}
	err = http.Serve(l,httpServer)
	if err != nil{
		logger.Error(err)
		return err
	}
	// }
	return nil
}


