package client

import(
	"context"
		
	"github.com/pigfall/tzzGoUtil/async"
	"github.com/pigfall/tzzGoUtil/log"
)



func Run (
	ctx context.Context,
	rawLogger log.Logger_Log,
	cfg *RunCfg,
)error{
	ctx,cancel := context.WithCancel(ctx)
	logger := log.NewHelper("Run",rawLogger,log.LevelDebug)
	logger.Info("")
	asyncCtrl := &async.Ctrl{}
	asyncCtrl.AppendCancelFuncs(cancel)
	asyncCtrl.OnRoutineQuit(
			func(){
				asyncCtrl.Cancel()
			},
	)
	asyncCtrl.AsyncDo(
			ctx,
			func(ctx context.Context){
				handleConnData(ctx)
			},
	)

	asyncCtrl.Wait()
	return nil
}
