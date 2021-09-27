package server

import(
	"fmt"
	"net/http"
	"context"
	"github.com/pigfall/tzzGoUtil/log"
	ws "github.com/gorilla/websocket"
)


// serve once
// ready tun interface
// ready connection
func Serve(ctx context.Context,rawLogger log.Logger_Log)error{
	logger := log.NewHelper("Serve",rawLogger,log.LevelDebug)
	// { ready tun interface
	tunIfce, err := tunReady(ctx,logger)
	if err != nil{
		err = fmt.Errorf("准备 tun ifce 环境失败: %w")
		logger.Error(err)
		return err
	}
	defer tunIfce.Close()
	// }

	connCtrl := newConnCtrl()

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
			err = connCtrl.Serve(ctx,conn)
			if err != nil{
				logger.Error(err)
			}
		},
	)
	// }
}
