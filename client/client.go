package client

import(
	"os"
	"net/url"
	"context"
	stdnet "net"
		
	"github.com/pigfall/tzzGoUtil/async"
	"github.com/pigfall/tzzGoUtil/log"
ws "github.com/gorilla/websocket"
)



func Run (
	ctx context.Context,
	rawLogger log.Logger_Log,
	cfg *RunCfg,
)error{
	ctx,cancel := context.WithCancel(ctx)
	defer cancel()
	logger := log.NewHelper("Run",rawLogger,log.LevelDebug)
	logger.Info("Connecting to address ",cfg.ServerAddr)
	svrUrl,err := url.Parse(cfg.ServerAddr)
	if err != nil{
		logger.Errorf("Config error, invalid server url %s",cfg.ServerAddr)
		os.Exit(1)
	}
	host :=svrUrl.Hostname()
	svrIp := stdnet.ParseIP(host)
	
	if svrIp == nil{
		logger.Errorf("Config error, server ip invalid %s",host)
		os.Exit(1)
	}
	conn,_,err := ws.DefaultDialer.Dial(cfg.ServerAddr,nil)
	if err != nil{
		logger.Error(err)
		return err
	}
	logger.Info("Connected")
	asyncCtrl := &async.Ctrl{}
	asyncCtrl.AppendCancelFuncs(func(){conn.Close()})
	asyncCtrl.AppendCancelFuncs(cancel)
	asyncCtrl.OnRoutineQuit(
			func(){
				asyncCtrl.Cancel()
			},
	)
	asyncCtrl.AsyncDo(
			ctx,
			func(ctx context.Context){
				handleConnData(ctx,rawLogger,conn,asyncCtrl,svrIp)
			},
	)

	asyncCtrl.Wait()
	return nil
}
