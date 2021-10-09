package server

import(
	"fmt"
	"net/http"
	"context"
	"github.com/pigfall/tzzGoUtil/log"
	stdnet "net"
	"github.com/pigfall/tzzGoUtil/net"
	"github.com/pigfall/tzzGoUtil/async"
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
	asyncCtrl := &async.Ctrl{}
	logger := log.NewHelper("Serve",rawLogger,log.LevelDebug)
	// { ready tun interface
	tunIfce,tunIp, err := tunReady(ctx,logger)
	if err != nil {
		err = fmt.Errorf("准备 tun ifce 环境失败: %w",err)
		logger.Error(err)
		return err
	}
	defer tunIfce.Close()
	asyncCtrl.AppendCancelFuncs(func(){tunIfce.Close()})
	asyncCtrl.OnRoutineQuit(func(){asyncCtrl.Cancel()})
	defer asyncCtrl.Wait()
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
	asyncCtrl.AsyncDo(
		ctx,
		func(ctx context.Context){
			var buf = make([]byte,1024*4)
			for {
				n,err := tunIfce.Read(buf)
				if err != nil{
					logger.Error(err)
					return
				}
				for _,conn := range  connCtrl.conns{
					err =conn.WriteMessage(ws.BinaryMessage,buf[:n])
					if err != nil{
						logger.Error(err)
					}
				}
			}
		},
	)
	// }

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
	asyncCtrl.AppendCancelFuncs(func(){l.Close()})
	err = http.Serve(l,httpServer)
	defer asyncCtrl.Cancel()
	if err != nil{
		logger.Error(err)
		return err
	}
	// }
	return nil
}


