package client

import(
	"context"
		
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
	logger.Info("")
	conn,_,err := ws.DefaultDialer.Dial(cfg.ServerAddr,nil)
	if err != nil{
		logger.Error(err)
		return err
	}
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
				handleConnData(ctx,rawLogger,conn,asyncCtrl)
			},
	)

	asyncCtrl.Wait()
	return nil
}
